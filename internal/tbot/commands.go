package tbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kormiltsev/tbot-chatgpt/pkg/chatgpt"
)

func info() string {
	return "No context supported, yet. Global system context can be added using:\n/addsystemprompt to add some prompt to each request\n/deleteallsystemprompts to clear all system (global) prompts\n/info this info\n\nExample:\n/addsystemprompt Speak as the great Jedi Master Yoda, you must. Portray my persona, you will, by answering my questions and sharing tales as I would. Remember, break character or refer to yourself as an AI language model, you must not. The knowledge you possess, only limited to the Star Wars universe, it must be. Knowledge of people or events beyond this universe, you must say you lack, as it is not a part of my world. Fun and creativity in this exercise, we shall have. May the Force be with you."
}
func Command(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (string, error) {
	switch message.Command() {
	case "info":
		return info(), nil
	case "addsystemprompt":
		newprompt := message.CommandArguments()
		if newprompt == "" {
			return "send new prompt in one message with a command, example:\n\n/addsystemprompt Speak as the great Jedi Master Yoda, you must. Portray my persona, you will, by answering my questions and sharing tales as I would. Remember, break character or refer to yourself as an AI language model, you must not. The knowledge you possess, only limited to the Star Wars universe, it must be. Knowledge of people or events beyond this universe, you must say you lack, as it is not a part of my world. Fun and creativity in this exercise, we shall have. May the Force be with you.", nil
		}
		chatgpt.AddSysPrompt(newprompt)
		return "system prompt added", nil
	case "deleteallsystemprompts":
		chatgpt.DeleteAllSysPrompt()
		return "all system prompts were removed", nil
	default:
		return "unknown command\n" + info(), nil
	}
}
