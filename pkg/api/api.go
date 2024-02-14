package api

import (
	"github.com/nlsra/message-system/pkg/service"
	"log"
	"net/http"
)

type MessageAPI struct {
	MessageService service.MessageService
}

func NewMessageAPI(msg service.MessageService) MessageAPI {
	return MessageAPI{MessageService: msg}
}

func (msg MessageAPI) FindAllSentMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		msgs, err := msg.MessageService.FindAllSentMessages()
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, msgs)
	}
}

func (msg MessageAPI) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res := "message sending started successfully!."

		rs := msg.MessageService.Start()
		if rs != true {
			msg.MessageService.Stop()
			res = "message sending stopped successfully!."
		}

		RespondWithJSON(w, http.StatusOK, res)
	}
}

func (msg MessageAPI) Migrate() {
	err := msg.MessageService.Migrate()
	if err != nil {
		log.Println(err)
	}
}
