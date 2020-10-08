// +build integration kind

package test

import(
	"regexp"
)

func statusRunning(service, branch string, statusOutput []byte) bool {
	reg, _ := regexp.Compile(service + "\\s+" + branch + "\\s+\\S+\\s+running")
	return reg.Match(statusOutput)
}