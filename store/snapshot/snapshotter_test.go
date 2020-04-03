package snapshot

import (
	"testing"
	"time"

	"github.com/micro/go-micro/v2/store"
)

func TestFileSnapshotter(t *testing.T) {
	f := NewFileSnapshotter(Destination("invalid"))
	if err := f.Init(); err == nil {
		t.Error(err)
	}
	if err := f.Init(Destination("file:///tmp/test-snapshot")); err != nil {
		t.Error(err)
	}

	recordChan, err := f.Start()
	if err != nil {
		t.Error(err)
	}

	for _, td := range testData {
		recordChan <- td
	}
	close(recordChan)
	f.Wait()

	r := NewFileRestorer(Source("invalid"))
	if err := r.Init(); err == nil {
		t.Error(err)
	}
	if err := r.Init(Source("file:///tmp/test-snapshot")); err != nil {
		t.Error(err)
	}

	returnChan, err := r.Start()
	if err != nil {
		t.Error(err)
	}
	var receivedData []*store.Record
	for r := range returnChan {
		println("decoded", r.Key)
		receivedData = append(receivedData, r)
	}

}

var testData = []*store.Record{
	{
		Key:    "foo",
		Value:  []byte(`foo`),
		Expiry: time.Until(time.Now().Add(5 * time.Second)),
	},
	{
		Key:    "bar",
		Value:  []byte(`bar`),
		Expiry: time.Until(time.Now().Add(5 * time.Second)),
	},
	{
		Key:    "baz",
		Value:  []byte(`baz`),
		Expiry: time.Until(time.Now().Add(5 * time.Second)),
	},
}
