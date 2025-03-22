package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kormiltsev/tbot-chatgpt/configs"
	"github.com/kormiltsev/tbot-chatgpt/internal/chatgpt"
	"github.com/kormiltsev/tbot-chatgpt/internal/dal/bolt"
	"github.com/kormiltsev/tbot-chatgpt/internal/dal/sql"
	"github.com/kormiltsev/tbot-chatgpt/internal/tbot"
	"github.com/kormiltsev/tbot-chatgpt/internal/utils"
)

func main() {
	ctx := context.Background()
	// get env configs
	envtemp := os.Getenv("CHATGPT_API_TEMPERATURE")
	if envtemp != "" {
		envtempfloat, err := strconv.ParseFloat(envtemp, 32)
		if err != nil {
			log.Println("ENV wrong format", "CHATGPT_API_TEMPERATURE", envtemp)
		}
		configs.Temperature = float32(envtempfloat)
	}

	envmaxtok := os.Getenv("CHATGPT_API_MAX_TOKENS")
	if envmaxtok != "" {
		envmaxtokint, err := strconv.Atoi(envmaxtok)
		if err != nil {
			log.Println("ENV wrong format", "CHATGPT_API_MAX_TOKENS", envmaxtok)
		}
		configs.MaxTokens = envmaxtokint
	}

	envadminid := os.Getenv("CHATGPT_API_ADMIN_USERID")
	if envadminid != "" {
		adminID, err := strconv.ParseInt(envadminid, 10, 64)
		if err != nil {
			log.Println("ENV wrong format", "CHATGPT_API_MAX_TOKENS", envmaxtok)
		}
		configs.AdminID = adminID
	}

	// init chat gpt client
	apiKey := os.Getenv("CHATGPT_API_TOKEN")
	if apiKey == "" {
		panic("CHATGPT_API_TOKEN is required")
	}

	model := "gpt-4o-mini"
	client := chatgpt.NewClient(apiKey, model)

	// init db
	envdburi := os.Getenv("CHATGPT_API_DB_URI")
	if apiKey == "" {
		panic("CHATGPT_API_DB_URI is required")
	}
	db, err := gorm.Open(mysql.Open(envdburi), &gorm.Config{
		Logger: gormlogger.New(
			log.New(os.Stderr, "\r\n", log.LstdFlags), // io writer
			gormlogger.Config{
				SlowThreshold:             time.Second,      // Slow SQL threshold
				LogLevel:                  gormlogger.Error, // Log level // make it Silent,Error, Warn or Info
				IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,             // Don't include params in the SQL log
				Colorful:                  false,            // Disable color
			},
		),
	})
	if err != nil {
		panic(err)
	}
	// not sure it required, just in case
	defer func(db *gorm.DB) {
		sqlDB, err := db.DB()
		if err != nil {
			log.Println("Failed to get database instance", "error:", err.Error())
		}
		defer sqlDB.Close()
	}(db)
	gormdb := sql.New(db)
	if err := gormdb.Migrate(); err != nil {
		panic(err)
	}
	// init bolt for content
	boltFileAddress := os.Getenv("CHATGPT_API_BOLT_FILE")
	if boltFileAddress == "" {
		panic("CHATGPT_API_BOLT_FILE is required")
	}
	newboltdb, err := bolt.NewBoltDB(boltFileAddress)
	if err != nil {
		panic(err)
	}
	dbbolt, err := bolt.New(newboltdb, "messages")
	if err != nil {
		panic(err)
	}
	// init tgbot
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		panic("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		var replytext string
		var reportError string
		if update.Message.IsCommand() {
			replytext, err = tbot.Command(bot, update.Message)
			if err != nil {
				reportError = " ❌TbotCommandError" + err.Error()
			}
		} else {
			msg := update.Message.Text
			log.Printf("[ %s ] %s", update.Message.From.UserName, update.Message.Text)

			// save to content storage
			msgID := utils.NewUuidBytes()
			err := dbbolt.Put(msgID, "user", msg)
			if err == nil {
				err := gormdb.SaveMessage(ctx, update.Message.Chat.ID, msgID)
				if err != nil {
					reportError = reportError + " ❌GormPutError: " + err.Error()
				}
			} else {
				reportError = reportError + " ❌DbBoltError: " + err.Error()
			}

			msgesid, err := gormdb.GetMessagesByUserId(ctx, update.Message.Chat.ID)
			if err != nil {
				reportError = reportError + " ❌GormGetError: " + err.Error()
			}

			promptmessages := make([]chatgpt.Message, len(msgesid))
			for i, oldrecord := range msgesid {
				oldmsg, err := dbbolt.Get(oldrecord.MessageId[:])
				if err != nil {
					reportError = reportError + " ❌BoltGetError: " + err.Error()
					continue
				}
				promptmessages[i] = chatgpt.Message{Role: oldmsg.Role, Content: oldmsg.Message}
			}

			replytext, err = client.Chat(promptmessages, configs.Temperature, configs.MaxTokens)
			if err != nil {
				reportError = reportError + " ❌SendClientRequestError: " + err.Error()
			} else {
				botmsgID := utils.NewUuidBytes()
				if err := dbbolt.Put(botmsgID, "assistant", replytext); err != nil {
					reportError = reportError + " ❌BoltPutResponseError: " + err.Error()
				} else {
					if err := gormdb.SaveMessage(ctx, update.Message.Chat.ID, botmsgID); err != nil {
						reportError = reportError + " ❌GormSaveResponseError: " + err.Error()
					}
				}
			}
		}
		// send errors to admin
		if reportError != "" && configs.AdminID > 0 {
			reportError = "User: " + strconv.FormatInt(update.Message.Chat.ID, 10) + "\nUsername: " + update.Message.Chat.UserName + reportError
			errorMessageToAdmin := tgbotapi.NewMessage(configs.AdminID, reportError)
			if _, err := bot.Send(errorMessageToAdmin); err != nil {
				log.Println("Failed to send message:", err)
			}
		}

		if replytext == "" {
			replytext = "❌Internal error❌\nPlease try again later."
		}
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, replytext)
		if _, err := bot.Send(reply); err != nil {
			log.Println("Failed to send message:", err)
		}

	}
}
