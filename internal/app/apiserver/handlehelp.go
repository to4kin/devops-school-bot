package apiserver

import "gopkg.in/tucnak/telebot.v3"

func (srv *server) handleHelp(c telebot.Context) error {
	if c.Message().Private() {
		return c.Send(msgHelpPrivate, &telebot.SendOptions{ParseMode: "HTML"})
	}

	return c.Reply(msgHelpGroup, &telebot.SendOptions{ParseMode: "HTML"})
}
