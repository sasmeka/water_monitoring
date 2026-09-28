package pkg

import (
	"github.com/sasmeka/water_monitoring/config"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int         `json:"status"`
	Message interface{} `json:"pesan,omitempty"`
	Token   interface{} `json:"token,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Step    interface{} `json:"step,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func (r *Response) Send(ctx *gin.Context) {
	ctx.JSON(r.Status, r)
	ctx.Abort()
	return
}

func Responses(code int, data *config.Result) *Response {
	var response = Response{
		Status: code,
	}
	response.Message = getStatus(code)
	if data.Message != nil {
		response.Message = data.Message
	}
	if data.Data != nil {
		response.Data = data.Data
	}
	if data.Step != nil {
		response.Step = data.Step
	}
	if data.Meta != nil {
		response.Meta = data.Meta
	}
	if data.Token != nil {
		response.Token = data.Token
	}

	return &response
}

func getStatus(status int) string {
	var desc string
	switch status {
	case 200:
		desc = "OK"
		break
	case 400:
		desc = "Bad Request"
		break
	case 401:
		desc = "Unauthorized"
		break
	case 403:
		desc = "Forbidden"
		break
	case 404:
		desc = "Not Found"
		break
	case 500:
		desc = "Internal Server Error"
		break
	case 501:
		desc = "Bad Gateway"
		break
	case 304:
		desc = "Not Modified"
		break
	default:
		desc = ""
	}

	return desc
}
