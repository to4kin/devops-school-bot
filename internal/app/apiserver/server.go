package apiserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

type server struct {
	router *mux.Router
	store  store.Store
	bot    *telebot.Bot
}

func newServer(store store.Store) *server {
	srv := &server{
		router: mux.NewRouter(),
		store:  store,
	}

	srv.configureRouter()

	return srv
}

func (srv *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(rw, r)
}

func (srv *server) configureRouter() {
	srv.router.HandleFunc("/", srv.botWebHookHandler()).Methods("POST")
}

func (srv *server) configureBotHandler() {
	srv.bot.Handle(telebot.OnText, srv.onTextHanlder)
}

func (srv *server) botWebHookHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var u telebot.Update
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			srv.error(rw, r, http.StatusBadRequest, err)
			return
		}

		srv.bot.ProcessUpdate(u)
	}
}

func (srv *server) error(rw http.ResponseWriter, r *http.Request, code int, err error) {
	srv.respond(rw, r, code, map[string]string{"error": err.Error()})
}

func (srv *server) respond(rw http.ResponseWriter, r *http.Request, code int, data interface{}) {
	rw.WriteHeader(code)
	if data != nil {
		json.NewEncoder(rw).Encode(data)
	}
}

func (srv *server) onTextHanlder(c telebot.Context) error {
	text := strings.ToLower(c.Message().Text)
	logrus.Debug("Text: " + text)
	for _, entity := range c.Message().Entities {
		logrus.Debug("Entity: " + entity.Type)
		logrus.Debug(text[entity.Offset : entity.Offset+entity.Length])
	}

	return nil
}
