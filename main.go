package main

import (
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"

	"DiscordMinecraftHelper/internal/bot"
	"DiscordMinecraftHelper/internal/commands"
	botrcon "DiscordMinecraftHelper/internal/server"
)

var GuildID string
var BotToken string

var s *discordgo.Session
var Logger *slog.Logger

func main() {
	// Logger init
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(Logger)
	Logger.Info("hello, world")

	// .env init
	err := godotenv.Load()
	if err != nil {
		Logger.Error("error loading .env file", "error", err)
	}

	serverEnv := botrcon.NewServerEnv()

	GuildID = os.Getenv("GUILD_ID")
	BotToken = os.Getenv("BOT_TOKEN")

	// Bot init
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		Logger.Error("error creating bot", "error", err)
	}

	server := botrcon.Server{
		Logger: Logger,
		Env:    serverEnv,
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		Logger.Info("successfully logged in", "user", s.State.User.Username)
	})

	err = s.Open()
	if err != nil {
		Logger.Error("error opening Discord session", "error", err)
	}

	// Slash command init
	commands.AddCommandHandlers(s, &server, Logger)
	commands.RegisterCommands(s, GuildID, Logger)

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	Logger.Info("press Ctrl+C to exit")

	// Cron job init
	c, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		c.Shutdown()
		Logger.Error("error starting cron scheduler", "error", err)
	}

	c.Start()

	// Status init
	statusTicker := time.NewTicker(10 * time.Minute)
	go func(s *discordgo.Session, server *botrcon.Server) {
		for {
			select {
			case <-statusTicker.C:
				bot.UpdateBotStatus(s, server)
			case <-stop:
				return
			}
		}
	}(s, &server)

	bot.UpdateBotStatus(s, &server)

	// Shutdown
	<-stop

	c.Shutdown()

	Logger.Info("stopping statusTicker")

	statusTicker.Stop()

	Logger.Info("removing commands")

	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, GuildID)
	if err != nil {
		Logger.Error("error fetching commands", "error", err)
	}

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
		if err != nil {
			Logger.Error("error deleting command", "command", v.Name, "error", err)
		}
	}

	Logger.Info("gracefully shutting down")
}
