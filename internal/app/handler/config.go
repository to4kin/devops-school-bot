package handler

import (
	"gopkg.in/tucnak/telebot.v3"
)

func (handler *Handler) configureBotHandlers() {
	handler.bot.Handle("/schools", handler.handleSchools)
	handler.bot.Handle("/startschool", handler.handleStartSchool)
	handler.bot.Handle("/stopschool", handler.handleStopSchool)
	handler.bot.Handle("/report", handler.handleReport)
	handler.bot.Handle("/fullreport", handler.handleFullReport)
	handler.bot.Handle("/homeworks", handler.handleHomework)

	handler.bot.Handle("/users", handler.handleUsers)
	handler.bot.Handle("/setsuperuser", handler.handleSetSuperuser)
	handler.bot.Handle("/unsetsuperuser", handler.handleUnsetSuperuser)

	handler.bot.Handle("/start", handler.handleStart)
	handler.bot.Handle("/help", handler.handleHelp)

	handler.bot.Handle("/joinstudent", handler.handleJoinStudent)
	handler.bot.Handle("/joinmodule", handler.handleJoinModule)
	handler.bot.Handle("/myreport", handler.handleMyReport)

	handler.bot.Handle(telebot.OnText, handler.handleOnText)
	handler.bot.Handle(telebot.OnPhoto, handler.handleOnText)
	handler.bot.Handle(telebot.OnDocument, handler.handleOnText)
	handler.bot.Handle(telebot.OnCallback, handler.handleCallback)

	handler.bot.Handle(telebot.OnEdited, handler.handleOnText)

	handler.logger.Info("bot handlers were successfully registered")
}
