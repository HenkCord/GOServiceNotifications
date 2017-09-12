package email

import (
	"net/http"
	"strconv"

	"github.com/HenkCord/notifications/db"
	"github.com/HenkCord/notifications/errors"
	"github.com/HenkCord/notifications/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// GetTemplates controller
// Method: GET
// URL: /email/template
// Header:
// Content-type
// Params:
// Query: [limit int], [skip int]
// Body:
// Return:
//	list						array
//  	template 				object
//		template._id			string
//		template.name			string
//  	template.description	string
//		template.subject 		object
//		template.message 		object
//		template.updatedAt		int64
//		template.createdAt		int64
// Errors:
//  500 "internal_server_error" "error_server"
//  404 "not_found" "not_found"
func GetTemplates(c *gin.Context) {
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))

	collection := db.EmailTemplatesCollection()

	resultItem := []EmailTemplate{}

	err := collection.Find(bson.M{}).Skip(skip).Limit(limit).All(&resultItem)
	if err != nil {
		//Ошибка может быть как из за отсутствия в бд записи, так и в случае отсутствия подключения к бд
		if err.Error() != "not found" {
			panic(err.Error())
		}
		c.JSON(errors.NotFound(errNotFound, err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": resultItem})
}

// GetTemplate controller
// Method: GET
// URL: /email/template/:name
// Header:
// Content-type
// Params:
//	name 		string
// Query:
// Body:
// Return:
//  template 					object
//	template._id				string
//  template.description		string
//	template.name				string
//	template.subject 			object
//	template.message 			object
//	template.updatedAt			int64
//	template.createdAt			int64
// Errors:
//  500 "internal_server_error" "error_server"
//  404 "not_found" "template_not_found"
func GetTemplate(c *gin.Context) {
	name := c.Param("name")

	collection := db.EmailTemplatesCollection()

	resultItem := EmailTemplate{}

	err := collection.Find(bson.M{"name": name}).One(&resultItem)
	if err != nil {
		//Ошибка может быть как из за отсутствия в бд записи, так и в случае отсутствия подключения к бд
		if err.Error() != "not found" {
			panic(err.Error())
		}
		c.JSON(errors.NotFound(errTemplateNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, resultItem)
}

// CreateTemplate controller
// Method: POST
// URL: /email/template
// Header:
// Content-type
// Params:
// Query:
// Body:
// 	name     		string
//  description		string
// 	subject   		object
//	subject.ru		string
//	subject....		string
// 	message  		object
//	message.ru		string
//	message....		string
// Return:
//  [lastInsertId] 	string
//	[error]			string
//  success 		bool
// Errors:
//  500 "internal_server_error" "error_server"
//  400 "bad_request" 			"invalid_arguments"
//  400 "bad_request" 			"invalid_argument_title"
//  400 "bad_request" 			"invalid_argument_message"
func CreateTemplate(c *gin.Context) {

	var item EmailTemplate

	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(errors.BadRequest(errInvalidArguments, err.Error()))
		return
	}

	err = CreateEmailValidation(&item)
	if err != nil {
		c.JSON(errors.BadRequest(err.Error(), ""))
		return
	}

	collection := db.EmailTemplatesCollection()

	resultItem := EmailTemplate{}

	//Проверка существования записи с указанным полем "name"
	err = collection.Find(bson.M{"name": item.Name}).One(&resultItem)
	if err != nil {
		//Ошибка может быть как из за отсутствия в бд записи, так и в случае подключения к бд
		if err.Error() != "not found" {
			panic(err.Error())
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": errTemplateExists})
		return
	}

	err = collection.Insert(item)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "lastInsertId": item.ID})
}

// UpdateTemplate controller
// Method: PUT
// URL: /email/template/:name
// Header:
// Content-type
// Params:
//	name 		string
// Query:
// Body:
//  [description]	string
// 	[subject]   	object
//	subject.ru		string
//	subject....		string
// 	[message]  		object
//	message.ru		string
//	message....		string
// Return:
//  success 		bool
//	updateId		string
// Errors:
//  500 "internal_server_error" "error_server"
//  400 "bad_request" 			"invalid_arguments"
//  400 "bad_request" 			"invalid_argument_subject"
//  400 "bad_request" 			"invalid_argument_message"
func UpdateTemplate(c *gin.Context) {

	name := c.Param("name")

	var item EmailTemplate

	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(errors.BadRequest(errInvalidArguments, err.Error()))
		return
	}

	err = UpdateEmailValidation(&item)
	if err != nil {
		c.JSON(errors.BadRequest(err.Error(), ""))
		return
	}

	collection := db.EmailTemplatesCollection()

	resultItem := EmailTemplate{}

	err = collection.Find(bson.M{"name": name}).One(&resultItem)
	if err != nil {
		//Ошибка может быть как из за отсутствия в бд записи, так и в случае подключения к бд
		if err.Error() != "not found" {
			panic(err.Error())
		}
	} else {
		//merge
		item.Subject = utils.MergeMaps(resultItem.Subject, item.Subject)
		item.Message = utils.MergeMaps(resultItem.Message, item.Message)
	}
	err = collection.Update(bson.M{"name": name}, bson.M{"$set": item})
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "updateName": name})
}

//DeleteTemplate controller
// Method: DELETE
// URL: /email/template/:name
// Header:
// Content-type
// Params:
//	name string
// Query:
// Body:
// Return:
//  success 		bool
//	removeName		string
// Errors:
//  500 "internal_server_error" "error_server"
//  404 "not_found" "template_not_found"
func DeleteTemplate(c *gin.Context) {

	collection := db.EmailTemplatesCollection()

	name := c.Param("name")

	err := collection.Remove(bson.M{"name": name})
	if err != nil {
		//Ошибка может быть как из за отсутствия в бд записи, так и в случае подключения к бд
		if err.Error() != "not found" {
			panic(err.Error())
		}

		c.JSON(errors.NotFound(errTemplateNotFound, err.Error()))
		return
	}

	c.JSON(200, gin.H{"success": true, "removeName": name})
}
