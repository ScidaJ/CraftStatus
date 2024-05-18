package commands

import (
	botrcon "DiscordMinecraftHelper/server"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func PlayerListHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if !botrcon.ServerRunning() {
		log.Println("Server not running.")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Server not running.",
			},
		})
		notifyAdmin(s, i.ChannelID)
		return
	}

	response, err := botrcon.ListPlayers()
	if err != nil {
		log.Println(discordgo.ErrCodeActionRequiredVerifiedAccount)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func RestartServerHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	conn, err := botrcon.RconConnect()
	if err != nil {
		log.Printf("Error restarting server: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unable to restart server.",
			},
		})
		notifyAdmin(s, i.ChannelID)
		return
	}

	_, err = conn.Execute("/say The server will restart in 10 seconds")
	if err != nil {
		log.Printf("Error restarting server: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unable to restart server.",
			},
		})
		notifyAdmin(s, i.ChannelID)
		return
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Restarting server in 10 seconds. Please wait at least 5 minutes before attempting to restart the server again. If something went wrong then I'll notify the admin.",
			},
		})
		err = botrcon.RestartServer(conn)

		if err != nil {
			notifyAdmin(s, i.ChannelID)
		}

		conn.Close()

		s.ChannelMessageSend(i.ChannelID, "Server has restarted.")
	}
}

func StartServerHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	conn, err := botrcon.RconConnect()
	if conn != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Server is already running.",
			},
		})
		conn.Close()
		return
	}

	if err != nil {
		log.Printf("Error starting server: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unable to start server.",
			},
		})
		notifyAdmin(s, i.ChannelID)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Starting server",
		},
	})

	err = botrcon.StartServer()

	if err != nil {
		log.Printf("Error starting server: %v", err)
		notifyAdmin(s, i.ChannelID)
		return
	}
}

func ServerAddressHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	conn, err := botrcon.RconConnect()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Server not running.",
			},
		})
	} else {
		conn.Close()

		address := botrcon.GetServerAddress()

		if len(address) == 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error retrieving server address. Service may be down.",
				},
			})
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Server Address: %v:25565", address),
				},
			})
		}
	}
}

func notifyAdmin(s *discordgo.Session, c string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	admin := os.Getenv("ADMIN")

	s.ChannelMessageSend(c, fmt.Sprintf("@%v There is a problem with the server.", admin))
}
