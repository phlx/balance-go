package validations

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"balance/internal/exchangerates"
)

func RegisterCurrency(log *logrus.Logger) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
			currency, ok := fl.Field().Interface().(string)
			if ok {
				_, err := exchangerates.StringToCurrency(currency)
				return err == nil
			}
			return false
		})
		if err != nil {
			log.Fatal("Unable to register validation")
		}
	}
}
