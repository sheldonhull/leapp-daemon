package custom_validator

import (
	"github.com/go-playground/validator/v10"
	"leapp_daemon/core/aws/region"
)

//
//var awsRegionValidator validator.Func = func(fl validator.FieldLevel) bool {
//  return region.IsRegionValid(fl.Field().String())
//}

func AwsRegionValidator(fl validator.FieldLevel) bool {
	return region.IsRegionValid(fl.Field().String())
}
