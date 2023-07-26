package common

import (
	"github.com/gin-gonic/gin"
)

type GenericResponseSerializer struct {
	C *gin.Context
}

type GenericResponse struct {
	RequestID     string       `json:"request_id"`
	LeaseID       string       `json:"lease_id"`
	LeaseDuration int          `json:"lease_duration"`
	Renewable     bool         `json:"renewable"`
	Warnings      *[]string    `json:"warnings"`
	Data          interface{}  `json:"data"`
	WrapInfo      *interface{} `json:"wrap_info"`
	Auth          *interface{} `json:"auth"`
}

func (s *GenericResponseSerializer) Response(data interface{}) GenericResponse {
	requestID := s.C.MustGet("request_id").(string)
	response := GenericResponse{
		RequestID:     requestID,
		LeaseID:       "",
		LeaseDuration: 0,
		Renewable:     false,
		Warnings:      nil,
		Data:          data,
	}
	return response
}

func NewGenericResponse(c *gin.Context, data interface{}) GenericResponse {
	gr := GenericResponseSerializer{C: c}
	return gr.Response(data)
}

type StatusResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewStatusResponse(code int, status string) StatusResponse {
	response := StatusResponse{
		Code:   code,
		Status: status,
	}
	return response
}
