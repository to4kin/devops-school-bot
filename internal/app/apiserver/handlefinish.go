package apiserver

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleFinish(c telebot.Context) error {
	if c.Message().Private() {
		return nil
	}

	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		srv.logger.Error(err)
		return nil
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

	if !account.Superuser {
		srv.logger.WithFields(account.LogrusFields()).Debug("account has insufficient permissions")
		return c.Reply(msgUserInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
	}

	srv.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Debug("get school by chat_id")
	school, err := srv.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		srv.logger.Error(err)
		return c.Reply(msgSchoolNotFound, &telebot.SendOptions{ParseMode: "HTML"})
	}
	srv.logger.WithFields(school.LogrusFields()).Debug("school found")

	if school.Finished {
		srv.logger.WithFields(school.LogrusFields()).Debug("school already finished")
		return c.Reply(fmt.Sprintf(msgSchoolAlreadyFinished, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
	}

	if err := srv.store.School().Finish(school); err != nil {
		srv.logger.Error(err)
	}
	srv.logger.WithFields(school.LogrusFields()).Debug("school finished")
	return c.Reply(fmt.Sprintf(msgSchoolFinished, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
}
