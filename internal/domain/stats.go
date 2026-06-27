package domain

type SiteStats struct {
	MachineTypes   int `json:"machine_types"`
	ExperienceYears int `json:"experience_years"`
	MachinesBuilt  int `json:"machines_built"`
	WorksProduced  int `json:"works_produced"`
}