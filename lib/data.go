// Library for slack handler data used by the Yieldbot Infrastructure
// teams in sensu.
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package lib

// SlackChannels A 1:1 map of the name of a channel and its internal ID.
// All references to channels by slack is done using the ID, but there is currently
// not a programatic way to do this without great expense.
var SlackChannels = map[string]string{
	"monitoring-test": "xxx",
	"devops-alerts":   "xxx",
}
