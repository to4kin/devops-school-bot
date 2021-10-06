package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
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
		logger: configureLogger(),
		store:  store,
	}

	srv.configureRouter()

	return srv
}

func configureLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetReportCaller(true)
	logger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	logger.SetLevel(logrus.InfoLevel)
	return logger
}

func (srv *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(rw, r)
}

func (srv *server) configureRouter() {
	srv.router.HandleFunc("/", srv.botWebHookHandler()).Methods("POST")
	srv.logger.Info("router handlers were successfully registered")
}

func (srv *server) configureCron(schedule string, fullreport bool) {
	cron := gocron.NewScheduler(time.UTC)
	if _, err := cron.Cron(schedule).Do(func() {
		schools, err := srv.store.School().FindByActive(true)
		if err != nil {
			srv.logger.Error(err)
		}

		hlpr := helper.NewHelper(srv.store, srv.logger)
		for _, school := range schools {
			schoolChat, err := srv.bot.ChatByID(school.ChatID)
			if err != nil {
				srv.logger.Error(err)
				continue
			}

			var reportMessage string
			if fullreport {
				reportMessage, err = hlpr.GetFullReport(school)
			} else {
				reportMessage, err = hlpr.GetReport(school)
			}

			if err != nil {
				srv.logger.Error(err)
				continue
			}

			srv.logger.WithFields(logrus.Fields{
				"school_title": school.Title,
				"fullreport":   fullreport,
			}).Debug("sent report by cron")
			srv.bot.Send(schoolChat, fmt.Sprintf("School <b>%v</b>\n\n%v", school.Title, reportMessage), &telebot.SendOptions{ParseMode: "HTML"})
		}

	}); err != nil {
		srv.logger.Error(err)
	} else {
		srv.logger.WithFields(logrus.Fields{
			"schedule":   schedule,
			"fullreport": fullreport,
		}).Info("cron was successfully registered")
		cron.StartAsync()
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

	srv.bot.Handle("/joinstudent", srv.handleJoinStudent)
	srv.bot.Handle("/joinmodule", srv.handleJoinModule)
	srv.bot.Handle("/myreport", srv.handleMyReport)

	srv.bot.Handle(telebot.OnText, srv.handleOnText)
	srv.bot.Handle(telebot.OnPhoto, srv.handleOnText)
	srv.bot.Handle(telebot.OnDocument, srv.handleOnText)
	srv.bot.Handle(telebot.OnCallback, srv.handleCallback)

	srv.bot.Handle(telebot.OnEdited, srv.handleOnText)

	srv.logger.Info("bot handlers were successfully registered")
}

func (srv *server) botWebHookHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var u telebot.Update
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			srv.error(rw, r, http.StatusBadRequest, err)
			return
		}

		srv.logger.WithFields(logrus.Fields{
			"update_id": u.ID,
		}).Info("new message received")
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
