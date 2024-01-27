package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
)

var (
	btoken  string
	chatID  int
	mgs     string
	file    string
	passive bool
	active  bool
)

func init() {
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("sendTele, my bugbounty telegramBot notification.")
	flagSet.CreateGroup("input", "INPUT",
		flagSet.StringVar(&btoken, "t", "", "Bot token"),
		flagSet.IntVar(&chatID, "c", 0, "Chat ID"),
		flagSet.StringVar(&mgs, "m", "", "Message to be sent"),
		flagSet.StringVar(&file, "f", "", "To get the len of a file"),
	)
	flagSet.CreateGroup("notify","NOTIFY",
		flagSet.BoolVarP(&passive, "p", "passive", false, "Notify me if I'm performing a passive scan."),
		flagSet.BoolVarP(&active, "a", "active", false, "Notify me if I'm performing a active scan."),
	)

	_ = flagSet.Parse()
}

func sendMessage(btoken string, chatID int64, message string) {
	bot, err := tgbotapi.NewBotAPI(btoken)
	if err != nil {
		gologger.Fatal().Msgf("Error initializing Telegram bot: %v", err)
	}

	msg := tgbotapi.NewMessage(chatID, message)

	_, err = bot.Send(msg)
	if err != nil {
		gologger.Error().Msgf("%v", err)
	}
}

func readFileLen(filepath string) int {
	file, err := os.Open(filepath)
	if err != nil {
		errFatalE(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		errFatalE(err)
	}

	return len(lines)
}

func teleCheck(btoken string, chatID int) bool {
	if btoken != "" && chatID != 0 {
		return true
	} else {
		return false
	}
}

func errFatal(err string) {
	gologger.Fatal().Msg(err)
}

func errFatalE(err error) {
	gologger.Fatal().Msgf("%v", err)
}

func main() {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	tch := teleCheck(btoken, chatID)
	if tch && file == "" && mgs == "" {
		errFatal("Please specify file/message using -m/-f or do -h for help")
	}

	if file != "" && strings.HasSuffix(file, ".txt") {
		fileLen := readFileLen(file)
		if tch {
			if passive {
				msgs := fmt.Sprintf("[INFO] %s \nPassive Subdomain Enumeration Started\n[-] Domain: %d", currentTime, fileLen)
				sendMessage(btoken, int64(chatID), msgs)
			} else if active {
				msgs := fmt.Sprintf("[INFO] %s \nActive Subdomain Enumeration Started\n[-] Domain: %d", currentTime, fileLen)
				sendMessage(btoken, int64(chatID), msgs)
			} else {
				errFatal("Please specify Passive/Active using -p/-a or do -h for help")
			}
		} else {
			errFatal("Please specify ChatID and Token using -c & -t or do -h for help")
		}
	}

	if tch {
		if mgs != "" && strings.HasPrefix(mgs, ".") {
			msgs := fmt.Sprintf("[INFO] %s: %s", currentTime, strings.ReplaceAll(mgs, "{nl}", "\n"))
			sendMessage(btoken, int64(chatID), msgs)
		} else if mgs != "" {
			msgs := fmt.Sprintf("[INFO] %s - %s now", currentTime, mgs)
			sendMessage(btoken, int64(chatID), msgs)
		}
	} else {
		errFatal("Please specify ChatID and Token using -c & -t or do -h for help")
	}
}
