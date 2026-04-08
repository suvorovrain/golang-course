package dto

import (
	"repo-stat/api/internal/domain"
)

type ServicesInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResponse struct {
	Status   string         `json:"status"`
	Services []ServicesInfo `json:"services"`
}

func CreatePingResponce(system_status string, processor_status domain.PingStatus, subscriber_status domain.PingStatus) *PingResponse {

	return &PingResponse{
		Status: system_status,
		Services: []ServicesInfo{
			ServicesInfo{
				Name:   "processor",
				Status: string(processor_status),
			},
			ServicesInfo{
				Name:   "subscriber",
				Status: string(subscriber_status),
			},
		},
	}
}
