package websocket

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/maronfranc/poc-golang-ddd/domain/dto"
)

var upgrader = websocket.Upgrader{
	// ReadBufferSize:  256,
	// WriteBufferSize: 256,
	// WriteBufferPool: &sync.Pool{},
}

func LoadHelloWebsocket() chi.Router {
	r := chi.NewRouter()

	r.Get("/hello", onWebsocket)

	return r
}

func process(ws *websocket.Conn) {
	defer ws.Close()
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("[Read:error]: %v", err)
			break
		}
		log.Printf("[Read:message]: %s", msg)

		for i := 1; i <= 5; i++ {
			time.Sleep(1 * time.Second)
			msg := fmt.Sprintf("Response from server %d", i)

			if err := ws.WriteJSON(dto.ResponseMessage{Message: msg}); err != nil {
				log.Printf("[Response:error]: %v", err)
			}
		}
	}
}

func onWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[Error:upgrade]: %v", err)
		return
	}
	go process(c)
}
