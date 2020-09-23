package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"balance/internal/utils"
)

type ValidationErrorWithResponse struct {
	Response ErrorResponse
}

func (v *ValidationErrorWithResponse) Error() string {
	result, err := json.Marshal(v.Response)
	if err != nil {
		return fmt.Sprintf("Unable to marshal validation response: [%+v]", v.Response)
	}
	return "Validation error: " + string(result)
}

func Validate(context *gin.Context, request interface{}, binding binding.Binding) *ValidationErrorWithResponse {
	if err := context.ShouldBindWith(request, binding); err != nil {
		var errors []ErrorResponseItem

		fieldErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return &ValidationErrorWithResponse{
				Response: ErrorResponse{
					Errors: append(errors, ErrorResponseItem{
						Code:   fmt.Sprintf(ErrorCodeValidation, "unknown"),
						Detail: fmt.Sprintf("Validation failed cause of: %+v'", err),
					}),
				},
			}
		}

		for _, fieldError := range fieldErrors {
			field := utils.ToSnakeCase(fieldError.Field())
			rule := fieldError.Tag()

			errors = append(errors, ErrorResponseItem{
				Code:   fmt.Sprintf(ErrorCodeValidation, field),
				Detail: fmt.Sprintf("Validation for field '%s' failed rule '%s'", field, rule),
			})
		}

		return &ValidationErrorWithResponse{
			Response: ErrorResponse{
				Errors: errors,
			},
		}
	}

	return nil
}
