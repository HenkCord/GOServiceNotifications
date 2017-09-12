package sms

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/HenkCord/notifications/db"
	"github.com/HenkCord/notifications/errors"
	"github.com/HenkCord/notifications/utils"
	"github.com/gin-gonic/gin"
)

// GetCode controller
// Method: POST
// URL: /sms/getCode
// Header:
// Content-type
// Params:
// Query:
// Body:
// 	phone     		string
//	userId			string
// 	code   			string
//  [lang]			string
// Return:
//  200 "success"				true
// Errors:
//  500 "internal_server_error" "error_server"
//  400 "bad_request" 			"invalid_arguments"
//  400 "bad_request" 			"invalid_argument_phone"
//  400 "bad_request" 			"invalid_argument_userId"
//  400 "bad_request" 			"invalid_argument_code"
//  400 "bad_request" 			"invalid_number"
//  404 "not_found" 			"template_not_found"
func GetCode(c *gin.Context) {

	var item InputData

	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(errors.BadRequest(errInvalidArguments, err.Error()))
		return
	}

	err = GetCodeValidation(&item)
	if err != nil {
		c.JSON(errors.BadRequest(err.Error(), ""))
		return
	}

	templateData := struct {
		Code string
	}{
		Code: item.Code,
	}

	templateName := "getCode"

	//Get Template
	collection := db.SMSTemplatesCollection()
	tpl := SMSTemplate{}
	err = collection.Find(bson.M{"name": templateName}).One(&tpl)
	if err != nil {
		//Ошибка может быть как из за отсутствия в бд записи, так и в случае отсутствия подключения к бд
		if err.Error() != "not found" {
			panic(err.Error())
		}
		c.JSON(errors.NotFound(errTemplateNotFound, err.Error()))
		return
	}

	res := NewSMS(item.Phone)
	err = res.ParseTemplate(templateName, utils.LangGetField(tpl.Message, item.Lang), templateData)
	if err != nil {
		panic(err.Error())
	}

	success, err := res.Send()
	if !success {
		//Удаленный sms сервис не отвечает, либо прислал критическую ошибку.
		c.JSON(errors.BadRequest(errRemoteServiceSentError, err.Error()))
		return
	}
	//Удаленный sms сервис ответил с ошибкой о некорректном номере телефона
	if err != nil {
		c.JSON(errors.BadRequest(errInvalidNumber, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
