// Take well-formed json from either stdin or an input file and create an attachment notification for Slack
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package main

import (
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	dracky "github.com/yieldbot/sensu-yieldbot-library/src"
	"os"
	"strconv"
	"strings"
	"time"
)

func setColor(status int) string {
	switch status {
	case 0:
		return "#33CC33"
	case 1:
		return "warning"
	case 2:
		return "#FF0000"
	case 3:
		return "#FF6600"
	default:
		return "#FF6600"
	}
}

func cleanOutput(output string) string {
	return strings.Split(output, ":")[0]
}

func acquireChannelID(channel string) string {
	fmt.Printf("%v", channel)
	return channel
}

func main() {

	slackTokenPtr := flag.String("token", "", "the slack integration token")
	channelPtr := flag.String("channel", "monitoring-test", "the channel to post notifications to")
	// stdinPtr := flag.Bool("read-stdin", true, "read input from stdin")

	flag.Parse()

	slackToken := *slackTokenPtr
	channelName := *channelPtr
  channelID := "000000"
	// rd_stdin := *stdinPtr

	sensuEvent := new(dracky.SensuEvent)
	sensuEvent = sensuEvent.AcquireSensuEvent()

	_ = acquireChannelID(channel)

	// YELLOW
	// this is ugly, needs to be a better way to do this
	if slackToken == "" {
		fmt.Print("Please enter a slack integration token")
		os.Exit(1)
	}

  for k, v := range dracky.SlackChannels {
    if channelName == k {
      channelID = v
    }
  }

  if channelID == "000000" {
    fmt.Printf("%v is not mapped, please see the infra team")
    os.exit(127)
  }

  api_ := slack.New(slackToken)
    // If you set debugging, it will log all requests to the console
    // Useful when encountering issues
    // api.SetDebug(true)
    groups, err := api_.GetChannelInfo(channelID)
    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }
        fmt.Printf("%v", groups)

	api := slack.New(slackToken)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color: setColor(sensuEvent.Check.Status),

		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Monitored Instance",
				Value: sensuEvent.AcquireMonitoredInstance(),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Sensu Client",
				Value: sensuEvent.Client.Name,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Check Name",
				Value: dracky.CreateCheckName(sensuEvent.Check.Name),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Check State",
				Value: dracky.DefineStatus(sensuEvent.Check.Status),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Event Time",
				Value: time.Unix(sensuEvent.Check.Issued, 0).Format(time.RFC3339),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Check State Duration",
				Value: strconv.Itoa(dracky.DefineCheckStateDuration()),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Check Output",
				Value: cleanOutput(sensuEvent.Check.Output),
				Short: true,
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage("C09JY7W0P", "", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
