package httpapi

import "github.com/gin-gonic/gin"

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func apiError(c *gin.Context, status int, code, msg string, fields map[string]string) {
	c.AbortWithStatusJSON(status, gin.H{"error": ErrorBody{Code: code, Message: msg, Fields: fields}})
}
