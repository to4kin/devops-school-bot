package apiserver

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleSchools(c telebot.Context) error {
	if !c.Message().Private() {
		return c.EditOrReply(helper.ErrWrongChatType, &telebot.SendOptions{ParseMode: "HTML"})
	}

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

	callback := &model.Callback{
		ID:          0,
		Type:        "school",
		Command:     "get",
		ListCommand: "get",
	}

	replyMessage, replyMarkup, err := helper.GetSchoolsList(srv.store, callback)
	if err != nil {
		srv.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	return c.EditOrReply(replyMessage, replyMarkup)
}
