package snapshot

import (
	"encoding/gob"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/micro/go-micro/v2/store"
	"github.com/pkg/errors"
)

// Snapshotter creates snapshots of a go-micro store
type Snapshotter interface {
	// Init validates the Snapshotter options and returns an error if they are invalid.
	// Init must be called before the Snapshotter is used
	Init(opts ...SnapshotterOption) error
	// Start opens a channel that receives *store.Record, adding any incoming records to a backup
	// close() the channel to commit the results.
	Start() (chan<- *store.Record, error)
	// Wait waits for any operations to be committed to underlying storage
	Wait()
}

// SnapshotterOptions configure a snapshotter
type SnapshotterOptions struct {
	Destination string
}

// SnapshotterOption is an individual option
type SnapshotterOption func(s *SnapshotterOptions)

// Destination is the URL to snapshot to, e.g. file:///path/to/file
func Destination(dest string) SnapshotterOption {
	return func(s *SnapshotterOptions) {
		s.Destination = dest
	}
}

// FileSnapshotter backs up incoming records to a File
type FileSnapshotter struct {
	Options SnapshotterOptions

	records chan *store.Record
	path    string
	encoder *gob.Encoder
	file    *os.File
	wg      *sync.WaitGroup
}

// NewFileSnapshotter returns a FileSnapshotter
func NewFileSnapshotter(opts ...SnapshotterOption) Snapshotter {
	f := &FileSnapshotter{wg: &sync.WaitGroup{}}
	for _, o := range opts {
		o(&f.Options)
	}
	return f
}

// Init validates the options
func (f *FileSnapshotter) Init(opts ...SnapshotterOption) error {
	for _, o := range opts {
		o(&f.Options)
	}
	u, err := url.Parse(f.Options.Destination)
	if err != nil {
		return errors.Wrap(err, "destination is invalid")
	}
	if u.Scheme != "file" {
		return errors.Errorf("unsupported scheme %s (wanted file)", u.Scheme)
	}
	if f.wg == nil {
		f.wg = &sync.WaitGroup{}
	}
	f.path = u.Path
	return nil
}

// Start opens a channel which recieves *store.Record and writes them to storage
func (f *FileSnapshotter) Start() (chan<- *store.Record, error) {
	if f.records != nil || f.encoder != nil || f.file != nil {
		return nil, errors.New("Snapshotter is already in use")
	}
	fi, err := os.OpenFile(f.path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o600)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't open file %s", f.path)
	}
	f.encoder = gob.NewEncoder(fi)
	f.file = fi
	f.records = make(chan *store.Record)
	go f.receiveRecords(f.records)
	return f.records, nil
}

// Wait waits for the snapshotter to commit the backups to persistent storage
func (f *FileSnapshotter) Wait() {
	f.wg.Wait()
}

func (f *FileSnapshotter) receiveRecords(rec <-chan *store.Record) {
	f.wg.Add(1)
	for {
		r, more := <-rec
		if !more {
			println("Stopping FileSnapshotter")
			f.file.Close()
			f.encoder = nil
			f.file = nil
			f.records = nil
			break
		}
		ir := record{
			Key: r.Key,
		}
		if r.Expiry != 0 {
			ir.ExpiresAt = time.Now().Add(r.Expiry)
		}
		ir.Value = make([]byte, len(r.Value))
		copy(ir.Value, r.Value)
		if err := f.encoder.Encode(ir); err != nil {
			// only thing to do here is panic
			panic(errors.Wrap(err, "couldn't write to file"))
		}
		println("encoded", ir.Key)
	}
	f.wg.Done()
}

// record is a store.Record when serialised to persistent storage.
type record struct {
	Key       string
	Value     []byte
	ExpiresAt time.Time
}
