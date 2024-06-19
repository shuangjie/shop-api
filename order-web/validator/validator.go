package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 1 开头，第二位是 3,4,5,6,7,8,9，后面 9 位随意
	result, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	return result
}
