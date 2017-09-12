package email

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/HenkCord/notifications/db"
	"github.com/HenkCord/notifications/errors"
	"github.com/HenkCord/notifications/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ConfirmEmailAddress controller
// Method: POST
// URL: /email/confirmEmail
// Header:
// Content-type
// Params:
// Query:
// Body:
//	email 			string
// 	username     	string
// 	activationCode  string
//	userId			string
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
func ConfirmEmailAddress(c *gin.Context) {

	var item InputData

	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(errors.BadRequest(errInvalidArguments, err.Error()))
		return
	}

	err = ConfirmEmailAddressValidation(&item)
	if err != nil {
		c.JSON(errors.BadRequest(err.Error(), ""))
		return
	}

	templateName := "confirmEmail"

	//Get Template
	collection := db.EmailTemplatesCollection()
	tpl := EmailTemplate{}
	err = collection.Find(bson.M{"name": templateName}).One(&tpl)
	if err != nil {
		//Ошибка может быть как из за отсутствия в бд записи, так и в случае отсутствия подключения к бд
		if err.Error() != "not found" {
			panic(err.Error())
		}
		c.JSON(errors.NotFound(errTemplateNotFound, err.Error()))
		return
	}

	templateData := struct {
		Username             string
		URLActivationAccount string
		URLFeedback          string
		SupportEmailAddress  string
	}{
		Username:             item.Username,
		URLActivationAccount: viper.GetString("services.email.urlActivationAccount") + item.ActivationCode,
		URLFeedback:          viper.GetString("services.services.all.urlFeedback"),
		SupportEmailAddress:  viper.GetString("services.services.all.supportEmail"),
	}

	res := NewEmail(item.Email, utils.LangGetField(tpl.Subject, item.Lang), "")
	err = res.ParseTemplate(templateName, utils.LangGetField(tpl.Message, item.Lang), templateData)
	if err != nil {
		panic(err.Error())
	}

	success, err := res.Send()
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"success": success})
}
