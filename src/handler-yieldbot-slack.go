// Take well-formed json from either stdin or an input file and create an attachment notification for Slack
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"github/yieldbot/dhuran"
	"github/yieldbot/dracky"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func set_color(status int) string {
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

func clean_output(output string) string {
	return strings.Split(output, ":")[0]
}

func main() {

	slack_tokenPtr := flag.String("token", "", "the slack integration token")
	channelPtr := flag.String("channel", "monitoring-test", "the channel to post notifications to")
	stdinPtr := flag.Bool("read-stdin", true, "read input from stdin")
	input_filePtr := flag.String("input-file", "", "file to read json in from, check docs for proper format")

	flag.Parse()

	slack_token := *slack_tokenPtr
	channel := *channelPtr
	rd_stdin := *stdinPtr
	input_file := *input_filePtr

	// I don't want to call these if they are not needed
	sensu_event := new(dracky.Sensu_Event)
	user_event := new(dracky.User_Event)

	// YELLOW
	// this is ugly, needs to be a better way to do this
	if slack_token == "" {
		fmt.Print("Please enter a slack integration token")
		os.Exit(1)
	}

	if (rd_stdin == false) && (input_file != "") {
		user_input, err := ioutil.ReadFile(input_file)
		if err != nil {
			dhuran.Check(err)
		}
		err = json.Unmarshal(user_input, &user_event)
		if err != nil {
			dhuran.Check(err)
		}
	} else if (rd_stdin == false) && (input_file == "") {
		fmt.Printf("Please enter a file to read from")
		os.Exit(1)
	} else {
		sensu_event = sensu_event.Acquire_sensu_event()
	}

	api := slack.New(slack_token)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color: set_color(sensu_event.Check.Status),

		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Monitored Instance",
				Value: sensu_event.Acquire_monitored_instance(),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Sensu Client",
				Value: sensu_event.Client.Name,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Check Name",
				Value: dracky.Create_check_name(sensu_event.Check.Name),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Check State",
				Value: dracky.Define_status(sensu_event.Check.Status),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Event Time",
				Value: time.Unix(sensu_event.Check.Issued, 0).Format(time.RFC822Z),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Check State Duration",
				Value: dracky.Define_check_state_duration(),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Check Output",
				Value: clean_output(sensu_event.Check.Output),
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
