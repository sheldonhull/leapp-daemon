package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/request_dto/g_suite_auth_dto"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/logging"
	"leapp_daemon/service"
	"net/http"
	"net/url"
	"strings"
)

type GSuiteAuthFirstStepResponse struct {
	CaptchaForm       url.Values `json:"captchaForm"`
	CaptchaInputId    string `json:"captchaInputId"`
	CaptchaPictureURL string `json:"captchaPictureURL"`
	CaptchaURL        string `json:"captchaURL"`
	LoginForm         url.Values `json:"loginForm"`
	LoginURL          string `json:"loginURL"`
}

type GSuiteAuthSecondStepResponse struct {
	IsMfaTokenRequested bool `json:"isMfaTokenRequested"`
	ResponseForm url.Values `json:"responseForm"`
	SubmitURL string `json:"submitURL"`
}

func GSuiteAuthFirstStepController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := g_suite_auth_dto.GSuiteAuthFirstStepRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	captchaForm, captchaInputId, captchaPictureURL, captchaURL, loginForm, loginURL := service.GSuiteAuthFirstStepService(requestDto.Username, requestDto.Password)

	gSuiteAuthFirstStepResponse := GSuiteAuthFirstStepResponse{
		CaptchaForm:       captchaForm,
		CaptchaInputId:    captchaInputId,
		CaptchaPictureURL: captchaPictureURL,
		CaptchaURL:        captchaURL,
		LoginForm:         loginForm,
		LoginURL:          loginURL,
	}

	marshalledGSuiteAuthFirstStepResponse, _ := json.Marshal(&gSuiteAuthFirstStepResponse)
	marshalledGSuiteAuthFirstStepResponseString := string(marshalledGSuiteAuthFirstStepResponse)
	marshalledGSuiteAuthFirstStepResponseString = strings.Replace(marshalledGSuiteAuthFirstStepResponseString, "\\u003c", "<", -1)
	marshalledGSuiteAuthFirstStepResponseString = strings.Replace(marshalledGSuiteAuthFirstStepResponseString, "\\u003e", ">", -1)
	marshalledGSuiteAuthFirstStepResponseString = strings.Replace(marshalledGSuiteAuthFirstStepResponseString, "\\u0026", "&", -1)

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: marshalledGSuiteAuthFirstStepResponseString}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func GSuiteAuthSecondStepController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := g_suite_auth_dto.GSuiteAuthSecondStepRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	var captchaForm url.Values
	captchaFormString := []byte(requestDto.CaptchaForm)
	_ = json.Unmarshal(captchaFormString, &captchaForm)

	var loginForm url.Values
	loginFormString := []byte(requestDto.LoginForm)
	_ = json.Unmarshal(loginFormString, &loginForm)

	isMfaTokenRequested, responseForm, submitURL := service.GSuiteAuthSecondStepService(requestDto.Captcha, requestDto.CaptchaInputId, requestDto.CaptchaUrl,
		captchaForm, requestDto.Password, loginForm, requestDto.LoginUrl)

	gSuiteAuthSecondStepResponse := GSuiteAuthSecondStepResponse{
		IsMfaTokenRequested: isMfaTokenRequested,
		ResponseForm: responseForm,
		SubmitURL: submitURL,
	}

	marshalledGSuiteAuthSecondStepResponse, _ := json.Marshal(&gSuiteAuthSecondStepResponse)

	marshalledGSuiteAuthSecondStepResponseString := string(marshalledGSuiteAuthSecondStepResponse)
	marshalledGSuiteAuthSecondStepResponseString = strings.Replace(marshalledGSuiteAuthSecondStepResponseString, "\\u003c", "<", -1)
	marshalledGSuiteAuthSecondStepResponseString = strings.Replace(marshalledGSuiteAuthSecondStepResponseString, "\\u003e", ">", -1)
	marshalledGSuiteAuthSecondStepResponseString = strings.Replace(marshalledGSuiteAuthSecondStepResponseString, "\\u0026", "&", -1)

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: marshalledGSuiteAuthSecondStepResponseString}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func GSuiteAuthThirdStepController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := g_suite_auth_dto.GSuiteAuthThirdStepRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	var responseForm url.Values
	loginFormString := []byte(requestDto.ResponseForm)
	_ = json.Unmarshal(loginFormString, &responseForm)

	samlAssertion := service.GSuiteAuthThirdStepService(requestDto.IsMfaTokenRequested, responseForm,
		requestDto.SubmitURL, requestDto.Token)

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: samlAssertion}
	context.JSON(http.StatusOK, responseDto.ToMap())
}