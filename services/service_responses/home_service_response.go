package service_responses

import "leapp_daemon/controllers/response_dto"

type HomeServiceResponse struct {
	Data string
}

func (serviceResponse *HomeServiceResponse) ToDto() *response_dto.HomeResponseDto {
	return &response_dto.HomeResponseDto{
		Message: "success",
		Data: serviceResponse.Data,
	}
}
