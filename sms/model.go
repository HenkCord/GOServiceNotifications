package sms

import (
	"time"

	"github.com/HenkCord/notifications/utils"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"

	"errors"
)

var errInvalidPhone = "invalid_argument_phone"
var errInvalidUserID = "invalid_argument_userId"
var errInvalidCode = "invalid_argument_code"
var errInvalidMessage = "invalid_argument_message"
var errInvalidName = "invalid_argument_name"
var errTemplateExists = "template_exists"
var errInvalidArguments = "invalid_arguments"
var errNotFound = "not_found"
var errTemplateNotFound = "template_not_found"
var errInvalidNumber = "invalid_number"
var errRemoteServiceSentError = "remote_service_sent_error"
var errInvalidPlaceName = "invalid_placeName"
var errInvalidPlaceType = "invalid_placeType"
var errInvalidPlaceStreet = "invalid_placeStreet"
var errInvalidNearWorkTime = "invalid_nearWorkTime"
var errInvalidDate = "invalid_date"
var errInvalidTime = "invalid_time"
var errInvalidPeopleCount = "invalid_peopleCount"
var errInvalidFirstName = "invalid_firstName"

//SMS is schema collection
// type SMS struct {
// 	ID        bson.ObjectId `json:"_id" bson:"_id,omitempty"`
// 	Phone     string        `json:"phone" bson:"phone,omitempty"`
// 	UserID    string        `json:"userId" bson:"userId,omitempty"`
// 	Message   string        `json:"message" bson:"message,omitempty"`
// 	Status    string        `json:"status" bson:"status,omitempty"`
// 	SMSID     string        `json:"smsId" bson:"smsId,omitempty"`
// 	UpdatedAt int64         `json:"updatedAt" bson:"updatedAt,omitempty"`
// 	CreatedAt int64         `json:"createdAt" bson:"createdAt,omitempty"`
// }

//EmailTemplate is schema collection
type SMSTemplate struct {
	ID          bson.ObjectId          `json:"_id" bson:"_id,omitempty"`
	Name        string                 `json:"name" bson:"name,omitempty"`
	Description string                 `json:"description" bson:"description,omitempty"`
	Message     map[string]interface{} `json:"message" bson:"message,omitempty"`
	UpdatedAt   int64                  `json:"updatedAt" bson:"updatedAt,omitempty"`
	CreatedAt   int64                  `json:"createdAt" bson:"createdAt,omitempty"`
}

//InputData is schema input data
type InputData struct {
	Phone        string `json:"phone" bson:"phone,omitempty"`
	UserID       string `json:"userId" bson:"userId,omitempty"`
	Code         string `json:"code" bson:"code,omitempty"`
	PlaceType    string `json:"placeType" bson:"placeType,omitempty"`
	PlaceName    string `json:"placeName" bson:"placeName,omitempty"`
	PlaceStreet  string `json:"placeStreet" bson:"placeStreet,omitempty"`
	NearWorkTime string `json:"nearWorkTime" bson:"nearWorkTime,omitempty"`
	SMSID        string `json:"smsId" bson:"smsId,omitempty"`
	Date         string `json:"date" bson:"time,omitempty"`
	Time         string `json:"time" bson:"time,omitempty"`
	PeopleCount  string `json:"peopleCount" bson:"peopleCount,omitempty"`
	FirstName    string `json:"firstName" bson:"firstName,omitempty"`
	Lang         string `json:"lang" bson:"lang,omitempty"`
}

//CreateSMSValidation check required arguments
func CreateSMSValidation(item *SMSTemplate) (err error) {
	t := time.Now().Unix()

	if item.Name == "" {
		return errors.New(errInvalidName)
	}

	if !utils.LangCheckFields(item.Message) {
		return errors.New(errInvalidMessage)
	}

	item.UpdatedAt = t
	item.CreatedAt = t

	item.ID = bson.NewObjectId()

	return
}

//UpdateSMSValidation check required arguments
func UpdateSMSValidation(item *SMSTemplate) (err error) {

	if len(item.Message) != 0 && !utils.LangCheckFields(item.Message) {
		return errors.New(errInvalidMessage)
	}

	item.UpdatedAt = time.Now().Unix()

	return
}

//GetCodeValidation check required arguments
func GetCodeValidation(item *InputData) (err error) {

	if item.UserID == "" {
		return errors.New(errInvalidUserID)
	}

	if item.Code == "" {
		return errors.New(errInvalidCode)
	}

	if item.Phone == "" {
		return errors.New(errInvalidPhone)
	}

	if PhoneNormalize(item.Phone) != nil {
		return errors.New(errInvalidPhone)
	}

	if item.Lang == "" {
		item.Lang = viper.GetString("services.sms.lang")
	}

	return
}
