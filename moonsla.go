package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

func getChannels(api *slack.Client) (channels map[string]string) {
	channels = make(map[string]string)
	chans, _ := api.GetChannels(true)
	for _, c := range chans {
		channels[c.ID] = c.Name
	}
	return channels
}

func getUsers(api *slack.Client) (users map[string]string) {
	users = make(map[string]string)
	allUsers, _ := api.GetUsers()
	for _, u := range allUsers {
		users[u.ID] = u.RealName
	}
	return users
}

func formatTimeStamp(ts string) (timeStamp string) {
	fmt.Println(ts)
	i, err := strconv.ParseInt(strings.Split(ts, ".")[0], 10, 64)
	if err != nil {
		panic(err)
	}
	t := time.Unix(i, 0)
	timeStamp = fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
	return timeStamp
}

func formatMentions(msg string, users map[string]string) string {
	re := regexp.MustCompile("<@U.*?>")
	matches := re.FindAllString(msg, -1)
	for _, m := range matches {
		userID := m[2:(len(m) - 1)]
		username, ok := users[userID]
		if ok {
			username = "@" + username
			msg = strings.Replace(msg, m, username, -1)
		}
	}
	return msg
}

func filterChannel(name string, channels map[string]string, whitelist []string) (whitelisted bool, cName string) {
	whitelisted = false

	cName, ok := channels[name]
	if ok {
		for _, w := range whitelist {
			if cName == w {
				whitelisted = true
			}
		}
	} else {
		whitelisted = true
		cName = name
	}

	if len(whitelist) == 0 {
		whitelisted = true
	}

	return whitelisted, cName
}

func main() {

	slackToken, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		fmt.Println("Please set your SLACK_TOKEN")
	}
	api := slack.New(slackToken)

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	channels := getChannels(api)
	fmt.Printf("Found %v channels\n", len(channels))

	users := getUsers(api)
	fmt.Printf("Found %v users\n", len(users))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	whitelist := strings.Split(os.Getenv("SLACK_CHANNELS"), ",")
	fmt.Printf("Channel whitelist: %v\n", whitelist)

	for msg := range rtm.IncomingEvents {

		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:

			// Skip empty messages
			if ev.Text == "" {
				continue
			}

			whitelisted, cName := filterChannel(ev.Channel, channels, whitelist)
			if !whitelisted {
				continue
			}

			// Map the users ID to a username if it exists
			uName, ok := users[ev.User]
			if !ok {
				uName = ev.User
			}

			timeStamp := formatTimeStamp(ev.EventTimestamp)

			msg := formatMentions(ev.Text, users)

			fmt.Printf("%v - %v - %v: %v\n", timeStamp, cName, uName, msg)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
