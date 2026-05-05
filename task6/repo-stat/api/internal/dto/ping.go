package dto

type PingResponse struct {
	Status   string          `json:"status"`
	Services []ServiceStatus `json:"services"`
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
