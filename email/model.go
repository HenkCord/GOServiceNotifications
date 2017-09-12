package email

import (
	"time"

	"github.com/spf13/viper"
	"github.com/HenkCord/GOServiceNotifications/utils"
	"gopkg.in/mgo.v2/bson"

	"errors"
)

var errInvalidActivationCode = "invalid_argument_activationCode"
var errInvalidEmail = "invalid_argument_email"
var errInvalidUserID = "invalid_argument_userId"
var errInvalidUsername = "invalid_argument_username"
var errInvalidName = "invalid_argument_name"
var errInvalidSubject = "invalid_argument_subject"
var errInvalidMessage = "invalid_argument_message"
var errTemplateExists = "template_exists"
var errInvalidArguments = "invalid_arguments"
var errNotFound = "not_found"
var errTemplateNotFound = "template_not_found"

//EmailTemplate is schema collection
type EmailTemplate struct {
	ID          bson.ObjectId          `json:"_id" bson:"_id,omitempty"`
	Name        string                 `json:"name" bson:"name,omitempty"`
	Description string                 `json:"description" bson:"description,omitempty"`
	Subject     map[string]interface{} `json:"subject" bson:"subject,omitempty"`
	Message     map[string]interface{} `json:"message" bson:"message,omitempty"`
	UpdatedAt   int64                  `json:"updatedAt" bson:"updatedAt,omitempty"`
	CreatedAt   int64                  `json:"createdAt" bson:"createdAt,omitempty"`
}

//InputData is schema input data
type InputData struct {
	Username       string `json:"username" bson:"username,omitempty"`
	UserID         string `json:"userId" bson:"userId,omitempty"`
	Email          string `json:"email" bson:"email,omitempty"`
	ActivationCode string `json:"activationCode" bson:"activationCode,omitempty"`
	Lang           string `json:"lang" bson:"lang,omitempty"`
}

//CreateEmailValidation check required arguments
func CreateEmailValidation(item *EmailTemplate) (err error) {
	t := time.Now().Unix()

	if item.Name == "" {
		return errors.New(errInvalidName)
	}

	if !utils.LangCheckFields(item.Subject) {
		return errors.New(errInvalidSubject)
	}

	if !utils.LangCheckFields(item.Message) {
		return errors.New(errInvalidMessage)
	}

	item.UpdatedAt = t
	item.CreatedAt = t

	item.ID = bson.NewObjectId()

	return
}

//UpdateEmailValidation check required arguments
func UpdateEmailValidation(item *EmailTemplate) (err error) {

	if len(item.Subject) != 0 && !utils.LangCheckFields(item.Subject) {
		return errors.New(errInvalidSubject)
	}

	if len(item.Message) != 0 && !utils.LangCheckFields(item.Message) {
		return errors.New(errInvalidMessage)
	}

	item.UpdatedAt = time.Now().Unix()

	return
}

//ConfirmEmailAddressValidation check required arguments
func ConfirmEmailAddressValidation(item *InputData) (err error) {

	if item.UserID == "" {
		return errors.New(errInvalidUserID)
	}

	item.Email = NormalizeEmail(item.Email)
	if item.Email == "" {
		return errors.New(errInvalidEmail)
	}

	if item.ActivationCode == "" {
		return errors.New(errInvalidActivationCode)
	}

	if item.Username == "" {
		return errors.New(errInvalidUsername)
	}

	if item.Lang == "" {
		item.Lang = viper.GetString("services.email.lang")
	}

	return
}
