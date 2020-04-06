package snapshot

import (
	"net/url"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	"github.com/pkg/errors"
)

// Snapshot in the entrypoint for micro store snapshot
func Snapshot(ctx *cli.Context) error {
	s := store.DefaultStore
	log := logger.DefaultLogger
	dest := ctx.String("destination")
	var sn Snapshotter

	if len(dest) == 0 {
		return errors.New("destination flag must be set")
	}
	u, err := url.Parse(dest)
	if err != nil {
		return errors.Wrap(err, "destination is invalid")
	}
	switch u.Scheme {
	case "file":
		sn = NewFileSnapshotter(Destination(dest))
	default:
		return errors.Errorf("unsupported destination scheme: %s", u.Scheme)
	}
	err = sn.Init()
	if err != nil {
		return errors.Wrap(err, "failed to initialise the snapshotter")
	}

	log.Logf(logger.InfoLevel, "Snapshotting store %s", s.String())
	recordChan, err := sn.Start()
	if err != nil {
		return errors.Wrap(err, "couldn't start the snapshotter")
	}
	keys, err := s.List()
	if err != nil {
		return errors.Wrap(err, "couldn't List() from store "+s.String())
	}
	log.Logf(logger.DebugLevel, "Snapshotting %d keys", len(keys))

	for _, key := range keys {
		r, err := s.Read(key)
		if err != nil {
			return errors.Wrapf(err, "couldn't read key %s", key)
		}
		if len(r) != 1 {
			return errors.Errorf("reading %s from %s returned 0 records", key, s.String())
		}
		recordChan <- r[0]
	}
	close(recordChan)
	sn.Wait()
	return nil
}

// Restore is the entrypoint for micro store restore
func Restore(ctx *cli.Context) error {
	s := store.DefaultStore
	log := logger.DefaultLogger
	var rs Restorer
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
		rs = NewFileRestorer(Source(source))
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
