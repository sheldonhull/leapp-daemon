package services

import (
	"fmt"
	"leapp_daemon/services/service_requests"
	"leapp_daemon/services/service_responses"
)

func HomeService(serviceRequest service_requests.HomeServiceRequest) (*service_responses.HomeServiceResponse, error) {
	res := service_responses.HomeServiceResponse{
		Data: fmt.Sprintf("hello %s", serviceRequest.Name),
	}
	return &res, nil
}
