package data_analysis

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kaeba0616/es_backend/internal/erapi"
)

type Config struct {
	EsapiClient erapi.Client
	CurrentUser *erapi.User
	Users       []erapi.User
	Rankers     []erapi.User
	Nextgame    *int
}

func StartRepl(cfg *Config, args ...string) {
	time, _ := cfg.EsapiClient.TimeList()
	fmt.Println()
	fmt.Println("Time List")
	for i, item := range time {
		fmt.Printf("Row %d: Code=%d, Name=%s, Seconds=%d, Total=%d\n", i+1, item.Code, item.Name, item.Seconds, item.Total)
	}

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("ER > ")
		reader.Scan()
		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}
		commandName := words[0]
		commandName = strings.ToLower(commandName)

		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		command, exists := getCommands()[commandName]
		if !exists {
			fmt.Printf("Command %s not found\n", commandName)
			continue
		}
		err := command.callback(cfg, args...)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}

	}
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	return words
}

type command struct {
	name        string
	descrpition string
	callback    func(*Config, ...string) error
}

func getCommands() map[string]command {
	return map[string]command{
		"help": {
			name:        "help",
			descrpition: "Displays this help message",
			callback:    commandHelp,
		},
		"config": {
			name:        "config",
			descrpition: "Displays config information",
			callback:    commandConfig,
		},
		"metatype": {
			name:        "metatype",
			descrpition: "Displays metatype information",
			callback:    commandMetatype,
		},
		"characterinfo": {
			name:        "characterinfo",
			descrpition: "Displays character information",
			callback:    commandCharacterInfo,
		},
		"weaponinfo": {
			name:        "weaponinfo",
			descrpition: "Displays weapon information",
			callback:    commandWeaponInfo,
		},
		"currentuser": {
			name:        "currentuser",
			descrpition: "Displays current user information",
			callback:    commandCurrentUser,
		},
		"user": {
			name:        "user",
			descrpition: "search usernum by nickname",
			callback:    commandUser,
		},
		"userlist": {
			name:        "userlist",
			descrpition: "Displays users information before I found out",
			callback:    commandUserList,
		},
		"usergame": {
			name:        "usergame",
			descrpition: "search user game list",
			callback:    commandUserGame,
		},
		"gamedetail": {
			name:        "gamedetail",
			descrpition: "Displays game information about all of the users in the game",
			callback:    commandGameDetail,
		},
		"toprankuserinfo": {
			name:        "toprankuserinfo",
			descrpition: "Displays top rank user information",
			callback:    commandTopRankUserInfo,
		},
		"userrankinfo": {
			name:        "userrankinfo",
			descrpition: "Displays user rank information",
			callback:    commandUserRankInfo,
		},
		"characterteaminfo": {
			name:        "characterteaminfo",
			descrpition: "Displays character team information",
			callback:    commandCharacterTeamInfo,
		},
		"statistics": {
			name:        "statistics",
			descrpition: "Displays statistics",
			callback:    commandUserComboStatistics,
		},
		"exit": {
			name:        "exit",
			descrpition: "Exits the program",
			callback:    commandExit,
		},
	}
}
