// Library for slack handler functions used by the Yieldbot Infrastructure
// teams in sensu.
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package lib

import (
	"strings"
)

// CleanOutput will return a truncated version of the output from a check result
// to hhelp keep the Slack message clear and concise.
func CleanOutput(output string) string {
	return strings.Split(output, ":")[0]
}
