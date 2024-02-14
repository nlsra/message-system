package config

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nlsra/message-system/pkg/api"
	"github.com/nlsra/message-system/pkg/service"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type App struct {
	Router      *mux.Router
	DB          *gorm.DB
	RedisClient *redis.Client
	Msg         *api.MessageAPI
}

var msg = service.MessageService{}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("Getting the env values")
	}

	a := App{}
	a.DBConnection()
	a.RedisConnection()
	a.InitializeRoutes()
	a.Msg.MessageService.StartSending()
	a.Start()

}

func (a *App) Start() {
	port := os.Getenv("APP_PORT")
	fmt.Println("Listening to port ", port)
	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}
