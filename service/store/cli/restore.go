package cli

import (
	"net/url"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/store/snapshot"
	"github.com/pkg/errors"
)

// Restore is the entrypoint for micro store restore
func Restore(ctx *cli.Context) error {
	s, err := makeStore(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't construct a store")
	}
	log := logger.DefaultLogger
	var rs snapshot.Restore
	source := ctx.String("source")

	if len(source) == 0 {
		return errors.New("source flag must be set")
	}
	u, err := url.Parse(source)
	if err != nil {
		return errors.Wrap(err, "source is invalid")
	}
	switch u.Scheme {
	case "file":
		rs = snapshot.NewFileRestore(snapshot.Source(source))
	default:
		return errors.Errorf("unsupported source scheme: %s", u.Scheme)
	}

	err = rs.Init()
	if err != nil {
		return errors.Wrap(err, "failed to initialise the restorer")
	}

	recordChan, err := rs.Start()
	if err != nil {
		return errors.Wrap(err, "couldn't start the restorer")
	}
	counter := uint64(0)
	for r := range recordChan {
		err := s.Write(r)
		if err != nil {
			log.Logf(logger.ErrorLevel, "couldn't write key %s to store %s", r.Key, s.String())
		} else {
			counter++
		}
	}
	log.Logf(logger.DebugLevel, "Restored %d records", counter)
	return nil
}
