package controller

/*
// swagger:response getConfigurationResponse
type getConfigurationResponseWrapper struct {
	// in: body
	Body getConfigurationResponse
}

type getConfigurationResponse struct {
	Message string
	Data    configuration.Configuration
}

func CreateConfigurationController(context *gin.Context) {
	// swagger:route POST /configuration configuration createConfiguration
	// Create a new configuration
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	service := configuration.facade{
		ConfigurationRepository: &repository.FileConfigurationRepository{
			FileSystem: &file_system.FileSystem{},
			Encryption: &encryption.Encryption{},
		},
	}

	err := service.Create()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func ReadConfigurationController(context *gin.Context) {
	// swagger:route GET /configuration configuration getConfiguration
	// Create a new configuration
	//   Responses:
	//     200: getConfigurationResponse

	logging.SetContext(context)

	service := configuration.facade{
		ConfigurationRepository: &repository.FileConfigurationRepository{
			FileSystem: &file_system.FileSystem{},
			Encryption: &encryption.Encryption{},
		},
	}

	config, err := service.Get()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: config}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
*/
