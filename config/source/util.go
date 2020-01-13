package mucp

import (
	"time"

	"github.com/micro/go-micro/config/source"
	proto "github.com/micro/micro/config/proto/config"
)

func toChangeSet(c *proto.ChangeSet) *source.ChangeSet {
	return &source.ChangeSet{
		Data:      c.Data,
		Checksum:  c.Checksum,
		Format:    c.Format,
		Timestamp: time.Unix(c.Timestamp, 0),
		Source:    c.Source,
	}
}
