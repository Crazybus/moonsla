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
	"github.com/logrusorgru/aurora"
)

func getChannels(api *slack.Client) (channels map[string]string) {
	channels = make(map[string]string)
	chans, _ := api.GetChannels(true)
	for _, c := range chans {
		channels[c.ID] = c.Name
	}
	return channels
}

func getDMs(api *slack.Client, users map[string]string) (channels map[string]string) {
	channels = make(map[string]string)
	chans, _ := api.GetIMChannels()
	for _, c := range chans {
		channels[c.ID] = users[c.User]
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

func getTimeStamp(ts string) (timeStamp time.Time, err error) {
	i, err := strconv.ParseInt(strings.Split(ts, ".")[0], 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	timeStamp = time.Unix(i, 0)
	return timeStamp, nil
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

func formatUrls(msg string) string {
	// Formats slack url into hyperlinks https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
	// Setting MOONSLA_NO_HYPERLINKS=true will disable this for terminals which don't support it yet.

	if os.Getenv("MOONSLA_NO_HYPERLINKS") != "" {
		return msg
	}

	re := regexp.MustCompile("<http.*?>")
	matches := re.FindAllString(msg, -1)
	for _, m := range matches {
		split := strings.Split(m[1:len(m)-1], "|")
		// If this is just a plain url continue since we can't format it
		if len(split) == 1 {
			continue
		}
		url := split[0 : len(split)-1][0]
		title := split[len(split)-1]
		formatted := fmt.Sprintf("\x1b]8;;%s\a%s\x1b]8;;\a", url, title)
		msg = strings.Replace(msg, m, formatted, -1)
	}
	return msg
}

func formatAttachments(attachments []slack.Attachment) string {

	var messages []string

	for _, a := range attachments {

		text := a.Text
		if a.Title != "" {
			text = a.Title + ": " + text
		}

		messages = append(messages, text)
	}
	return strings.Join(messages, "\n")
}

func filterChannel(name string, channels map[string]string, whitelist []string, blacklist []string) (whitelisted bool, cName string) {
	whitelisted = false
	var blacklisted bool = false

	cName, ok := channels[name]
	if ok {
		for _, w := range whitelist {
			if cName == w {
				whitelisted = true
			}
		}
		for _, w := range blacklist {
			if cName == w {
				blacklisted = true
			}
		}
	} else {
		whitelisted = true
		cName = name
	}

	if len(whitelist) == 1 && whitelist[0] == "" {
		whitelisted = true
	}

	if blacklisted {
		return false, cName
	} else {
		return whitelisted, cName
	}
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

	dms := getDMs(api, users)
	fmt.Printf("Found %v DMs\n", len(dms))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	whitelist := strings.Split(os.Getenv("SLACK_CHANNELS"), ",")
	fmt.Printf("Channel whitelist: %v\n", whitelist)

	blacklist := strings.Split(strings.TrimSpace(os.Getenv("SLACK_BLACKLIST_CHANNELS")), ",")
	fmt.Printf("Channel blacklist: %v\n", blacklist)

	for msg := range rtm.IncomingEvents {

		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:

			whitelisted, cName := filterChannel(ev.Channel, channels, whitelist, blacklist)
			var is_dm bool = false

			// Map the users ID to a username if it exists
			uName, ok := users[ev.User]
			if !ok {
				uName = ev.User
			}

			if ev.Username != "" {
				uName = ev.Username
			}

			dmName, present := dms[ev.Channel]
			if present {
				cName = dmName
				is_dm = true
			}

			t, err := getTimeStamp(ev.EventTimestamp)
			var timeStamp string = "00:00:00"
			if err == nil {
				timeStamp = fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
			}

			text := ev.Text

			if len(ev.Attachments) > 0 {
				text = formatAttachments(ev.Attachments)
			}

			msg := formatMentions(text, users)

			msg = formatUrls(msg)
			if !whitelisted {
				continue
			}

			msgC := aurora.Gray(msg)
			if is_dm {
				msgC = aurora.Red(msg)
			}

			fmt.Printf("%v - %v - %v: %v\n", timeStamp, aurora.Green(cName), aurora.Blue(uName), msgC)

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
