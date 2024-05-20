package main

import (
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/lpernett/godotenv"

	"DiscordMinecraftHelper/commands"
	botrcon "DiscordMinecraftHelper/server"
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

	GuildID = os.Getenv("GUILD_ID")
	BotToken = os.Getenv("BOT_TOKEN")

	// Bot init
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		Logger.Error("error creating bot", "error", err)
	}

	server := botrcon.Server{
		Logger: Logger,
	}

	err = server.StartServer()
	if err != nil {
		Logger.Warn("error starting server", "error", err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		Logger.Info("successfully logged in", "user", s.State.User.Username)
	})

	err = s.Open()
	if err != nil {
		Logger.Error("error opening Discord session", "error", err)
	}

	addCommandHandlers(s, server)
	registerCommands()

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	Logger.Info("press Ctrl+C to exit")

	c, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		c.Shutdown()
		Logger.Error("error starting cron scheduler", "error", err)
	}

	addCronJobs(c, server)

	c.Start()

	// statusTicker := time.NewTicker(10 * time.Second)
	// go func(s *discordgo.Session, guildID string) {
	// 	for {
	// 		select {
	// 		case <-statusTicker.C:
	// 			userID := "Test User ID"
	// 			precense, err := s.State.Presence(guildID, userID)

	// 			if err != nil {
	// 				log.Println("Error")
	// 				log.Println(err)
	// 			}

	// 			log.Println(precense)
	// 			log.Println("Trying to update status")
	// 			playerCount, _ := botrcon.GetPlayerCount()
	// 			if botrcon.ServerRunning() {
	// 				activity := discordgo.Activity{
	// 					Name:    "All the Mods 9",
	// 					Type:    discordgo.ActivityTypeWatching,
	// 					State:   "Online",
	// 					Details: fmt.Sprintf("%v player(s) online!", playerCount),
	// 				}
	// 				presence := discordgo.Presence{
	// 					User:   s.State.User,
	// 					Status: discordgo.StatusOnline,
	// 					Activities: []*discordgo.Activity{
	// 						&activity,
	// 					},
	// 					ClientStatus: discordgo.ClientStatus{
	// 						Desktop: discordgo.StatusOnline,
	// 						Mobile:  discordgo.StatusOnline,
	// 						Web:     discordgo.StatusOnline,
	// 					},
	// 				}

	// 				err = s.State.PresenceAdd(guildID, &presence)
	// 				if err != nil {
	// 					log.Println(err)
	// 				}
	// 			} else {
	// 				activity := discordgo.Activity{
	// 					Name:    "All the Mods 9",
	// 					Type:    discordgo.ActivityTypeWatching,
	// 					State:   "Offline",
	// 					Details: "No players online.",
	// 				}
	// 				presence := discordgo.Presence{
	// 					User:   s.State.User,
	// 					Status: discordgo.StatusOnline,
	// 					Activities: []*discordgo.Activity{
	// 						&activity,
	// 					},
	// 					ClientStatus: discordgo.ClientStatus{
	// 						Desktop: discordgo.StatusOnline,
	// 						Mobile:  discordgo.StatusOnline,
	// 						Web:     discordgo.StatusOnline,
	// 					},
	// 				}

	// 				err = s.State.PresenceAdd(guildID, &presence)
	// 				if err != nil {
	// 					log.Println(err)
	// 				}
	// 			}
	// 		case <-stop:
	// 			return
	// 		}
	// 	}
	// }(s, GuildID)

	<-stop

	c.Shutdown()

	//statusTicker.Stop()

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

func registerCommands() []*discordgo.ApplicationCommand {

	commands := commands.GetCommands()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd := &discordgo.ApplicationCommand{
			Name:        v.Name,
			Description: v.Description,
			Options:     v.Options,
		}
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, cmd)

		if err != nil {
			Logger.Error("error adding command", "command", v.Name, "error", err)
		}
		registeredCommands[i] = cmd
	}

	return registeredCommands
}

func addCommandHandlers(s *discordgo.Session, server botrcon.Server) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandHandlers := commands.GetCommandsHandlers()
		Logger.Info("command received", "command", i.ApplicationCommandData().Name, "user", i.Member.User.GlobalName)
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, server)
		}
	})
}

func addCronJobs(c gocron.Scheduler, server botrcon.Server) {
	c.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(2, 0, 0),
			),
		),
		gocron.NewTask(
			func() { server.DailyRestart() },
		),
	)
	c.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(1, 55, 0),
			),
		),
		gocron.NewTask(
			func() {
				conn, err := server.RconConnect()
				if err != nil {
					return
				} else {
					conn.Execute("/say Server will restart in 5 minutes")
				}
			},
		),
	)
	c.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(1, 30, 0),
			),
		),
		gocron.NewTask(
			func() {
				conn, err := server.RconConnect()
				if err != nil {
					return
				} else {
					conn.Execute("/say Server will restart in 30 minutes")
				}
			},
		),
	)
}