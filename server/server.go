// @TODO: Capture CMD process on start. Adjust restart function. Add /stop command with admin protection
// @TODO: Add StopServer function which sends SIGINT to Cmd process. V2 of bot
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
	"syscall"
	"time"

	"github.com/gorcon/rcon"
)

type Server struct {
	Logger  *slog.Logger
	Env     ServerEnv
	Process *os.Process
}

func (s *Server) DailyRestart() {
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
	}
}

func (s *Server) GetPlayerCount() (int, error) {
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
func (s *Server) GetServerAddress() string {
	sLogger := s.Logger.With("function", "get_server_address")

	if s.Env.SERVER_ADDRESS != "" {
		return s.Env.SERVER_ADDRESS
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

func (s *Server) ListPlayers() (string, error) {
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

func (s *Server) RconConnect() (*rcon.Conn, error) {
	conn, err := rcon.Dial(s.Env.RCON_ADDRESS, s.Env.RCON_PASSWORD)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("raw error: %v", err))
		return nil, errors.New("server offline")
	}

	return conn, nil
}

func (s *Server) RestartServer(conn *rcon.Conn) error {
	sLogger := s.Logger.With("function", "restart_server")
	sLogger.Info("restarting server")

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

func (s *Server) ServerRunning() bool {
	conn, err := s.RconConnect()
	if err != nil {
		if s.Process != nil {
			s.Logger.Warn(fmt.Sprintf("server process running, unable to connect to server through rcon: %v", err))
		}
		return false
	}

	defer conn.Close()

	return conn != nil
}

func (s *Server) StartServer() error {
	sLogger := s.Logger.With("function", "start_server")
	if !s.ServerRunning() {
		cmd := exec.Command("cmd.exe", "/C", "Start", s.Env.START_SERVER_PATH)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}
		err := cmd.Start()
		if err != nil {
			sLogger.Warn("unable to start server", "error", err)
			return err
		}

		if cmd.Process != nil {
			sLogger.Info("attached server process")
			s.Process = cmd.Process
		}

		for i := 0; i < 10; i++ {
			time.Sleep(30 * time.Second)

			_, err := s.RconConnect()
			if err == nil {
				return nil
			}
		}
		return err
	}
	return nil
}

func (s *Server) StopServer() error {
	sLogger := s.Logger.With("function", "stop_server")
	if s.Process == nil {
		sLogger.Warn("unable to stop server, server process not running")
		return errors.New("server process not running")
	}

	_, err := s.RconConnect()
	if err != nil {
		return err
	}

	err = s.Process.Signal(syscall.SIGKILL)
	if err != nil {
		sLogger.Warn(fmt.Sprintf("unable to kill process %v", err))
		return err
	}

	time.Sleep(30 * time.Second)

	if s.Process != nil {
		sLogger.Warn("process still running")
		return errors.New("server process not stopped")
	}

	return nil
}

func (s *Server) nameDecoder(usernameList []string) (string, error) {
	var nameList strings.Builder

	for _, v := range usernameList {
		v = strings.TrimSuffix(v, "\n")

		playerName, ok := s.Env.PLAYER_LIST[v]
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
