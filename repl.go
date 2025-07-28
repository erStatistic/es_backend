package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kaeba0616/es_backend/internal/erapi"
)

type config struct {
	esapiClient erapi.Client
	currentUser *erapi.User
	users       []erapi.User
	nextgame    *int
}

func startRepl(cfg *config, args ...string) {
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
	callback    func(*config, ...string) error
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
		"exit": {
			name:        "exit",
			descrpition: "Exits the program",
			callback:    commandExit,
		},
	}
}
