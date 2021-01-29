package error_handling

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"leapp_daemon/rest_api/controllers/utils"
	"log"
	"net/http"
)

type errorHandler struct{}

func (*errorHandler) Handle(context *gin.Context, err error) {
	var code int

	switch err.(type) {
	case BadRequestError:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}

	errorMap := gin.H{ "statusCode": code, "error": err.Error(), "context": utils.NewContext(context) }
	errorJson, marshallingError := json.MarshalIndent(errorMap, "", "  ")

	if marshallingError != nil {
		log.Println(fmt.Sprintf("%+v", errorMap))
		context.JSON(code, errorMap)
	} else {
		log.Println(fmt.Sprintf("%+v", string(errorJson)))
		context.JSON(code, errorMap)
	}
}

var ErrorHandler = &errorHandler{}
