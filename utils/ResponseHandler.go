package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Response struct {
	Code   int          `json:"code"`
	Status string       `json:"status"`
	Data   interface{}  `json:"data"`
	Meta   *interface{} `json:"meta"`
	Error  interface{}  `json:"error,omitempty"`
}
type Pagination struct {
	Page      int `json:"page"`
	Size      int `json:"size"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

func HandleResponse(c *gin.Context, code int, status string, data interface{}, meta interface{}, err error) {
	if err != nil {
		HandleError(c, err, code)
	}
	c.JSON(code, Response{
		Code:   code,
		Status: status,
		Data:   data,
		Meta:   &meta,
	})

}

// HandleError is a simplified error response handler
func HandleError(c *gin.Context, err error, code int) {
	log.Printf("Error: %v", err) // Simple logging for traceability
	c.JSON(http.StatusInternalServerError, Response{
		Status: "error",
		Error:  err.Error(),
		Code:   code,
	})
}
