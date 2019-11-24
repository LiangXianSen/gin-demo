package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.bj.sensetime.com/SenseGo/debug-kit/errors"
)

func FaultHandler(c *gin.Context) {
	c.Next()

	if !c.Writer.Written() {
		errorG := c.MustGet("error")
		switch err := errorG.(type) {
		case *errors.Error:
			resp := struct {
				RequestID string `json:"request_id"`
				*errors.Error
			}{
				RequestID: c.GetString("request_id"),
				Error:     err,
			}
			if err.IsBadRequest() {
				c.JSON(http.StatusBadRequest, resp)
			} else {
				c.JSON(http.StatusOK, resp)
			}
		default:
			resp := struct {
				RequestID string `json:"request_id"`
				Code      int    `json:"code"`
				Message   string `json:"message"`
			}{
				RequestID: c.GetString("request_id"),
				Code:      errors.InternalErrorCode,
				Message:   errorG.(error).Error(),
			}
			c.JSON(http.StatusInternalServerError, resp)
		}
	}
}
