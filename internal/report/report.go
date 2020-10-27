// Package report contains CLI error reporting and handling and cleanup methods.
// This is mostly required due to the impossibility to catch `os.Exit`s
// (https://stackoverflow.com/questions/39509447/trap-os-exit-in-golang)
// and other funkiness around error paths.
//
// This package currently only tracks the m3o platform calls.
// Please use `1` event values for failure and `0` for success to be consistent
// in our Google Analytics alerts.
package report

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro/micro/v3/internal/helper"
	pb "github.com/micro/micro/v3/proto/alert"
	"github.com/micro/micro/v3/service/client"
	"github.com/urfave/cli/v2"
)

// Error is a helper function to record error events
func Error(ctx *cli.Context, a ...interface{}) {
	val := uint64(1)
	err := TrackEvent(ctx, TrackingData{
		Category: getTrackingCategory(ctx),
		Action:   "error",
		Label:    fmt.Sprint(a...),
		Value:    &val,
	})
	if err != nil {
		fmt.Println(err)
	}
}

// Errorf is a helper function to record error events
func Errorf(ctx *cli.Context, format string, a ...interface{}) {
	val := uint64(1)
	err := TrackEvent(ctx, TrackingData{
		Category: getTrackingCategory(ctx),
		Action:   "error",
		Label:    fmt.Sprintf(format, a...),
		Value:    &val,
	})
	if err != nil {
		fmt.Println(err)
	}
}

// Success is a helper function to record success events
func Success(ctx *cli.Context, a ...interface{}) {
	val := uint64(0)
	err := TrackEvent(ctx, TrackingData{
		Category: getTrackingCategory(ctx),
		Action:   "success",
		Label:    fmt.Sprint(a...),
		Value:    &val,
	})
	if err != nil {
		fmt.Println(err)
	}
}

// Successf is a helper function to record success events
func Successf(ctx *cli.Context, format string, a ...interface{}) {
	val := uint64(0)
	err := TrackEvent(ctx, TrackingData{
		Category: getTrackingCategory(ctx),
		Action:   "success",
		Label:    fmt.Sprintf(format, a...),
		Value:    &val,
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
	Value    *uint64
}

func getTrackingCategory(ctx *cli.Context) string {
	if ctx == nil {
		return "cli"
	}
	command := ctx.Command.Name
	subcommand := helper.Subcommand(ctx)
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
func TrackEvent(ctx *cli.Context, td TrackingData) error {
	sendEvent(ctx, td)
	return nil
}

// send event to alert service
func sendEvent(ctx *cli.Context, td TrackingData) error {
	alertService := pb.NewAlertService("alert", client.DefaultClient)
	val := uint64(0)
	if td.Value != nil {
		val = *td.Value
	}
	_, err := alertService.ReportEvent(context.TODO(), &pb.ReportEventRequest{
		Event: &pb.Event{
			Category: td.Category,
			Action:   td.Action,
			Label:    td.Label,
			Value:    val,
		},
	})
	return err
}
