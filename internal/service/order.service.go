package service

import (
	"BE/internal/domain"
	"BE/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
}

func NewOrderService(orderRepo *repository.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

type CreateOrderRequest struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	ProductName string `json:"product_name"`
	ProjectType string `json:"project_type"`
	Budget      string `json:"budget"`
	Description string `json:"description"`
}

// GetOrdersRequest untuk filter dan pagination
type GetOrdersRequest struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Status string `json:"status"`
	Search string `json:"search"`
	From   string `json:"from"`
	To     string `json:"to"`
}

// OrderStats untuk statistik order
type OrderStats struct {
	Total       int64 `json:"total"`
	Pending     int64 `json:"pending"`
	Processed   int64 `json:"processed"`
	Completed   int64 `json:"completed"`
	Cancelled   int64 `json:"cancelled"`
	TodayCount  int64 `json:"today_count"`
	WeekCount   int64 `json:"week_count"`
	MonthCount  int64 `json:"month_count"`
}

func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*domain.OrderLog, error) {
	// Validasi
	if req.Name == "" {
		return nil, errors.New("nama wajib diisi")
	}
	if req.Phone == "" {
		return nil, errors.New("nomor WhatsApp wajib diisi")
	}

	order := &domain.OrderLog{
		ID:           uuid.New(),
		VisitorName:  req.Name,
		VisitorWA:    req.Phone,
		ProductName:  req.ProductName,
		ProjectType:  req.ProjectType,
		Budget:       req.Budget,
		Description:  req.Description,
		Status:       domain.OrderStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetAllOrders() ([]domain.OrderLog, error) {
	return s.orderRepo.FindAll()
}

// GetAllOrdersWithPagination - dengan pagination
func (s *OrderService) GetAllOrdersWithPagination(limit, offset int) ([]domain.OrderLog, int64, error) {
	orders, err := s.orderRepo.FindAllWithPagination(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.orderRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetOrdersWithFilters - dengan filter status, search, dan date range
func (s *OrderService) GetOrdersWithFilters(req GetOrdersRequest) ([]domain.OrderLog, int64, error) {
	// Set default limit
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	orders, err := s.orderRepo.FindWithFilters(
		req.Limit,
		req.Offset,
		req.Status,
		req.Search,
		req.From,
		req.To,
	)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.orderRepo.CountWithFilters(req.Status, req.Search, req.From, req.To)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *OrderService) GetOrderByID(id string) (*domain.OrderLog, error) {
	return s.orderRepo.FindByID(id)
}

func (s *OrderService) UpdateOrderStatus(id string, status string) error {
	// Validasi status
	validStatuses := map[string]bool{
		domain.OrderStatusPending:   true,
		domain.OrderStatusProcessed: true,
		domain.OrderStatusCompleted: true,
		domain.OrderStatusCancelled: true,
	}
	if !validStatuses[status] {
		return errors.New("status tidak valid")
	}

	return s.orderRepo.UpdateStatus(id, status)
}

func (s *OrderService) DeleteOrder(id string) error {
	return s.orderRepo.Delete(id)
}

func (s *OrderService) CountOrders() (int64, error) {
	return s.orderRepo.Count()
}

// GetOrderStats - mendapatkan statistik order
func (s *OrderService) GetOrderStats() (*OrderStats, error) {
	total, err := s.orderRepo.Count()
	if err != nil {
		return nil, err
	}

	pending, err := s.orderRepo.CountByStatus(domain.OrderStatusPending)
	if err != nil {
		return nil, err
	}

	processed, err := s.orderRepo.CountByStatus(domain.OrderStatusProcessed)
	if err != nil {
		return nil, err
	}

	completed, err := s.orderRepo.CountByStatus(domain.OrderStatusCompleted)
	if err != nil {
		return nil, err
	}

	cancelled, err := s.orderRepo.CountByStatus(domain.OrderStatusCancelled)
	if err != nil {
		return nil, err
	}

	// Today
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayCount, err := s.orderRepo.CountByDateRange(startOfDay, now)
	if err != nil {
		return nil, err
	}

	// This week
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
	weekCount, err := s.orderRepo.CountByDateRange(startOfWeek, now)
	if err != nil {
		return nil, err
	}

	// This month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthCount, err := s.orderRepo.CountByDateRange(startOfMonth, now)
	if err != nil {
		return nil, err
	}

	return &OrderStats{
		Total:      total,
		Pending:    pending,
		Processed:  processed,
		Completed:  completed,
		Cancelled:  cancelled,
		TodayCount: todayCount,
		WeekCount:  weekCount,
		MonthCount: monthCount,
	}, nil
}

// GetRecentOrders - ambil N order terbaru
func (s *OrderService) GetRecentOrders(limit int) ([]domain.OrderLog, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return s.orderRepo.FindRecent(limit)
}

// UpdateOrder - update order (selain status)
func (s *OrderService) UpdateOrder(id string, req CreateOrderRequest) (*domain.OrderLog, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("order tidak ditemukan")
	}

	// Update fields
	if req.Name != "" {
		order.VisitorName = req.Name
	}
	if req.Phone != "" {
		order.VisitorWA = req.Phone
	}
	if req.ProductName != "" {
		order.ProductName = req.ProductName
	}
	if req.ProjectType != "" {
		order.ProjectType = req.ProjectType
	}
	if req.Budget != "" {
		order.Budget = req.Budget
	}
	if req.Description != "" {
		order.Description = req.Description
	}
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}

// BulkDeleteOrders - hapus banyak order sekaligus
func (s *OrderService) BulkDeleteOrders(ids []string) (int64, error) {
	if len(ids) == 0 {
		return 0, errors.New("tidak ada ID yang dipilih")
	}
	return s.orderRepo.BulkDelete(ids)
}