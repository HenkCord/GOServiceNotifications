package push

import (
	"time"

	"github.com/HenkCord/notifications/utils"
	"gopkg.in/mgo.v2/bson"

	"errors"
)

var errInvalidUserID = "invalid_argument_userId"
var errInvalidPlayerID = "invalid_argument_playerId"
var errInvalidName = "invalid_argument_name"
var errInvalidTitle = "invalid_argument_title"
var errInvalidMessage = "invalid_argument_message"
var errTemplateExists = "template_exists"
var errInvalidArguments = "invalid_arguments"
var errRemoteServiceSentError = "remote_service_sent_error"
var errNotFound = "not_found"
var errTemplateNotFound = "template_not_found"

//PushTemplate is schema collection
type PushTemplate struct {
	ID          bson.ObjectId          `json:"_id" bson:"_id,omitempty"`
	Name        string                 `json:"name" bson:"name,omitempty"`
	Description string                 `json:"description" bson:"description,omitempty"`
	Title       map[string]interface{} `json:"title" bson:"title,omitempty"`
	Message     map[string]interface{} `json:"message" bson:"message,omitempty"`
	UpdatedAt   int64                  `json:"updatedAt" bson:"updatedAt,omitempty"`
	CreatedAt   int64                  `json:"createdAt" bson:"createdAt,omitempty"`
}

//InputData is schema input data
type InputData struct {
	UserID   string `json:"userId" bson:"userId,omitempty"`
	PlayerID string `json:"playerId" bson:"playerId,omitempty"`
	Lang     string `json:"lang" bson:"lang,omitempty"`
}

//CreatePushValidation check required arguments
func CreatePushValidation(item *PushTemplate) (err error) {
	t := time.Now().Unix()

	if item.Name == "" {
		return errors.New(errInvalidName)
	}

	// For IOS 10
	// if !utils.LangCheckFields(item.Subject) {
	// 	return errors.New(errInvalidSubject)
	// }

	if len(item.Title) != 0 && !utils.LangCheckFields(item.Title) {
		return errors.New(errInvalidTitle)
	}

	if !utils.LangCheckFields(item.Message) {
		return errors.New(errInvalidMessage)
	}

	item.UpdatedAt = t
	item.CreatedAt = t

	item.ID = bson.NewObjectId()

	return
}

//UpdatePushValidation check required arguments
func UpdatePushValidation(item *PushTemplate) (err error) {

	// For IOS 10
	// if len(item.Subject) != 0 && !utils.LangCheckFields(item.Subject) {
	// 	return errors.New(errInvalidSubject)
	// }

	if len(item.Title) != 0 && !utils.LangCheckFields(item.Title) {
		return errors.New(errInvalidTitle)
	}

	if len(item.Message) != 0 && !utils.LangCheckFields(item.Message) {
		return errors.New(errInvalidMessage)
	}

	item.UpdatedAt = time.Now().Unix()

	return
}

//GiveReviewValidation check required arguments
func GiveReviewValidation(item *InputData) (err error) {

	if item.UserID == "" {
		return errors.New(errInvalidUserID)
	}

	if item.PlayerID == "" {
		return errors.New(errInvalidPlayerID)
	}

	// if item.Lang == "" {
	// 	item.Lang = viper.GetString("services.push.lang")
	// }

	return
}
