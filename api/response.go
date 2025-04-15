package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	PerPage      int `json:"per_page"`
	TotalRecords int `json:"total_records"`
}

func RespondWithSuccess(c *gin.Context, data any, message string, statusCode ...int) {
	resp := Response{
		Success: true,
		Data:    data,
	}

	if message != "" {
		resp.Message = message
	}

	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	c.JSON(code, resp)
}

func RespondWithPagination(c *gin.Context, data any, meta PaginationMeta) {
	resp := Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
	c.JSON(http.StatusOK, resp)
}

func RespondWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
