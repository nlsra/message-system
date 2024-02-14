package config

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/nlsra/message-system/pkg/api"
	"github.com/nlsra/message-system/pkg/repository/message"
	"github.com/nlsra/message-system/pkg/service"
	"gorm.io/gorm"
)

func (a *App) InitializeRoutes() {
	a.Router = mux.NewRouter()

	a.Msg = InitAPI(a.DB, a.RedisClient)
	s := a.Router.PathPrefix("/api").Subrouter()

	s.HandleFunc("/message", a.Msg.FindAllSentMessages()).Methods("GET")
	s.HandleFunc("/message/send", a.Msg.Send()).Methods("POST")
}

func InitAPI(db *gorm.DB, client *redis.Client) *api.MessageAPI {
	msgRepository := message.NewRepository(db, client)
	msgService := service.NewMessageService(msgRepository)
	msgAPI := api.NewMessageAPI(msgService)
	msgAPI.Migrate()
	return &msgAPI
}
