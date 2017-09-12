package conf

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Init() {
	InitConfig()
	InitFlags()
	ChangeConfig()
}

//InitFlags read console flags (--environment=production)
//and set Go environment
func InitFlags() {
	environment := flag.String("environment", "default", "a string")
	flag.Parse()
	env := Setenv(*environment)
	os.Setenv("Env", env) // "development" "production" "local"
	if env == "production" {
		gin.SetMode(gin.ReleaseMode) //default gin.DebugMode
	}
}

func Setenv(environment string) string {
	if environment == "default" {
		envDefault := viper.GetString("environment.default")
		if envDefault != "" {
			return envDefault
		}
	} else if environment != "" {
		return environment
	}
	return "development"
}

func InitConfig() {
	viper.SetConfigFile("./conf/config.yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func ChangeConfig() {
	env := os.Getenv("Env")
	server := viper.Get("environment." + env + ".server")
	viper.Set("server", server)
	db := viper.Get("environment." + env + ".db")
	viper.Set("db", db)
}
