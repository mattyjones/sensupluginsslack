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

package sensupluginsslack

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
	"github.com/yieldbot/sensuplugin/sensuhandler"
	"github.com/yieldbot/sensuplugin/sensuutil"
	//"github.com/yieldbot/sensupluginsslack/version"

	"github.com/spf13/cobra"
)

// slack channel to post messages to
var channelID string

// handlerSlackCmd represents the handlerSlack command
var handlerSlackCmd = &cobra.Command{
	Use:   "handlerSlack --token <token> --channel <slack channel>",
	Short: "Post Sensu check results to a slack channel",
	Long: `Read in the Sensu check result and condense the output and post it
	 as a Slack attachment to a given channel`,
	Run: func(sensupluginsslack *cobra.Command, args []string) {

		// Bring in the environmant details
		sensuEnv := new(sensuhandler.EnvDetails)
		sensuEnv = sensuEnv.SetSensuEnv()

		if slackToken == "" {
			syslogLog.WithFields(logrus.Fields{
				"check":      "sensupluginsslack",
				"client":     host,
				//"version":    version.AppVersion(),
				"slackToken": slackToken,
			}).Error(`Please enter a valid slack token`)
			sensuutil.Exit("RUNTIMEERROR")
		}
		// read in the event data from the sensu server
		sensuEvent := new(sensuhandler.SensuEvent)
		sensuEvent = sensuEvent.AcquireSensuEvent()

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
					Value: sensuEvent.Check.Name,
					Short: true,
				},
				slack.AttachmentField{
					Title: "Check State",
					Value: sensuhandler.DefineStatus(sensuEvent.Check.Status),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Event Time (UTC)",
					Value: time.Unix(sensuEvent.Check.Issued, 0).Format(time.RFC3339),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Check State Duration",
					Value: strconv.Itoa(sensuhandler.DefineCheckStateDuration()),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Current Threshold",
					Value: sensuEvent.AcquireThreshold(),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Uchiwa",
					Value: sensuEnv.AcquireUchiwa(sensuEvent.AcquireMonitoredInstance(), sensuEvent.Check.Name),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Playbook",
					Value: sensuEvent.AcquirePlaybook(),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Sensu Environment",
					Value: sensuhandler.DefineSensuEnv(sensuEnv.Sensu.Environment),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Check Event ID",
					Value: sensuEvent.ID,
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
		_, _, err := api.PostMessage(channelID, "", params)
		if err != nil {
			syslogLog.WithFields(logrus.Fields{
				"check":   "sensupluginsslack",
				"client":  host,
				//"version": version.AppVersion(),
				"error":   err,
			}).Error(`Slack attachment could not be sent`)
			sensuutil.Exit("RUNTIMEERROR")
		}
		syslogLog.WithFields(logrus.Fields{
			"check":   "sensupluginsslack",
			"client":  host,
			//"version": version.AppVersion(),
		}).Error(`Slack attachment has been sent`)
		sensuutil.Exit("OK")

	},
}

func init() {
	RootCmd.AddCommand(handlerSlackCmd)
	handlerSlackCmd.Flags().StringVarP(&channelID, "channel", "", "", "the Slack channel ID")

}
