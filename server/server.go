package botrcon

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorcon/rcon"
	"github.com/joho/godotenv"
)

func DailyRestart() {
	if ServerRunning() {
		conn, err := RconConnect()
		if err != nil {
			return
		}

		RestartServer(conn)
	} else {
		StartServer()
	}
}

func GetPlayerCount() (int, error) {
	conn, err := RconConnect()
	if err != nil {
		log.Print(err)
		return 0, err
	}

	response, err := conn.Execute("/list")
	if err != nil {
		log.Print(err)
		return 0, err
	}

	responses := strings.Split(response, ":")
	responseLeft := strings.Split(responses[0], " ")
	playerNumber := responseLeft[2]
	playerCount, _ := strconv.Atoi(playerNumber)

	return playerCount, nil
}

// This assumes that the bot is running on the same machine as the server. Would not be needed if hosted on dedicated server.
func GetServerAddress() string {
	ipService := "https://api.ipify.org"
	resp, err := http.Get(ipService)
	if err != nil {
		log.Printf("Unable to get server address: %v", err)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

func ListPlayers() (string, error) {
	conn, err := RconConnect()
	if err != nil {
		log.Print(err)
		return err.Error(), err
	}

	defer conn.Close()

	response, err := conn.Execute("/list")
	if err != nil {
		log.Print(err)
		return err.Error(), err
	}

	responses := strings.Split(response, ":")
	responseLeft := strings.Split(responses[0], " ")
	responseRight := strings.ReplaceAll(responses[1], " ", "")
	usernameList := strings.Split(responseRight, ",")
	usernames, err := nameDecoder(usernameList)
	if err != nil {
		log.Println("Error listing players")
		return err.Error(), err
	}

	if responseLeft[2] == fmt.Sprint(0) {
		return "There are no players online.", nil
	} else {
		return fmt.Sprintf("There are %s player(s) online, %s", responseLeft[2], usernames), nil
	}
}

func RconConnect() (*rcon.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rconAddress := os.Getenv("RCON_ADDRESS")
	rconPassword := os.Getenv("RCON_PASSWORD")

	conn, err := rcon.Dial(rconAddress, rconPassword)
	if err != nil {
		log.Println("Error connecting to server")
		return nil, err
	}

	return conn, nil
}

func RestartServer(conn *rcon.Conn) error {
	time.Sleep(10 * time.Second)

	conn.Execute("/stop")
	conn.Close()

	for i := 0; i < 10; i++ {
		time.Sleep(30 * time.Second)

		_, err := RconConnect()

		if err == nil {
			return nil
		}
	}

	return nil
}

func ServerRunning() bool {
	conn, err := RconConnect()

	if conn != nil {
		return true
	}

	if err != nil {
		log.Println("Server Not Running")
		return false
	}

	return true
}

func StartServer() error {
	if !ServerRunning() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		serverPath := os.Getenv("START_SERVER_PATH")
		log.Println(serverPath)

		log.Println("Starting server")

		c := exec.Command("cmd.exe", "/C", "Start", serverPath)
		err = c.Start()

		if err != nil {
			log.Printf("Unable to start server: %v", err)
			return err
		}

		time.Sleep(2 * time.Minute)

		conn, err := RconConnect()

		if err != nil {
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
