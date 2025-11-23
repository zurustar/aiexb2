// backend/internal/util/validator.go
package util

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator はバリデーション機能を提供する構造体
type Validator struct {
	validate *validator.Validate
}

// NewValidator は新しいValidatorインスタンスを作成します
func NewValidator() *Validator {
	v := validator.New()

	// JSONタグ名をフィールド名として使用する設定
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{validate: v}
}

// ValidateStruct は構造体のバリデーションを行います
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

// FormatError はバリデーションエラーを読みやすいメッセージに変換します
// 簡易的な実装です。必要に応じて国際化対応などを追加できます。
func FormatError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			tag := fieldError.Tag()
			param := fieldError.Param()

			switch tag {
			case "required":
				errors[field] = "This field is required"
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				errors[field] = fmt.Sprintf("Must be at least %s characters long", param)
			case "max":
				errors[field] = fmt.Sprintf("Must be at most %s characters long", param)
			case "uuid":
				errors[field] = "Invalid UUID format"
			default:
				errors[field] = fmt.Sprintf("Failed validation on tag: %s", tag)
			}
		}
	} else {
		errors["error"] = err.Error()
	}

	return errors
}
