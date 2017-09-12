package push

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/HenkCord/GOServiceNotifications/db"
	"github.com/HenkCord/GOServiceNotifications/errors"
	"github.com/HenkCord/GOServiceNotifications/utils"
	"github.com/gin-gonic/gin"
)

// GiveReview - отправка пуш уведомления с просьбой оставить
// отзыв по OneSignal playerId пользователя
// Method: POST
// URL: /push/giveReview
// Header:
// Content-type
// Params:
// Query:
// Body:
//	userId			string
// 	playerId   		string
//  [lang]			string
// Return:
//  200 "success"				true
// Errors:
//  500 "internal_server_error" "error_server"
//  400 "bad_request" 			"invalid_arguments"
//  400 "bad_request" 			"invalid_argument_userId"
//  400 "bad_request" 			"invalid_argument_playerId"
//  404 "not_found" 			"template_not_found"
func GiveReview(c *gin.Context) {

	var item InputData

	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(errors.BadRequest(errInvalidArguments, err.Error()))
		return
	}

	err = GiveReviewValidation(&item)
	if err != nil {
		c.JSON(errors.BadRequest(err.Error(), ""))
		return
	}

	templateName := "giveReview"

	//Get Template
	collection := db.PushTemplatesCollection()
	tpl := PushTemplate{}
	err = collection.Find(bson.M{"name": templateName}).One(&tpl)
	if err != nil {
		//Ошибка может быть как из за отсутствия в бд записи, так и в случае отсутствия подключения к бд
		if err.Error() != "not found" {
			panic(err.Error())
		}
		c.JSON(errors.NotFound(errTemplateNotFound, err.Error()))
		return
	}

	templateData := struct{}{}

	//BEGIN Special for OneSignal
	title, err := utils.MapInterfaceToString(tpl.Title)
	if err != nil {
		panic(err.Error())
	}
	message, err := utils.MapInterfaceToString(tpl.Message)
	if err != nil {
		panic(err.Error())
	}
	//END Special for OneSignal

	res := NewPush(item.PlayerID, title, "")
	err = res.ParseTemplate(templateName, message, templateData)
	if err != nil {
		panic(err.Error())
	}

	success, err := res.Send()
	if err != nil {
		c.JSON(errors.BadRequest(errRemoteServiceSentError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": success})
}
