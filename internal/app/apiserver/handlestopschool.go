package apiserver

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	msgSchoolFinished string = `school <b>%v</b> finished`
)

func (srv *server) handleStopSchool(c telebot.Context) error {
	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		srv.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

	if !account.Superuser {
		srv.logger.WithFields(account.LogrusFields()).Debug("account has insufficient permissions")
		return c.EditOrReply(helper.ErrInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
	}

	if c.Message().Private() {
		callback := &model.Callback{
			ID:          0,
			Type:        "school",
			Command:     "stop",
			ListCommand: "stop",
		}

		hlpr := helper.NewHelper(srv.store, srv.logger)
		replyMessage, replyMarkup, err := hlpr.GetSchoolsList(callback)
		if err != nil {
			srv.logger.Error(err)
			return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return c.EditOrReply(replyMessage, replyMarkup)
	}

	srv.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Debug("get school by chat_id")
	school, err := srv.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		srv.logger.Error(err)
		return c.EditOrReply(helper.ErrSchoolNotStarted, &telebot.SendOptions{ParseMode: "HTML"})
	}
	srv.logger.WithFields(school.LogrusFields()).Debug("school found")

	if !school.Active {
		srv.logger.WithFields(school.LogrusFields()).Debug("school already finished")
		return c.EditOrReply(fmt.Sprintf(msgSchoolFinished, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
	}

	school.Active = false
	if err := srv.store.School().Update(school); err != nil {
		srv.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	srv.logger.WithFields(school.LogrusFields()).Debug("school finished")
	return c.EditOrReply(fmt.Sprintf(msgSchoolFinished, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
}
