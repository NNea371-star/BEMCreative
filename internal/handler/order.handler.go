package handler

import (
	"BE/internal/hub"
	"BE/internal/service"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// POST /api/orders/whatsapp - Order via WhatsApp (public)
func (h *OrderHandler) OrderViaWhatsapp(c *fiber.Ctx) error {
	var req service.CreateOrderRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if req.Name == "" || req.Phone == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Nama dan nomor WhatsApp wajib diisi",
		})
	}

	order, err := h.orderService.CreateOrder(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan pesanan: " + err.Error(),
		})
	}

	// Broadcast ke semua admin via WebSocket
	hub.H.BroadcastToAdmins(fiber.Map{
		"type":         "order:new",
		"message":      fmt.Sprintf("📦 Pesanan baru dari %s", req.Name),
		"visitor_name": req.Name,
		"visitor_wa":   req.Phone,
		"product_name": req.ProductName,
		"project_type": req.ProjectType,
		"budget":       req.Budget,
		"created_at":   time.Now().Format(time.RFC3339),
	})

	return c.Status(201).JSON(fiber.Map{
		"message": "Pesanan berhasil dikirim",
		"order":   order,
	})
}

// GET /api/admin/orders - Get all orders (admin only)
func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	status := c.Query("status")
	search := c.Query("search")
	from := c.Query("from")
	to := c.Query("to")

	if status != "" || search != "" || from != "" || to != "" {
		req := service.GetOrdersRequest{
			Limit:  limit,
			Offset: offset,
			Status: status,
			Search: search,
			From:   from,
			To:     to,
		}
		orders, total, err := h.orderService.GetOrdersWithFilters(req)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Gagal mengambil data pesanan",
			})
		}
		return c.JSON(fiber.Map{
			"data":   orders,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		})
	}

	orders, total, err := h.orderService.GetAllOrdersWithPagination(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data pesanan",
		})
	}

	return c.JSON(fiber.Map{
		"data":   orders,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GET /api/admin/orders/stats - Get order statistics
func (h *OrderHandler) GetOrderStats(c *fiber.Ctx) error {
	stats, err := h.orderService.GetOrderStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil statistik pesanan",
		})
	}
	return c.JSON(stats)
}

// GET /api/admin/orders/recent - Get recent orders
func (h *OrderHandler) GetRecentOrders(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	orders, err := h.orderService.GetRecentOrders(limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil pesanan terbaru",
		})
	}
	return c.JSON(orders)
}

// GET /api/admin/orders/:id - Get order by ID
func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID pesanan wajib diisi",
		})
	}

	order, err := h.orderService.GetOrderByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Pesanan tidak ditemukan",
		})
	}

	return c.JSON(order)
}

// PUT /api/admin/orders/:id - Update order
func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID pesanan wajib diisi",
		})
	}

	var req service.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	order, err := h.orderService.UpdateOrder(id, req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Broadcast update ke admin lain
	hub.H.BroadcastToAdmins(fiber.Map{
		"type":  "order:updated",
		"order": order,
	})

	return c.JSON(fiber.Map{
		"message": "Pesanan berhasil diupdate",
		"order":   order,
	})
}

// PUT /api/admin/orders/:id/status - Update order status
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID pesanan wajib diisi",
		})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if req.Status == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Status wajib diisi",
		})
	}

	if err := h.orderService.UpdateOrderStatus(id, req.Status); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Broadcast status baru ke admin lain
	hub.H.BroadcastToAdmins(fiber.Map{
		"type":   "order:status_updated",
		"id":     id,
		"status": req.Status,
	})

	return c.JSON(fiber.Map{
		"message": "Status pesanan berhasil diupdate",
	})
}

// DELETE /api/admin/orders/:id - Delete order
func (h *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID pesanan wajib diisi",
		})
	}

	if err := h.orderService.DeleteOrder(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menghapus pesanan",
		})
	}

	// Broadcast penghapusan ke admin lain
	hub.H.BroadcastToAdmins(fiber.Map{
		"type": "order:deleted",
		"id":   id,
	})

	return c.JSON(fiber.Map{
		"message": "Pesanan berhasil dihapus",
	})
}

// POST /api/admin/orders/bulk-delete - Bulk delete orders
func (h *OrderHandler) BulkDeleteOrders(c *fiber.Ctx) error {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Tidak ada ID yang dipilih",
		})
	}

	deleted, err := h.orderService.BulkDeleteOrders(req.IDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Broadcast bulk delete ke admin lain
	hub.H.BroadcastToAdmins(fiber.Map{
		"type": "order:bulk_deleted",
		"ids":  req.IDs,
	})

	return c.JSON(fiber.Map{
		"message":       fmt.Sprintf("%d pesanan berhasil dihapus", deleted),
		"deleted_count": deleted,
	})
}