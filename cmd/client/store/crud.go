package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"github.com/dustin/go-humanize"
	"github.com/micro/micro/v5/cmd/client/util"
	"github.com/micro/micro/v5/service/store"
	"github.com/micro/micro/v5/util/namespace"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// read gets something from the store
func read(ctx *cli.Context) error {
	if err := initStore(ctx); err != nil {
		return err
	}
	if ctx.Args().Len() < 1 {
		return errors.New("Key arg is required")
	}

	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	// get the namespace
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}

	opts := []store.ReadOption{
		store.ReadFrom(ns, ctx.String("table")),
	}
	if ctx.Bool("prefix") {
		opts = append(opts, store.ReadPrefix())
	}
	if ctx.Bool("suffix") {
		opts = append(opts, store.ReadSuffix())
	}
	if ctx.Uint("limit") != 0 {
		opts = append(opts, store.ReadLimit(ctx.Uint("limit")))
	}
	if ctx.Uint("offset") != 0 {
		opts = append(opts, store.ReadLimit(ctx.Uint("offset")))
	}
	if v := ctx.String("order"); len(v) > 0 {
		order := store.OrderAsc
		if v == "desc" {
			order = store.OrderDesc
		}
		opts = append(opts, store.ReadOrder(order))
	}

	records, err := store.DefaultStore.Read(ctx.Args().First(), opts...)
	if err != nil {
		if err.Error() == "not found" {
			return err
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
	default:
		if ctx.Bool("verbose") {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			fmt.Fprintf(w, "%v \t %v \t %v\n", "KEY", "VALUE", "EXPIRY")
			for _, r := range records {
				var key, value, expiry string
				key = r.Key
				if isPrintable(r.Value) {
					value = string(r.Value)
				} else {
					value = fmt.Sprintf("%#x", r.Value[:20])
				}
				if r.Expiry == 0 {
					expiry = "None"
				} else {
					expiry = humanize.Time(time.Now().Add(r.Expiry))
				}
				fmt.Fprintf(w, "%v \t %v \t %v\n", key, value, expiry)
			}
			w.Flush()
			return nil
		}
		for _, r := range records {
			fmt.Println(string(r.Value))
		}
	}
	return nil
}

// write puts something in the store.
func write(ctx *cli.Context) error {
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

	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	// get the namespace
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}

	if err := store.DefaultStore.Write(record, store.WriteTo(ns, ctx.String("table"))); err != nil {
		return errors.Wrap(err, "couldn't write")
	}
	return nil
}

// list retrieves keys
func list(ctx *cli.Context) error {
	if err := initStore(ctx); err != nil {
		return err
	}

	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	// get the namespace
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}

	opts := []store.ListOption{
		store.ListFrom(ns, ctx.String("table")),
	}
	if ctx.Bool("prefix") {
		opts = append(opts, store.ListPrefix(ctx.Args().First()))
	}
	if ctx.Bool("suffix") {
		opts = append(opts, store.ListSuffix(ctx.Args().First()))
	}
	if ctx.Uint("limit") != 0 {
		opts = append(opts, store.ListLimit(ctx.Uint("limit")))
	}
	if ctx.Uint("offset") != 0 {
		opts = append(opts, store.ListLimit(ctx.Uint("offset")))
	}
	if v := ctx.String("order"); len(v) > 0 {
		order := store.OrderAsc
		if v == "desc" {
			order = store.OrderDesc
		}
		opts = append(opts, store.ListOrder(order))
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
	default:
		for _, key := range keys {
			fmt.Println(key)
		}
	}
	return nil
}

// delete deletes keys
func delete(ctx *cli.Context) error {
	if err := initStore(ctx); err != nil {
		return err
	}
	if len(ctx.Args().Slice()) == 0 {
		return errors.New("key is required")
	}

	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	// get the namespace
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}

	if err := store.DefaultStore.Delete(ctx.Args().First(), store.DeleteFrom(ns, ctx.String("table"))); err != nil {
		return errors.Wrapf(err, "couldn't delete key %s", ctx.Args().First())
	}
	return nil
}

func initStore(ctx *cli.Context) error {
	opts := []store.StoreOption{}

	if len(ctx.String("database")) > 0 {
		opts = append(opts, store.Database(ctx.String("database")))
	}
	if len(ctx.String("table")) > 0 {
		opts = append(opts, store.Table(ctx.String("table")))
	}

	if err := store.DefaultStore.Init(opts...); err != nil {
		return errors.Wrap(err, "couldn't reinitialise store with options")
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
