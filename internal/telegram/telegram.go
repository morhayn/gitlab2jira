package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Tocken string
var Channel int64 = -1001514590541

type Client struct {
	bot *tgbotapi.BotAPI
}

func New() *Client {
	// fmt.Println("Tocken", Tocken)
	bot, err := tgbotapi.NewBotAPI(Tocken)
	if err != nil {
		log.Panic("!!!", err)
	}

	return &Client{
		bot: bot,
	}
}

func (c *Client) SendMessage(text string, chatId int64) error {
	msg := tgbotapi.NewMessage(Channel, text)
	msg.ParseMode = "markdown"
	_, err := c.bot.Send(msg)
	return err
}
