package main

import (
	"log"
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

func init() {
	err := botrcon.StartServer()

	if err != nil {
		log.Printf("Unable to start server %v", err)
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GuildID = os.Getenv("GUILD_ID")
	BotToken = os.Getenv("BOT_TOKEN")
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	addCommandHandlers(s)
	registerCommands()

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")

	c, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		c.Shutdown()
		log.Fatalf("Cron scheduler failed to start %v", err)
	}

	c.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(2, 0, 0),
			),
		),
		gocron.NewTask(
			func() { botrcon.DailyRestart() },
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
				conn, err := botrcon.RconConnect()
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
				conn, err := botrcon.RconConnect()
				if err != nil {
					return
				} else {
					conn.Execute("/say Server will restart in 30 minutes")
				}
			},
		),
	)

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

	log.Println("Removing commands...")
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, GuildID)
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
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
			log.Fatalf("error adding command %v\n%v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	return registeredCommands
}

func addCommandHandlers(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandHandlers := commands.GetCommandsHandlers()
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
