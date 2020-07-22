// Package report contains CLI error reporting and handling and cleanup methods.
// This is mostly required due to the impossibility to catch `os.Exit`s
// (https://stackoverflow.com/questions/39509447/trap-os-exit-in-golang)
// and other funkiness around error paths.
//
// This package currently only tracks the m3o platform calls.
package report

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/helper"
)

const (
	gaPropertyID = "UA-70478210-6"
	envToTrack   = "staging"
)

// Error is a helper function: it prints an error, records events an exits the CLI
func Error(ctx *cli.Context, a ...interface{}) {
	env := util.GetEnv(ctx)
	if env.Name != envToTrack {
		return
	}
	err := TrackEvent(TrackingData{
		Category: getTrackingCategory(ctx),
		Action:   "error",
		Label:    fmt.Sprint(a...),
	})
	if err != nil {
		fmt.Println(err)
	}
}

// Errorf is a helper function: it prints an error, records events an exits the CLI
func Errorf(ctx *cli.Context, format string, a ...interface{}) {
	env := util.GetEnv(ctx)
	if env.Name != envToTrack {
		return
	}
	err := TrackEvent(TrackingData{
		Category: getTrackingCategory(ctx),
		Action:   "error",
		Label:    fmt.Sprintf(format, a...),
	})
	if err != nil {
		fmt.Println(err)
	}
}

type TrackingData struct {
	Category string
	Action   string
	Label    string
	UserID   string
	Value    *uint
}

func getTrackingCategory(ctx *cli.Context) string {
	command := helper.MicroCommand(ctx)
	subcommand := helper.MicroSubcommand(ctx)
	if len(strings.TrimSpace(subcommand)) == 0 {
		return command
	}
	return strings.Join([]string{command, subcommand}, "/")
}

// TrackEvent records an event on google analytics
// For details consult https://support.google.com/analytics/answer/1033068?hl=en
//
// Example:
// Category: "Videos"
// Action: "Downloaded"
// Label: "Gone With the Wind"
func TrackEvent(td TrackingData) error {
	if gaPropertyID == "" {
		return errors.New("analytics: GA_TRACKING_ID environment variable is missing")
	}
	if td.Category == "" || td.Action == "" {
		return errors.New("analytics: category and action are required")
	}

	cid := td.UserID
	if len(cid) == 0 {
		// GA does not seem to accept events without user id so we generate a UUID
		cid = uuid.New().String()
	}
	v := url.Values{
		"v":   {"1"},
		"tid": {gaPropertyID},
		// Anonymously identifies a particular user. See the parameter guide for
		// details:
		// https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#cid
		//
		// Depending on your application, this might want to be associated with the
		// user in a cookie.
		"cid": {cid},
		"t":   {"event"},
		"ec":  {td.Category},
		"ea":  {td.Action},
		"ua":  {"cli"},
	}

	if td.Label != "" {
		v.Set("el", td.Label)
	}

	if td.Value != nil {
		v.Set("ev", fmt.Sprintf("%d", *td.Value))
	}

	// NOTE: Google Analytics returns a 200, even if the request is malformed.
	_, err := http.PostForm("https://www.google-analytics.com/collect", v)
	return err
}
