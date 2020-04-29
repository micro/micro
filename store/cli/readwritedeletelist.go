package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/dustin/go-humanize"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/store"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

// Read gets something from the store
func Read(ctx *cli.Context) error {
	if err := initStore(ctx); err != nil {
		return err
	}
	if ctx.Args().Len() != 1 {
		return errors.New("Key arg is required")
	}
	opts := []store.ReadOption{}
	if ctx.Bool("prefix") {
		opts = append(opts, store.ReadPrefix())
	}

	records, err := store.DefaultStore.Read(ctx.Args().First(), opts...)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.New("not found")
		}
		return errors.Wrapf(err, "Couldn't read %s from store", ctx.Args().First())
	}
	switch ctx.String("output") {
	case "json":
		jsonRecords, err := json.MarshalIndent(records, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed marshalling JSON")
		}
		fmt.Printf("%s\n", string(jsonRecords))
	case "table":
		t := tablewriter.NewWriter(os.Stdout)
		t.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		t.SetCenterSeparator("|")
		t.SetHeader([]string{"Key", "Value", "Expiry"})
		for _, r := range records {
			var key, value, expiry string
			key = r.Key
			if isPrintable(r.Value) {
				value = string(r.Value)
				if len(value) > 50 {
					runes := []rune(value)
					value = string(runes[:50]) + "..."
				}
			} else {
				value = fmt.Sprintf("%#x", r.Value[:20])
			}
			if r.Expiry == 0 {
				expiry = "None"
			} else {
				expiry = humanize.Time(time.Now().Add(r.Expiry))
			}
			t.Append([]string{key, value, expiry})
		}
		t.SetFooter([]string{fmt.Sprintf("Total %d", len(records)), "", ""})
		t.Render()
	default:
		return errors.Errorf("%s is not a valid output format", ctx.String("output"))
	}
	return nil
}

// Write puts something in the store.
func Write(ctx *cli.Context) error {
	if err := initStore(ctx); err != nil {
		return err
	}
	if ctx.Args().Len() < 2 {
		return errors.New("Key and Value args are required")
	}
	record := &store.Record{
		Key:   ctx.Args().First(),
		Value: []byte(strings.Join(ctx.Args().Tail(), " ")),
	}
	if len(ctx.String("expiry")) > 0 {
		d, err := time.ParseDuration(ctx.String("expiry"))
		if err != nil {
			return errors.Wrap(err, "expiry flag is invalid")
		}
		record.Expiry = d
	}
	if err := store.DefaultStore.Write(record); err != nil {
		return errors.Wrap(err, "couldn't write")
	}
	return nil
}

// List retrieves keys
func List(ctx *cli.Context) error {
	if err := initStore(ctx); err != nil {
		return err
	}
	var opts []store.ListOption
	if ctx.Bool("prefix") {
		opts = append(opts, store.ListPrefix(ctx.Args().First()))
	}
	if ctx.Uint("limit") != 0 {
		opts = append(opts, store.ListLimit(ctx.Uint("limit")))
	}
	if ctx.Uint("offset") != 0 {
		opts = append(opts, store.ListLimit(ctx.Uint("offset")))
	}
	keys, err := store.DefaultStore.List(opts...)
	if err != nil {
		return errors.Wrap(err, "couldn't list")
	}
	switch ctx.String("output") {
	case "json":
		jsonRecords, err := json.MarshalIndent(keys, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed marshalling JSON")
		}
		fmt.Printf("%s\n", string(jsonRecords))
	case "table":
		t := tablewriter.NewWriter(os.Stdout)
		t.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		t.SetCenterSeparator("|")
		t.SetHeader([]string{"Key"})
		for _, k := range keys {
			t.Append([]string{k})
		}
		t.SetFooter([]string{fmt.Sprintf("Total %d", len(keys))})
		t.Render()
	default:
		return errors.Errorf("%s is not a valid output format", ctx.String("output"))
	}
	return nil
}

// Delete deletes keys
func Delete(ctx *cli.Context) error {
	if err := initStore(ctx); err != nil {
		return err
	}
	if len(ctx.Args().Slice()) == 0 {
		return errors.New("key is required")
	}
	if err := store.DefaultStore.Delete(ctx.Args().First()); err != nil {
		return errors.Wrapf(err, "couldn't delete key %s", ctx.Args().First())
	}
	return nil
}

func initStore(ctx *cli.Context) error {
	opts := []store.Option{}
	if len(ctx.String("database")) > 0 {
		opts = append(opts, store.Database(ctx.String("database")))
	}
	if len(ctx.String("table")) > 0 {
		opts = append(opts, store.Table(ctx.String("table")))
	}
	if len(opts) > 0 {
		if err := store.DefaultStore.Init(opts...); err != nil {
			return errors.Wrap(err, "couldn't reinitialise store with options")
		}
	}
	return nil
}

func isPrintable(b []byte) bool {
	s := string(b)
	for _, r := range []rune(s) {
		if r == utf8.RuneError {
			return false
		}
	}
	return true
}
