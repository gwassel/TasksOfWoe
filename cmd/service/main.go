package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/bot"
	add_handler "github.com/gwassel/TasksOfWoe/internal/handler/add"
	complete_handler "github.com/gwassel/TasksOfWoe/internal/handler/complete"
	list_handler "github.com/gwassel/TasksOfWoe/internal/handler/list"
	listall_handler "github.com/gwassel/TasksOfWoe/internal/handler/listall"
	add_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/add"
	complete_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/complete"
	list_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/list"
	listall_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/listall"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

func main() {
	// Read environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	// Connect to PostgreSQL
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		dbHost, dbUser, dbName, dbPassword, dbPort)

	// Initialize Telegram Bot
	botApi, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	botApi.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botApi.GetUpdatesChan(u)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"/var/log/task-tracker/app.log"}
	config.ErrorOutputPaths = []string{"/var/log/task-tracker/error.log"}
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Info("started")

	// usecase
	completeUsecase := complete_usecase.NewUsecase(db)
	addUsecase := add_usecase.NewUsecase(db)
	listUsecase := list_usecase.NewUsecase(db)
	listallUsecase := listall_usecase.NewUsecase(db)

	// handler
	completeHandler := complete_handler.New(sugar, botApi, completeUsecase)
	addHandler := add_handler.New(sugar, botApi, addUsecase)
	listHandler := list_handler.New(sugar, botApi, listUsecase)
	listallHandler := listall_handler.New(sugar, botApi, listallUsecase)

	handlersMap := map[string]interface {
		Handle(message *tgbotapi.Message)
	}{
		"com":     completeHandler,
		"add":     addHandler,
		"list":    listHandler,
		"listall": listallHandler,
	}

	handler := bot.NewBot(botApi, handlersMap)

	// Handle incoming updates
	for update := range updates {
		if update.Message != nil { // Ignore non-Message updates
			handler.HandleMessage(update.Message)
		}
	}
}
