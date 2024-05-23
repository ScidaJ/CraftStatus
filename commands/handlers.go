package commands

import (
	"DiscordMinecraftHelper/bot"
	botrcon "DiscordMinecraftHelper/server"
	"fmt"
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func PlayerListHandler(s *discordgo.Session, i *discordgo.InteractionCreate, g botrcon.Server) {
	if !g.ServerRunning() {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Server not running.",
			},
		})
		return
	}

	response, err := g.ListPlayers()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
	}
}

func RestartServerHandler(s *discordgo.Session, i *discordgo.InteractionCreate, g botrcon.Server) {
	bot.UpdateBotStatus(s, g, g.Logger.With("process", "bot_status"))

	conn, err := g.RconConnect()
	if err != nil {
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
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unable to restart server.",
			},
		})
		notifyAdmin(s, i.ChannelID)
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Restarting server in 10 seconds. Please wait at least 5 minutes before attempting to restart the server again. If something went wrong then I'll notify the admin.",
			},
		})

		err = g.RestartServer(conn)
		if err != nil {
			notifyAdmin(s, i.ChannelID)
		}

		conn.Close()
		s.ChannelMessageSend(i.ChannelID, "Server has restarted.")
	}
	bot.UpdateBotStatus(s, g, g.Logger.With("process", "bot_status"))
}

func StartServerHandler(s *discordgo.Session, i *discordgo.InteractionCreate, g botrcon.Server) {
	bot.UpdateBotStatus(s, g, g.Logger.With("process", "bot_status"))
	conn, _ := g.RconConnect()
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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Starting server",
		},
	})

	err := g.StartServer()
	if err != nil {
		notifyAdmin(s, i.ChannelID)
	}
	bot.UpdateBotStatus(s, g, g.Logger.With("process", "bot_status"))
}

func ServerAddressHandler(s *discordgo.Session, i *discordgo.InteractionCreate, g botrcon.Server) {
	message := ""
	conn, err := g.RconConnect()
	if err != nil {
		message += "The server is not running. "
	} else {
		conn.Close()
	}
	address := g.GetServerAddress()
	if len(address) == 0 {
		message += "Error retrieving server address. Service may be down."
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	} else {
		message += fmt.Sprintf("Server Address: %v:25565", address)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	}
}

func notifyAdmin(s *discordgo.Session, c string) {
	err := godotenv.Load()
	if err != nil {
		slog.Error("error notifying admin", "error", err)
	}

	admin := os.Getenv("ADMIN")
	s.ChannelMessageSend(c, fmt.Sprintf("<@%v> There is a problem with the server.", admin))
}
