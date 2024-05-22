package botrcon

import (
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
	Logger *slog.Logger
}

func (s Server) DailyRestart() {
	sLogger := s.Logger.With("process", "daily_restart")
	if s.ServerRunning() {
		conn, err := s.RconConnect()
		if err != nil {
			return
		}

		err = s.RestartServer(conn)
		if err != nil {
			sLogger.Warn("error restarting server", "process", "daily_restart", "error", err)
		}
	} else {
		err := s.StartServer()
		if err != nil {
			sLogger.Warn("error starting server", "error", err)
		}
	}
}

func (s Server) GetPlayerCount() (int, error) {
	sLogger := s.Logger.With("process", "get_player_count")
	conn, err := s.RconConnect()
	if err != nil {
		sLogger.Warn("error connecting to server", "error", err)
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

// This assumes that the bot is running on the same machine as the server. Would not be needed if hosted on dedicated server.
func (s Server) GetServerAddress() string {
	sLogger := s.Logger.With("process", "get_server_address")

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
	sLogger := s.Logger.With("process", "get_player_count")
	conn, err := s.RconConnect()
	if err != nil {
		sLogger.Warn("error connecting to server", "error", err)
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
	usernames, err := nameDecoder(usernameList)
	if err != nil {
		sLogger.Warn("error extracting player list", "error", err)
		return err.Error(), err
	}

	if responseLeft[2] == fmt.Sprint(0) {
		return "There are no players online.", nil
	} else {
		return fmt.Sprintf("There are %s player(s) online, %s", responseLeft[2], usernames), nil
	}
}

func (s Server) RconConnect() (*rcon.Conn, error) {
	sLogger := s.Logger.With("process", "rcon_connect")
	err := godotenv.Load()
	if err != nil {
		sLogger.Error("error loading .env file", "error", err)
	}

	rconAddress := os.Getenv("RCON_ADDRESS")
	rconPassword := os.Getenv("RCON_PASSWORD")

	conn, err := rcon.Dial(rconAddress, rconPassword)
	if err != nil {
		sLogger.Warn("error connecting to server", "error", err)
		return nil, err
	}

	return conn, nil
}

func (s Server) RestartServer(conn *rcon.Conn) error {
	sLogger := s.Logger.With("process", "restart_server")
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

	return true
}

func (s Server) StartServer() error {
	sLogger := s.Logger.With("process", "start_server")
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

		time.Sleep(2 * time.Minute)

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

func nameDecoder(usernames []string) (string, error) {
	var nameList strings.Builder
	names := map[string]string{
		"Beamsword":       "Kurt",
		"burgerdude9":     "Sean",
		"Rob1729":         "Rob",
		"ShermanTWilliam": "Nik",
		"ThatGuyinPJs":    "Jacob",
	}

	for _, v := range usernames {
		v = strings.TrimSuffix(v, "\n")

		_, err := nameList.WriteString(names[v])
		if err != nil {
			return err.Error(), err
		}

		_, err = nameList.WriteString(", ")
		if err != nil {
			return err.Error(), err
		}
	}
	return nameList.String(), nil
}
