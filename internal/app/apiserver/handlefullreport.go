package apiserver

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleFullReport(c telebot.Context) error {
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

	hlpr := helper.NewHelper(srv.store, srv.logger)

	if c.Message().Private() {
		callback := &model.Callback{
			ID:          0,
			Type:        "school",
			Command:     "full_report",
			ListCommand: "full_report",
		}

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

		if err == store.ErrRecordNotFound {
			return c.EditOrReply(helper.ErrSchoolNotStarted, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	srv.logger.WithFields(school.LogrusFields()).Debug("school found")

	reportMessage, err := hlpr.GetFullReport(school)
	if err != nil && err != store.ErrRecordNotFound {
		srv.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	if err == store.ErrRecordNotFound {
		reportMessage = helper.ErrReportNotFound
	}

	srv.logger.Debug("report sent")
	return c.EditOrReply(fmt.Sprintf("School <b>%v</b>\n\n%v", school.Title, reportMessage), &telebot.SendOptions{ParseMode: "HTML"})
}
