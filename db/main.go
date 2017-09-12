package db

import (
	"log"
	"time"

	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
)

// Service declare database abilities
type service struct {
	databaseName             string
	smsCollection            string
	smsTemplatesCollection   string
	emailCollection          string
	emailTemplatesCollection string
	pushCollection           string
	pushTemplatesCollection  string
	session                  *mgo.Session
}

var notificationsService = &service{
	databaseName:             "notifications",
	smsCollection:            "sms",
	smsTemplatesCollection:   "smsTemplates",
	emailCollection:          "email",
	emailTemplatesCollection: "emailTemplates",
	pushCollection:           "push",
	pushTemplatesCollection:  "pushTemplates",
}

//InitDBConnection init connection to DB and set Session variable
func InitDBConnection() {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{viper.GetString("db.mongo.host") + ":" + viper.GetString("db.mongo.port")},
		Timeout:  30 * time.Second,
		Database: notificationsService.databaseName,
	}
	var err error
	notificationsService.session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Println("MongoDB: Error session created")
		return
	}
	// Optional. Switch the session to a monotonic behavior.
	notificationsService.session.SetMode(mgo.Monotonic, true)
	log.Println("MongoDB: Session created")
}

//CloseDBConnection close connection to DB
func CloseDBConnection() {
	notificationsService.session.Close()
	log.Println("MongoDB: Session closed")
}

// GetDb return connection to database
func GetDb() *mgo.Database {
	return notificationsService.session.DB(notificationsService.databaseName)
}

// SMSCollection return connection to sms collection
func SMSCollection() *mgo.Collection {
	return notificationsService.session.DB(notificationsService.databaseName).C(notificationsService.smsCollection)
}

// SMSTemplatesCollection return connection to smsTemplates collection
func SMSTemplatesCollection() *mgo.Collection {
	return notificationsService.session.DB(notificationsService.databaseName).C(notificationsService.smsTemplatesCollection)
}

// EmailCollection return connection to email collection
func EmailCollection() *mgo.Collection {
	return notificationsService.session.DB(notificationsService.databaseName).C(notificationsService.emailCollection)
}

// EmailTemplatesCollection return connection to emailTemplates collection
func EmailTemplatesCollection() *mgo.Collection {
	return notificationsService.session.DB(notificationsService.databaseName).C(notificationsService.emailTemplatesCollection)
}

// PushCollection return connection to push collection
func PushCollection() *mgo.Collection {
	return notificationsService.session.DB(notificationsService.databaseName).C(notificationsService.pushCollection)
}

// PushTemplatesCollection return connection to pushTemplates collection
func PushTemplatesCollection() *mgo.Collection {
	return notificationsService.session.DB(notificationsService.databaseName).C(notificationsService.pushTemplatesCollection)
}
