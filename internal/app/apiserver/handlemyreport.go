package apiserver

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleMyReport(c telebot.Context) error {
	if c.Message().Private() {
		return c.EditOrReply(helper.ErrWrongChatType, &telebot.SendOptions{ParseMode: "HTML"})
	}

	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		srv.logger.Error(err)

		if err == store.ErrRecordNotFound {
			return c.EditOrReply(helper.ErrUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

	srv.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Debug("get school by chat_id")
	school, err := srv.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		srv.logger.Error(err)

		if err == store.ErrRecordNotFound {
			return c.EditOrReply(helper.ErrSchoolNotStarted, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	srv.logger.WithFields(school.LogrusFields()).Debug("school found")

	hlpr := helper.NewHelper(srv.store, srv.logger)
	reportMessage, err := hlpr.GetUserReport(account, school)
	if err != nil {
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	srv.logger.Debug("report sent")
	return c.EditOrReply(reportMessage, &telebot.SendOptions{ParseMode: "HTML"})
}
