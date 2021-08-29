package model

type Status struct {
	Status string `json:"status"`
	Redis  string `json:"redis"`
	Uptime string `json:"uptime"`
}
