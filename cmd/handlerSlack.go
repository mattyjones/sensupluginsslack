// Copyright © 2016 Yieldbot <devops@yieldbot.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/op/go-logging"

	"github.com/nlopes/slack"
	"github.com/yieldbot/sensuplugin/sensuhandler"
	"github.com/yieldbot/sensuplugin/sensuutil"

	"github.com/spf13/cobra"
)

var slackToken string
var channelID string

var syslogLog = logging.MustGetLogger("slackHandler")
var stderrLog = logging.MustGetLogger("slackHandler")

// handlerSlackCmd represents the handlerSlack command
var handlerSlackCmd = &cobra.Command{
	Use:   "handlerSlack --token <token> --channel <slack channel>",
	Short: "Post Sensu check results to a slack channel",
	Long: `Read in the Sensu check result and condense the output and post it
	 as a Slack attachment to a given channel`,
	Run: func(cmd *cobra.Command, args []string) {

		syslogBackend, _ := logging.NewSyslogBackend("handlerSlack")
		stderrBackend := logging.NewLogBackend(os.Stderr, "handlerSlack", 0)
		syslogBackendFormatter := logging.NewBackendFormatter(syslogBackend, sensuutil.SyslogFormat)
		stderrBackendFormatter := logging.NewBackendFormatter(stderrBackend, sensuutil.StderrFormat)
		logging.SetBackend(syslogBackendFormatter, stderrBackendFormatter)

		if slackToken == "" {
			syslogLog.Error("Please enter a slack integration token")
			sensuutil.Exit("CONFIGERROR")
		}

		sensuEvent := new(sensuhandler.SensuEvent)
		sensuEvent = sensuEvent.AcquireSensuEvent()

		sensuEnv := sensuhandler.SetSensuEnv()

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
				slack.AttachmentField{
					Title: "Sensu Environment",
					Value: sensuhandler.DefineSensuEnv(sensuEnv.Sensu.Environment),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Uchiwa",
					Value: sensuhandler.AcquireUchiwa("hostname", sensuEnv),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Runbook",
					Value: "",
					Short: true,
				},
			},
		}
		params.Attachments = []slack.Attachment{attachment}
		channelID, timestamp, err := api.PostMessage(channelID, "", params)
		if err != nil {
			syslogLog.Error(err)
		}
		fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

	},
}

func init() {
	RootCmd.AddCommand(handlerSlackCmd)
	handlerSlackCmd.Flags().StringVarP(&slackToken, "token", "", "", "the slack api token")
	handlerSlackCmd.Flags().StringVarP(&channelID, "channel", "", "", "the Slack channel ID")

}
