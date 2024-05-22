package bot

import (
	botrcon "DiscordMinecraftHelper/server"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
)

func AddCronJobs(c gocron.Scheduler, server botrcon.Server, logger *slog.Logger) {
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

func UpdateBotStatus(s *discordgo.Session, server botrcon.Server, logger *slog.Logger) {
	logger.Info("updating status")
	playerCount, _ := server.GetPlayerCount()
	if server.ServerRunning() {
		activity := discordgo.Activity{
			Name:    fmt.Sprintf("Players: %v online", playerCount),
			Type:    discordgo.ActivityTypeWatching,
			State:   "Online",
			Details: fmt.Sprintf("%v player(s) online!", playerCount),
		}
		presence := discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				&activity,
			},
			Status: string(discordgo.StatusOnline),
			AFK:    false,
		}
		s.UpdateStatusComplex(presence)
	} else {
		activity := discordgo.Activity{
			Name:    "Server offline",
			Type:    discordgo.ActivityTypeWatching,
			State:   "Online",
			Details: "Server offline",
		}
		presence := discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				&activity,
			},
			Status: string(discordgo.StatusOnline),
			AFK:    false,
		}
		s.UpdateStatusComplex(presence)
	}
}
