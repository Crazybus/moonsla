# Moonsla

[![Build Status](https://travis-ci.org/Crazybus/moonsla.svg?branch=master)](https://travis-ci.org/Crazybus/moonsla)

Moonsla is a small tool to display a stream of slack messages in a single view.

It looks something like this

```
10:42:37 - general - Michael Russell: Weird I never knew that slack threads were just normal messages
10:42:55 - random - Someone Else: Sweet Potato!
10:43:37 - general - John Smith: Can people please stop using threads for everything!
```

# Usage

If you don't have one already you will need to generate a [slack API token](https://api.slack.com/custom-integrations/legacy-tokens)

You need to set this to your `SLACK_TOKEN` environment variable
```
export SLACK_TOKEN='xoxp-1231231231232323-123123123123-123123123123123-c91238917239123'
moonsla
```

You can also set `SLACK_CHANNELS` to a comma separated list of channels to filter for
```
export SLACK_CHANNELS='general,random'
```

Or you can instead blacklist some channels:
```
export SLACK_BLACKLIST_CHANNELS='general,random'
```

If you are using a [terminal compatible with hyperlinks](https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda) URLs will be nicely formatted and clickable

```
09:31:58 - D7T4KE0VA - Michael Russell: test
09:31:58 - D7T4KE0VA - slackbot: I searched for that on our Help Center. Perhaps these articles will help:
• Troubleshoot desktop notifications
• Join Slack’s desktop app beta program
• Join Slack’s shared channel beta
```

If your terminal does not support this you can disable it with `MOONSLA_NO_HYPERLINKS=true`

# Why?

I'm not a fan of notifications because they are very intrusive. Instead I used to keep slack always open with the slackbot channel active (so I don't accidentally type shell commands into #general). Whenever I had a spare moment I would then check each slack channel to see if there was anything that needed my attention. This took a bit too long and quite often the message would be bot telling me I had just submitted a pull request

# Future

* Fix channel naming for slackbot and private channels
* Improve formatting of messages so that sub-teams are formatted as expected
* Automatically link to the slack message so it is easy to open up the message from moonsla in slack
* Support multiple slack workspaces
* Add an optional web interface to make things look pretty and allow displaying of images
* Add a nice way to see if a message is from a slack or not
