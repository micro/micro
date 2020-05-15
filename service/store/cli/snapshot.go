package cli

import (
	"net/url"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/store/snapshot"
	"github.com/pkg/errors"
)

// Snapshot in the entrypoint for micro store snapshot
func Snapshot(ctx *cli.Context) error {
	s, err := makeStore(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't construct a store")
	}
	log := logger.DefaultLogger
	dest := ctx.String("destination")
	var sn snapshot.Snapshot

	if len(dest) == 0 {
		return errors.New("destination flag must be set")
	}
	u, err := url.Parse(dest)
	if err != nil {
		return errors.Wrap(err, "destination is invalid")
	}
	switch u.Scheme {
	case "file":
		sn = snapshot.NewFileSnapshot(snapshot.Destination(dest))
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
