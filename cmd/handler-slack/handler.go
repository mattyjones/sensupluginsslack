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
	"github.com/yieldbot/sensuplugin/sensuhandler"
	"github.com/yieldbot/sensuplugin/sensuutil"
	//"github.com/yieldbot/sensuslack/lib"
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

func main() {

	slackTokenPtr := flag.String("token", "", "the slack integration token")
	channelPtr := flag.String("channel", "monitoring-test", "the channel to post notifications to")

	flag.Parse()

	slackToken := *slackTokenPtr
	channelName := *channelPtr
	channelID := "C09JY7W0P"

	sensuEvent := new(sensuhandler.SensuEvent)
	sensuEvent = sensuEvent.AcquireSensuEvent()

	// YELLOW
	// this is ugly, needs to be a better way to do this
	if slackToken == "" {
		fmt.Print("Please enter a slack integration token")
		sensuutil.Exit("CONFIGERROR")
	}

	fmt.Printf(channelName)

	// for k, v := range lib.SlackChannels {
	// 	if channelName == k {
	// 		channelID = v
	// 	}
	// }

	// if channelID == "000000" {
	// 	fmt.Printf("%v is not mapped, please see the infra team")
	// 	sensuutil.Exit("CONFIGERROR")
	// }

	// api := slack.New(slackToken)
	// // If you set debugging, it will log all requests to the console
	// // Useful when encountering issues
	// // api.SetDebug(true)
	// groups, err := api.GetChannelInfo(channelID)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// 	return
	// }
	// fmt.Printf("%v", groups)

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
				Value: sensuhandler.CreateCheckName(sensuEvent.Check.Name),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Check State",
				Value: sensuhandler.DefineStatus(sensuEvent.Check.Status),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Event Time",
				Value: time.Unix(sensuEvent.Check.Issued, 0).Format(time.RFC3339),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Check State Duration",
				Value: strconv.Itoa(sensuhandler.DefineCheckStateDuration()),
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
	channelID, timestamp, err := api.PostMessage(channelID, "", params)
	if err != nil {
		sensuutil.EHndlr(err)
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
