package botrcon

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorcon/rcon"
	"github.com/joho/godotenv"
)

type Server struct {
	Logger  *slog.Logger
	Players map[string]string
}

func (s Server) DailyRestart() {
	sLogger := s.Logger.With("function", "daily_restart")
	if s.ServerRunning() {
		conn, err := s.RconConnect()
		if err != nil {
			return
		}

		err = s.RestartServer(conn)
		if err != nil {
			sLogger.Warn("error restarting server", "error", err)
		}
	} else {
		err := s.StartServer()
		if err != nil {
			sLogger.Warn("error starting server", "error", err)
		}
	}
}

func (s Server) GetPlayerCount() (int, error) {
	sLogger := s.Logger.With("function", "get_player_count")
	conn, err := s.RconConnect()
	if err != nil {
		return 0, err
	}

	response, err := conn.Execute("/list")
	if err != nil {
		sLogger.Warn("error executing /list", "error", err)
		return 0, err
	}

	responses := strings.Split(response, ":")
	responseLeft := strings.Split(responses[0], " ")
	playerNumber := responseLeft[2]
	playerCount, _ := strconv.Atoi(playerNumber)

	return playerCount, nil
}

// This assumes that the bot is running on the same machine as the server. If SERVER_ADDRESS is supplied then it will return that.
func (s Server) GetServerAddress() string {
	sLogger := s.Logger.With("function", "get_server_address")

	err := godotenv.Load()
	if err != nil {
		sLogger.Error("error loading .env file", "error", err)
	}

	sAddress := os.Getenv("SERVER_ADDRESS")

	if sAddress != "" {
		return sAddress
	}

	ipService := "https://api.ipify.org"
	resp, err := http.Get(ipService)
	if err != nil {
		sLogger.Warn("error getting server address", "error", err)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		sLogger.Error("error reading response", "error", err)
	}

	return string(body)
}

func (s Server) ListPlayers() (string, error) {
	sLogger := s.Logger.With("function", "get_player_count")
	conn, err := s.RconConnect()
	if err != nil {
		return err.Error(), err
	}

	defer conn.Close()

	response, err := conn.Execute("/list")
	if err != nil {
		sLogger.Warn("error executing /list", "error", err)
		return err.Error(), err
	}

	responses := strings.Split(response, ":")
	responseLeft := strings.Split(responses[0], " ")
	responseRight := strings.ReplaceAll(responses[1], " ", "")
	usernameList := strings.Split(responseRight, ",")
	usernames, err := s.nameDecoder(usernameList)
	if err != nil {
		sLogger.Warn("error extracting player list", "error", err)
		return err.Error(), err
	}

	if responseLeft[2] == fmt.Sprint(0) {
		return "There are no players online.", nil
	} else {
		sLogger.Info("players online", "names", usernames)
		return fmt.Sprintf("There are %s player(s) online, %s", responseLeft[2], usernames), nil
	}
}

func LoadPlayerList(l *slog.Logger) (map[string]string, error) {
	l = l.With("function", "load_player_list")

	playerList := map[string]string{}

	err := godotenv.Load()
	if err != nil {
		l.Error("error loading .env file", "error", err)
		return playerList, err
	}

	playersString := os.Getenv("PLAYER_LIST")
	playersSlice := strings.Split(playersString, ",")

	for _, v := range playersSlice {
		player := strings.Split(v, ":")
		playerList[player[0]] = player[1]
	}

	return playerList, err
}

func (s Server) RconConnect() (*rcon.Conn, error) {
	sLogger := s.Logger.With("function", "rcon_connect")
	err := godotenv.Load()
	if err != nil {
		sLogger.Error("error loading .env file", "error", err)
	}

	rconAddress := os.Getenv("RCON_ADDRESS")
	rconPassword := os.Getenv("RCON_PASSWORD")

	conn, err := rcon.Dial(rconAddress, rconPassword)
	if err != nil {
		return nil, errors.New("server offline")
	}

	return conn, nil
}

func (s Server) RestartServer(conn *rcon.Conn) error {
	sLogger := s.Logger.With("function", "restart_server")
	sLogger.Info("restarting server")
	time.Sleep(10 * time.Second)

	_, err := conn.Execute("/stop")
	if err != nil {
		sLogger.Warn("error", "error", err)
		return err
	}

	err = conn.Close()
	if err != nil {
		sLogger.Warn("error", "error", err)
		return err
	}

	for i := 0; i < 10; i++ {
		time.Sleep(30 * time.Second)

		_, err := s.RconConnect()
		if err == nil {
			return nil
		}
	}

	return nil
}

func (s Server) ServerRunning() bool {
	conn, err := s.RconConnect()
	if err != nil {
		return false
	}

	defer conn.Close()

	if conn != nil {
		return true
	}

	return false
}

func (s Server) StartServer() error {
	sLogger := s.Logger.With("function", "start_server")
	if !s.ServerRunning() {
		err := godotenv.Load()
		if err != nil {
			sLogger.Error("error loading .env file", "error", err)
		}

		serverPath := os.Getenv("START_SERVER_PATH")
		c := exec.Command("cmd.exe", "/C", "Start", serverPath)
		err = c.Start()
		if err != nil {
			sLogger.Warn("unable to start server", "error", err)
			return err
		}

		time.Sleep(5 * time.Minute)

		conn, err := s.RconConnect()
		if err != nil {
			sLogger.Warn("unable to start server", "error", err)
			return err
		}
		conn.Close()

		return nil
	}
	return nil
}

func (s Server) nameDecoder(usernameList []string) (string, error) {
	var nameList strings.Builder

	for _, v := range usernameList {
		v = strings.TrimSuffix(v, "\n")

		playerName, ok := s.Players[v]
		if !ok {
			_, err := nameList.WriteString(v)
			if err != nil {
				return err.Error(), err
			}
		} else {
			_, err := nameList.WriteString(playerName)
			if err != nil {
				return err.Error(), err
			}
		}

		_, err := nameList.WriteString(", ")
		if err != nil {
			return err.Error(), err
		}
	}
	return nameList.String(), nil
}
