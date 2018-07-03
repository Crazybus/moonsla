package main

import "testing"

func TestFormatTimeStamp(t *testing.T) {
	var tests = []struct {
		description string
		timeStamp   string
		want        string
	}{
		{
			"Convert timestamp to something human readable",
			"1530593277.000080",
			"06:47:57",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := formatTimeStamp(test.timeStamp)
			want := test.want

			if got != want {
				t.Errorf("got '%s' want '%s'", got, want)
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
