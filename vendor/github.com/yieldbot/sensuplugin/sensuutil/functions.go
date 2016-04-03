// Library for general purpose functions used by the Yieldbot Infrastructure
// teams in sensu.
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package sensuutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/consul/api"
)

// EHndlr is for generic error handling in all Yieldbot monitoring packages.
func EHndlr(e error) {
	if e != nil {
		fmt.Printf("ERROR: %v", e)
	}
}

// AcquireSensuClient will return the address for the best possible sensu-client
// to connect to.
func AcquireSensuClient() {
	config := api.DefaultConfig()
	config.Address = "core-monitoring-general-rabbitmq-0:8500"
	config.Datacenter = "us-atlanta-1"

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	services, _, _ := client.Catalog().Service("sensu-client", "general", nil)

	for _, s := range services {
		fmt.Println(s)
	}
}

// Exit method for all sensu checks that will print the output and desired exit code
// To use, pass it in the state you want and an optional text you would want outputted with the check.
//Ex. sensuutil.Exit("ok", "Everything is fine")
//    sensuutil.Exit("critical", variable)
// A list of error codes currently supported can be found in common.go
func Exit(args ...interface{}) {
	// YELLOW need to make sure that condition exists
	// panic is bad
	var exitCode int
	output := ""

	if 1 > len(args) {
		panic("Not enough parameters.")
	}

	for i, p := range args {
		switch i {
		case 0: // name
			param, ok := p.(string)
			if !ok {
				panic("1st paramete not type string.")
			}

			for k := range MonitoringErrorCodes {
				if k == strings.ToUpper(param) {
					exitCode = MonitoringErrorCodes[k]
				}
			}

		case 1: // optional text
			param, ok := p.(string)
			if !ok {
				panic("2nd parameter not type string.")
			}
			output = param

		default:
			panic("Incorrect parameters")
		}
	}

	fmt.Printf("%v\n", output)
	os.Exit(exitCode)
}
