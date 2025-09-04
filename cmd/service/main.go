package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gwassel/TasksOfWoe/internal/bot"
	"github.com/gwassel/TasksOfWoe/internal/domain/encoder"
	add_handler "github.com/gwassel/TasksOfWoe/internal/handler/add"
	complete_handler "github.com/gwassel/TasksOfWoe/internal/handler/complete"
	description_handler "github.com/gwassel/TasksOfWoe/internal/handler/description"
	help_handler "github.com/gwassel/TasksOfWoe/internal/handler/help"
	list_handler "github.com/gwassel/TasksOfWoe/internal/handler/list"
	listall_handler "github.com/gwassel/TasksOfWoe/internal/handler/listall"
	take_handler "github.com/gwassel/TasksOfWoe/internal/handler/take"
	untake_handler "github.com/gwassel/TasksOfWoe/internal/handler/untake"
	"github.com/gwassel/TasksOfWoe/internal/persistence"
	add_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/add"
	complete_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/complete"
	description_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/description"
	help_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/help"
	list_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/list"
	listall_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/listall"
	take_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/take"
	untake_usecase "github.com/gwassel/TasksOfWoe/internal/usecase/untake"
	_ "github.com/lib/pq"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Read environment variables
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	encKey := os.Getenv("ENCRYPTION_KEY")

	// Connect to PostgreSQL

	db, err := persistence.NewDB()
	if err != nil {
		panic(err.Error())
	}
	defer func() { _ = db.Close() }()

	// Initialize Telegram Bot
	botApi, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	botApi.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botApi.GetUpdatesChan(u)

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
	defer lo.Must0(logger.Sync())
	sugar := logger.Sugar()
	sugar.Info("started")

	encoder, err := encoder.New(encKey)
	if err != nil {
		logger.Fatal("failed to create encoder")
	}

	// usecase
	completeUsecase := complete_usecase.NewUsecase(sugar, db)
	addUsecase := add_usecase.NewUsecase(db, encoder)
	listUsecase := list_usecase.NewUsecase(sugar, db, encoder)
	listallUsecase := listall_usecase.NewUsecase(db, encoder)
	takeUsecase := take_usecase.NewUsecase(sugar, db)
	untakeUsecase := untake_usecase.NewUsecase(sugar, db)
	descriptionUsecase := description_usecase.NewUsecase(sugar, db, encoder)
	helpUsecase := help_usecase.NewUsecase()

	// command descriptions
	descs := map[string]help_handler.HelpEntry{
		"com":         {Is_alias: true, Desc: completeUsecase.Desc},
		"complete":    {Is_alias: false, Desc: completeUsecase.Desc},
		"add":         {Is_alias: false, Desc: addUsecase.Desc},
		"list":        {Is_alias: false, Desc: listUsecase.Desc},
		"ls":          {Is_alias: true, Desc: listUsecase.Desc},
		"listall":     {Is_alias: false, Desc: listallUsecase.Desc},
		"la":          {Is_alias: true, Desc: listallUsecase.Desc},
		"take":        {Is_alias: false, Desc: takeUsecase.Desc},
		"untake":      {Is_alias: false, Desc: untakeUsecase.Desc},
		"description": {Is_alias: false, Desc: descriptionUsecase.Desc},
		"desc":        {Is_alias: true, Desc: descriptionUsecase.Desc},
		"help":        {Is_alias: false, Desc: helpUsecase.Desc},
	}

	// handler
	completeHandler := complete_handler.New(sugar, botApi, completeUsecase)
	addHandler := add_handler.New(sugar, botApi, addUsecase)
	listHandler := list_handler.New(sugar, botApi, listUsecase)
	listallHandler := listall_handler.New(sugar, botApi, listallUsecase)
	takeHandler := take_handler.New(sugar, botApi, takeUsecase)
	untakeHandler := untake_handler.New(sugar, botApi, untakeUsecase)
	descriptionHandler := description_handler.New(sugar, botApi, descriptionUsecase)
	helpHandler := help_handler.New(sugar, botApi, descs)

	handlersMap := map[string]interface {
		Handle(message *tgbotapi.Message)
	}{
		"add":     addHandler,
		"list":    listHandler,
		"listall": listallHandler,
		"desc":    descriptionHandler,
		"take":    takeHandler,
		"untake":  untakeHandler,
		"com":     completeHandler,
		"help":    helpHandler,
	}

	handler := bot.NewBot(botApi, sugar, handlersMap)

	// Handle incoming updates
	for update := range updates {
		if update.Message != nil { // Ignore non-Message updates
			handler.HandleMessage(update.Message)
		}
	}
}
