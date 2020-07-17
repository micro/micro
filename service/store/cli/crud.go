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
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/cmd"
	"github.com/micro/go-micro/v2/store"
	srvstore "github.com/micro/go-micro/v2/store/service"
	"github.com/micro/micro/v2/client/cli/namespace"
	cliutil "github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/client"

	"github.com/pkg/errors"
)

// Read gets something from the store
func Read(ctx *cli.Context) error {
	if ctx.Args().Len() < 1 {
		return errors.New("Key arg is required")
	}
	opts := []store.ReadOption{}
	if ctx.Bool("prefix") {
		opts = append(opts, store.ReadPrefix())
	}

	opts = append(opts, store.ReadFrom(databaseAndTable(ctx)))

	store, err := storeFromContext(ctx)
	if err != nil {
		return err
	}

	records, err := store.Read(ctx.Args().First(), opts...)
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

func databaseAndTable(ctx *cli.Context) (string, string) {
	// default db to namespace
	db, _ := namespace.Get(cliutil.GetEnv(ctx).Name)
	if dbCtx := ctx.String("database"); dbCtx != "" {
		db = dbCtx
	}
	table := ""
	if tableCtx := ctx.String("table"); tableCtx != "" {
		table = tableCtx
	}
	fmt.Printf("DB %s, table %s\n", db, table)
	return db, table
}

func storeFromContext(ctx *cli.Context) (store.Store, error) {
	var st store.Store
	if cliutil.IsLocal(ctx) {
		st = *cmd.DefaultCmd.Options().Store
	} else {
		st = srvstore.NewStore(store.WithClient(client.New(ctx)))
	}
	if err := initStore(ctx, st); err != nil {
		return nil, err
	}
	return st, nil

}

// Write puts something in the store.
func Write(ctx *cli.Context) error {
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
	st, err := storeFromContext(ctx)
	if err != nil {
		return err
	}
	if err := st.Write(record, store.WriteTo(databaseAndTable(ctx))); err != nil {
		return errors.Wrap(err, "couldn't write")
	}
	return nil
}

// List retrieves keys
func List(ctx *cli.Context) error {
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
	opts = append(opts, store.ListFrom(databaseAndTable(ctx)))
	store, err := storeFromContext(ctx)
	if err != nil {
		return err
	}
	keys, err := store.List(opts...)
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

// Delete deletes keys
func Delete(ctx *cli.Context) error {
	if len(ctx.Args().Slice()) == 0 {
		return errors.New("key is required")
	}
	st, err := storeFromContext(ctx)
	if err != nil {
		return err
	}
	if err := st.Delete(ctx.Args().First(), store.DeleteFrom(databaseAndTable(ctx))); err != nil {
		return errors.Wrapf(err, "couldn't delete key %s", ctx.Args().First())
	}
	return nil
}

func initStore(ctx *cli.Context, st store.Store) error {
	opts := []store.Option{}
	db, _ := namespace.Get(cliutil.GetEnv(ctx).Name)
	if dbCtx := ctx.String("database"); dbCtx != "" {
		db = dbCtx
	}
	if db != "" {
		opts = append(opts, store.Database(db))
	}
	if tbl := (ctx.String("table")); tbl != "" {
		opts = append(opts, store.Table(tbl))
	}
	if len(opts) > 0 {
		if err := st.Init(opts...); err != nil {
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
