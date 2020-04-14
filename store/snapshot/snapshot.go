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

// Snapshot creates snapshots of a go-micro store
type Snapshot interface {
	// Init validates the Snapshot options and returns an error if they are invalid.
	// Init must be called before the Snapshot is used
	Init(opts ...SnapshotOption) error
	// Start opens a channel that receives *store.Record, adding any incoming records to a backup
	// close() the channel to commit the results.
	Start() (chan<- *store.Record, error)
	// Wait waits for any operations to be committed to underlying storage
	Wait()
}

// SnapshotOptions configure a snapshotter
type SnapshotOptions struct {
	Destination string
}

// SnapshotOption is an individual option
type SnapshotOption func(s *SnapshotOptions)

// Destination is the URL to snapshot to, e.g. file:///path/to/file
func Destination(dest string) SnapshotOption {
	return func(s *SnapshotOptions) {
		s.Destination = dest
	}
}

// FileSnapshot backs up incoming records to a File
type FileSnapshot struct {
	Options SnapshotOptions

	records chan *store.Record
	path    string
	encoder *gob.Encoder
	file    *os.File
	wg      *sync.WaitGroup
}

// NewFileSnapshot returns a FileSnapshot
func NewFileSnapshot(opts ...SnapshotOption) Snapshot {
	f := &FileSnapshot{wg: &sync.WaitGroup{}}
	for _, o := range opts {
		o(&f.Options)
	}
	return f
}

// Init validates the options
func (f *FileSnapshot) Init(opts ...SnapshotOption) error {
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
func (f *FileSnapshot) Start() (chan<- *store.Record, error) {
	if f.records != nil || f.encoder != nil || f.file != nil {
		return nil, errors.New("Snapshot is already in use")
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
func (f *FileSnapshot) Wait() {
	f.wg.Wait()
}

func (f *FileSnapshot) receiveRecords(rec <-chan *store.Record) {
	f.wg.Add(1)
	for {
		r, more := <-rec
		if !more {
			println("Stopping FileSnapshot")
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
