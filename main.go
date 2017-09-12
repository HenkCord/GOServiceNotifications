package main

import (
	"fmt"

	"os"

	"github.com/HenkCord/notifications/conf"
	"github.com/HenkCord/notifications/db"
	e "github.com/HenkCord/notifications/email"
	"github.com/HenkCord/notifications/errors"
	p "github.com/HenkCord/notifications/push"
	s "github.com/HenkCord/notifications/sms"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	conf.Init()
	StartLog()
}

func main() {

	db.InitDBConnection()
	defer db.CloseDBConnection()

	router := gin.Default()
	router.Use(Logger())

	sms := router.Group("/sms")
	{
		smsTemplate := sms.Group("/template")
		{
			smsTemplate.GET("", s.GetTemplates)
			smsTemplate.GET("/:name", s.GetTemplate)
			smsTemplate.POST("", s.CreateTemplate)
			smsTemplate.PUT("/:name", s.UpdateTemplate)
			smsTemplate.DELETE("/:name", s.DeleteTemplate)
		}
		sms.POST("/getCode", s.GetCode)
		sms.POST("/reservationFromApp", s.ReservationFromApp)
		sms.POST("/nearWorkTime", s.NearWorkTime)
		sms.POST("/reservationApprovedWith15Minutes", s.ReservationApprovedWith15Minutes)
		sms.POST("/reservationApproved", s.ReservationApproved)
		sms.POST("/noResponse", s.NoResponse)
		sms.POST("/reservationCancelled", s.ReservationCancelled)
	}

	email := router.Group("/email")
	{
		emailTemplate := email.Group("/template")
		{
			emailTemplate.GET("", e.GetTemplates)
			emailTemplate.GET("/:name", e.GetTemplate)
			emailTemplate.POST("", e.CreateTemplate)
			emailTemplate.PUT("/:name", e.UpdateTemplate)
			emailTemplate.DELETE("/:name", e.DeleteTemplate)
		}
		email.POST("/confirmEmail", e.ConfirmEmailAddress)
	}

	push := router.Group("/push")
	{
		pushTemplate := push.Group("/template")
		{
			pushTemplate.GET("", p.GetTemplates)
			pushTemplate.GET("/:name", p.GetTemplate)
			pushTemplate.POST("", p.CreateTemplate)
			pushTemplate.PUT("/:name", p.UpdateTemplate)
			pushTemplate.DELETE("/:name", p.DeleteTemplate)
		}
		push.POST("/giveReview", p.GiveReview)
	}

	host := viper.GetString("server.host")
	port := viper.GetString("server.port")

	router.Run(host + ":" + port)

}

//StartLog write info in console
func StartLog() {
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	env := os.Getenv("Env")
	fmt.Println("--------------------------------")
	fmt.Println("Server started on", host, port)
	fmt.Println("Environment: " + env)
	if env != "production" {
		fmt.Println("Activation production mode: `--environment=production` in console")
	}
	fmt.Println("--------------------------------")
}

//Logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(errors.InternalServerError("server_error", r.(string)))
				return
			}
		}()

		// before request

		c.Next()

		// after request

	}
}
