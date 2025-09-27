package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

// ErrorResponse represents a structured error message.
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

// ErrorHandlerMiddleware is a global error handler that recovers from panics
// and formats error responses consistently.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Interface("error", err).Msg("Panic recovered")
				c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "An unexpected internal server error occurred."})
			}
		}()
		c.Next()

		// Capture any errors that were not handled by other middleware/handlers
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Check for validation errors
				if validationErrs, ok := e.Err.(validator.ValidationErrors); ok {
					errs := make(map[string]string)
					for _, fe := range validationErrs {
						errs[fe.Field()] = "Validation failed on tag: " + fe.Tag()
					}
					c.JSON(http.StatusBadRequest, ErrorResponse{
						Error:   "Invalid request payload",
						Details: errs,
					})
					return // Stop processing after handling validation errors
				}
			}
			// Fallback for other generic errors if not already handled
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "An internal server error occurred."})
		}
	}
}