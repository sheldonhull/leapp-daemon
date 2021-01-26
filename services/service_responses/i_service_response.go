package service_responses

import (
	"leapp_daemon/controllers/response_dto"
)

type IServiceResponse interface {
	ToDto() response_dto.IResponseDto
}
