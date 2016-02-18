// Take well-formed json from either stdin or an input file and create an attachment notification for Slack
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/nlopes/slack"
	"github.com/yieldbot/sensuplugin/sensuhandler"
	"github.com/yieldbot/sensuplugin/sensuutil"
	"os"
	"strconv"
	"time"
)

func main() {

	var slackToken string
	var channelID string
	var debug bool

	app := cli.NewApp()
	app.Name = "handler-slack"
	app.Usage = "Send notifications to a given Slack channel as an attachment"
	app.Action = func(c *cli.Context) {

		sensuEvent := new(sensuhandler.SensuEvent)
		sensuEvent = sensuEvent.AcquireSensuEvent()

		if slackToken == "" {
			fmt.Print("Please enter a slack integration token")
			sensuutil.Exit("CONFIGERROR")
		}

		if debug {
			fmt.Printf("For the value of the flags please see the handler configuration file")
			sensuutil.Exit("ok")
		}

		// This is done with an api token not an incoming webhook to a specific channel
		api := slack.New(slackToken)
		params := slack.PostMessageParameters{}
		// Build an attachment message for sending to the specified slack channel
		attachment := slack.Attachment{
			Color: sensuhandler.SetColor(sensuEvent.Check.Status),

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
					Value: sensuhandler.CleanOutput(sensuEvent.Check.Output),
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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "token, t",
			Value:       "",
			Usage:       "the slack integration token",
			EnvVar:      "SLACK_TOKEN",
			Destination: &slackToken,
		},

		cli.StringFlag{
			Name:        "channel, c",
			Value:       "",
			Usage:       "The channel ID that you wish to send this to",
			EnvVar:      "SLACK_CHANNEL_ID",
			Destination: &channelID,
		},

		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Set this to print debugging information. No notifications will be sent",
			Destination: &debug,
		},
	}
	app.Run(os.Args)
}
