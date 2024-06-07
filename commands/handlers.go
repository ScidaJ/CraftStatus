package commands

import (
	"DiscordMinecraftHelper/bot"
	botrcon "DiscordMinecraftHelper/server"
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

	conn, err := g.RconConnect()

	if err != nil {
		if err.Error() == "server offline" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Server is offline. Attempting to start server.",
				},
			})
			err = g.StartServer()
			if err != nil {
				notifyAdmin(s, i.ChannelID)
				return
			}

			s.ChannelMessageSend(i.ChannelID, "Server has started.")
			return
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Unable to restart server.",
				},
			})
		}
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
	bot.UpdateBotStatus(s, g)
}

func StartServerHandler(s *discordgo.Session, i *discordgo.InteractionCreate, g *botrcon.Server) {
	bot.UpdateBotStatus(s, g)

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
		return
	}
	s.ChannelMessageSend(i.ChannelID, "Server has started.")
	bot.UpdateBotStatus(s, g)
}

func StopServerHandler(s *discordgo.Session, i *discordgo.InteractionCreate, g *botrcon.Server) {
	admin, ok := os.LookupEnv("ADMIN")
	if !ok {
		s.ChannelMessageSend(i.ChannelID, "admin user not found.")
		return
	}

	if i.Member.User.ID != admin {
		g.Logger.Warn(fmt.Sprintf("non-admin user using /stop: %v", i.Member.User.Username))
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are not authorized to use this command.",
			},
		})
		return
	}

	conn, _ := g.RconConnect()
	if conn == nil && g.Process == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Server is not running.",
			},
		})
		return
	} else if conn == nil && g.Process != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unable to connect to server.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Stopping server",
		},
	})

	err := g.StopServer()
	if err != nil {
		notifyAdmin(s, i.ChannelID)
		return
	}
	s.ChannelMessageSend(i.ChannelID, "Server has stopped.")
	bot.UpdateBotStatus(s, g)
}

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
		g.Logger.Warn("error", err)
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
	admin, ok := os.LookupEnv("ADMIN")
	if !ok {
		s.ChannelMessageSend(c, "admin user not found.")
		return
	}

	s.ChannelMessageSend(c, fmt.Sprintf("<@%v> There is a problem with the server.", admin))
}
