package commands

import (
	"DiscordMinecraftHelper/internal/bot"
	botrcon "DiscordMinecraftHelper/internal/server"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func PlayerListHandler(s *discordgo.Session, i *discordgo.InteractionCreate, g *botrcon.Server) {
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

func RestartServerHandler(s *discordgo.Session, i *discordgo.InteractionCreate, g *botrcon.Server) {
	bot.UpdateBotStatus(s, g)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Command disabled.",
		},
	})

	// conn, err := g.RconConnect()

	// if err != nil {
	// 	if err.Error() == "server offline" {
	// 		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 			Data: &discordgo.InteractionResponseData{
	// 				Content: "Server is offline. Contacting admins.",
	// 			},
	// 		})
	// 	} else {
	// 		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 			Data: &discordgo.InteractionResponseData{
	// 				Content: "Unable to restart server.",
	// 			},
	// 		})
	// 	}
	// 	notifyAdmin(s, i.ChannelID)
	// 	return
	// }

	// _, err = conn.Execute("/say The server will restart in 10 seconds")
	// if err != nil {
	// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: &discordgo.InteractionResponseData{
	// 			Content: "Unable to restart server.",
	// 		},
	// 	})
	// 	notifyAdmin(s, i.ChannelID)
	// } else {
	// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: &discordgo.InteractionResponseData{
	// 			Content: "Restarting server in 10 seconds. Please wait at least 5 minutes before attempting to restart the server again. If something went wrong then I'll notify the admin.",
	// 		},
	// 	})

	// 	err = g.RestartServer(conn)
	// 	if err != nil {
	// 		notifyAdmin(s, i.ChannelID)
	// 	}

	// 	conn.Close()
	// 	s.ChannelMessageSend(i.ChannelID, "Server has restarted.")
	// }
	// bot.UpdateBotStatus(s, g)
}

// TODO: Verify functionality. Dummy address supplied for now
func ServerAddressHandler(s *discordgo.Session, i *discordgo.InteractionCreate, g *botrcon.Server) {
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
		g.Logger.Warn(err.Error())
	} else {
		message += fmt.Sprintf("Server Address: %v", address)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	}
}

func notifyAdmin(s *discordgo.Session, c string) {
	admin, ok := os.LookupEnv("ADMIN")
	if !ok {
		s.ChannelMessageSend(c, "admin user not found.")
		return
	}

	s.ChannelMessageSend(c, fmt.Sprintf("<@%v> There is a problem with the server.", admin))
}
