package main

import (
	"testing"

	"github.com/nlopes/slack"
)

func TestGetTimeStamp(t *testing.T) {
	var tests = []struct {
		description string
		timeStamp   string
		want        int64
	}{
		{
			"Convert timestamp to something human readable",
			"1530593277.000080",
			1530593277,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ts := getTimeStamp(test.timeStamp)
			got := ts.Unix()

			want := test.want

			if got != want {
				t.Errorf("got '%v' want '%v'", got, want)
			}
		})
	}
}

func TestFilterChannel(t *testing.T) {
	var tests = []struct {
		description string
		id          string
		channels    map[string]string
		whitelist   []string
		name        string
		whitelisted bool
	}{
		{
			"Channel that is whitelisted",
			"12345",
			map[string]string{
				"12345": "channel-name",
			},
			[]string{
				"channel-name",
			},
			"channel-name",
			true,
		},
		{
			"Channel that is not whitelisted",
			"12344",
			map[string]string{
				"12345": "channel-name",
				"12344": "spam-channel",
			},
			[]string{
				"channel-name",
			},
			"spam-channel",
			false,
		},
		{
			"Channel that is not in the channels list",
			"123",
			map[string]string{
				"12345": "channel-name",
				"12344": "spam-channel",
			},
			[]string{
				"channel-name",
			},
			"123",
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			whitelisted, name := filterChannel(test.id, test.channels, test.whitelist)

			if name != test.name {
				t.Errorf("got '%s' want '%s'", name, test.name)
			}
			if whitelisted != test.whitelisted {
				t.Errorf("got '%v' want '%v'", whitelisted, test.whitelisted)
			}
		})
	}
}

func TestFormatMentions(t *testing.T) {
	var tests = []struct {
		description string
		message     string
		users       map[string]string
		want        string
	}{
		{
			"Replace mentions on message",
			"hello <@U1234> how are you?",
			map[string]string{
				"U1234": "crazybus",
			},
			"hello @crazybus how are you?",
		},
		{
			"Replace multiple mentions in a message",
			"hello <@U1234> have you met <@U321> I think you would like them <@U1234>?",
			map[string]string{
				"U1234": "crazybus",
				"U321":  "notcrazybus",
			},
			"hello @crazybus have you met @notcrazybus I think you would like them @crazybus?",
		},
		{
			"Don't replace anything if there are no mentions",
			"hi",
			map[string]string{
				"U1234": "crazybus",
				"U321":  "notcrazybus",
			},
			"hi",
		},
		{
			"Leave the id if the user can't be found",
			"hi <@U999>!",
			map[string]string{
				"U1234": "crazybus",
				"U321":  "notcrazybus",
			},
			"hi <@U999>!",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := formatMentions(test.message, test.users)
			want := test.want

			if got != want {
				t.Errorf("got '%s' want '%s'", got, want)
			}
		})
	}
}

func TestFormatAttachments(t *testing.T) {
	var tests = []struct {
		description string
		attachments []slack.Attachment
		want        string
	}{
		{
			"Print attachment as single line",
			[]slack.Attachment{slack.Attachment{
				Title: "",
				Text:  "test message",
			}},
			"test message",
		},
		{
			"Print attachment with title",
			[]slack.Attachment{slack.Attachment{
				Title: "title",
				Text:  "test message",
			}},
			"title: test message",
		},
		{
			"Print multie attachments",
			[]slack.Attachment{
				slack.Attachment{
					Title: "",
					Text:  "first message",
				},
				slack.Attachment{
					Title: "",
					Text:  "second message",
				}},
			"first message\nsecond message",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

			got := formatAttachments(test.attachments)

			want := test.want

			if got != want {
				t.Errorf("got '%v' want '%v'", got, want)
			}
		})
	}
}

func TestFormatUrls(t *testing.T) {
	var tests = []struct {
		description string
		message     string
		want        string
	}{
		{
			"Message with no urls",
			"test message",
			"test message",
		},
		{
			"Message with a url",
			"hello <http://google.com|test> world",
			"hello \x1b]8;;http://google.com\atest\x1b]8;;\a world",
		},
		{
			"Message with multiple urls",
			"hello <http://google.com|test> world how <https://google.com|are you>",
			"hello \x1b]8;;http://google.com\atest\x1b]8;;\a world how \x1b]8;;https://google.com\aare you\x1b]8;;\a",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

			got := formatUrls(test.message)

			want := test.want

			if got != want {
				t.Errorf("got '%v' want '%v'", got, want)
			}
		})
	}
}
