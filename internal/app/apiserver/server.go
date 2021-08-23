package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
	bot    *telebot.Bot
}

func newServer(store store.Store) *server {
	srv := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
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

func (srv *server) configureLogger(logLevel string) {
	if level, err := logrus.ParseLevel(logLevel); err != nil {
		srv.logger.Error(err)
		srv.logger.SetLevel(logrus.InfoLevel)
	} else {
		srv.logger.SetLevel(level)
	}
}

func (srv *server) configureBotHandler() {
	srv.bot.Handle("/schools", srv.handleSchools)
	srv.bot.Handle("/startschool", srv.handleStartSchool)
	srv.bot.Handle("/stopschool", srv.handleStopSchool)
	srv.bot.Handle("/report", srv.handleReport)
	srv.bot.Handle("/fullreport", srv.handleFullReport)
	srv.bot.Handle("/homeworks", srv.handleHomework)

	srv.bot.Handle("/users", srv.handleUsers)
	srv.bot.Handle("/setsuperuser", srv.handleSetSuperuser)
	srv.bot.Handle("/unsetsuperuser", srv.handleUnsetSuperuser)

	srv.bot.Handle("/start", srv.handleStart)
	srv.bot.Handle("/help", srv.handleHelp)

	srv.bot.Handle("/joinstudent", srv.handleJoin)
	srv.bot.Handle("/joinmodule", srv.handleJoin)
	srv.bot.Handle("/myreport", srv.handleMyReport)

	srv.bot.Handle(telebot.OnText, srv.handleOnText)
	srv.bot.Handle(telebot.OnCallback, srv.handleCallback)
}

func (srv *server) botWebHookHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var u telebot.Update
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			srv.error(rw, r, http.StatusBadRequest, err)
			return
		}

		srv.bot.ProcessUpdate(u)
		srv.respond(rw, r, http.StatusOK, nil)
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
