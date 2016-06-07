## sensupluginsslack

## Commands
 * handlerSlack

## Usage

### handlerSlack
Send a subset of the Sensu check result to a given slack channel

Ex. `sensupluginsslack --token --channel`

In order to receive alerts in Slack you will need two data points from Slack, an api token and the channel_id for the channel you want to post to.

#### Token

The [Slack Web Api](https://api.slack.com/web) page will get you started for generating a token

#### Channel ID

The channel_id is the numeric representation of the channel according to Slack. This is not dependent on the channel name so it is a 1 time cost of getting it and then you can change the name of the channel that is represents. There is only one way to obtain the channel_id and that is via an [api call](https://api.slack.com/methods/channels.list). This will return a list of all channels in the org and in that list will be the channel_id. You will need your api token to obtain this.

## Installation

1. godep go build -o bin/sensupluginsslack
1. chmod +x sensupluginsslack
1. cp sensupluginsslack /usr/local/bin

## Notes

This handler will require a json file to provide additional values. The best way to create this file is with chef or another configuration management tool as outlined below.
The file location and name default to `/etc/sensu/conf.d/monitoring_infra.json`.

```json
{
    "sensu": {
        "environment": "<%=node.chef_environment%>",
        "fqdn": "<%=node.fqdn%>",
        "hostname": "<%=node.hostname%>",
        "consul": {
          "tags": "<%= node['devops_monitoring']['cluster'] %>",
          "datacenter": "<%= node['sensu2']['uchiwa']['dc_name'] %>"
        }
    }
}
```
