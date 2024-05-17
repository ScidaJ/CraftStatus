package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/lpernett/godotenv"
)

var GuildID string
var BotToken string

var s *discordgo.Session

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

	commands := []SlashCommand{
		PlayerList,
		StartServer,
		RestartServer,
		ServerAddress,
	}

	log.Println("Adding commands...")

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

	// Make and add command list here

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	// We need to fetch the commands, since deleting requires the command ID.
	// We are doing this from the returned commands on line 375, because using
	// this will delete all the commands, which might not be desirable, so we
	// are deleting only the commands that we added.
	registeredCommands, err = s.ApplicationCommands(s.State.User.ID, GuildID)
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
