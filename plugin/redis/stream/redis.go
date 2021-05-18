package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	"github.com/pkg/errors"
)

type redisStream struct {
	redisClient *redis.Client
}

func NewStream(opts ...Option) (events.Stream, error) {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	rc := redis.NewClient(&redis.Options{
		Addr:      options.Address,
		Username:  options.User,
		Password:  options.Password,
		TLSConfig: options.TLSConfig,
	})
	rs := &redisStream{
		redisClient: rc,
	}
	return rs, nil
}

func (r *redisStream) Publish(topic string, msg interface{}, opts ...events.PublishOption) error {
	// validate the topic
	if len(topic) == 0 {
		return events.ErrMissingTopic
	}
	options := events.PublishOptions{
		Timestamp: time.Now(),
	}
	for _, o := range opts {
		o(&options)
	}

	// encode the message if it's not already encoded
	var payload []byte
	if p, ok := msg.([]byte); ok {
		payload = p
	} else {
		p, err := json.Marshal(msg)
		if err != nil {
			return events.ErrEncodingMessage
		}
		payload = p
	}

	// construct the event
	event := &events.Event{
		ID:        uuid.New().String(),
		Topic:     topic,
		Timestamp: options.Timestamp,
		Metadata:  options.Metadata,
		Payload:   payload,
	}

	// serialize the event to bytes
	bytes, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "Error encoding event")
	}

	return r.redisClient.XAdd(context.Background(), &redis.XAddArgs{
		Stream: event.Topic,
		Values: []string{"event", string(bytes)},
	}).Err()

}

func (r *redisStream) Consume(topic string, opts ...events.ConsumeOption) (<-chan events.Event, error) {
	if len(topic) == 0 {
		return nil, events.ErrMissingTopic
	}

	options := events.ConsumeOptions{}
	for _, o := range opts {
		o(&options)
	}
	group := options.Group
	if len(group) == 0 {
		group = uuid.New().String()
	}
	return r.consumeWithGroup(topic, group, options)
}

func (r *redisStream) consumeWithGroup(topic, group string, options events.ConsumeOptions) (<-chan events.Event, error) {
	lastRead := "$"
	if !options.Offset.IsZero() {
		lastRead = fmt.Sprintf("%d", options.Offset.Unix()*1000)
	}
	if err := r.redisClient.XGroupCreateMkStream(context.Background(), topic, group, lastRead).Err(); err != nil {
		if !strings.HasPrefix(err.Error(), "BUSYGROUP") {
			return nil, err
		}
	}
	consumerName := uuid.New().String()
	ch := make(chan events.Event)
	go func() {
		for {
			res := r.redisClient.XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Group:    group,
				Consumer: consumerName,
				Streams:  []string{topic, ">"},
				Block:    0,
			})

			if err := r.processStreamRes(res, ch, topic, group, options.AutoAck); err != nil {
				close(ch)
				return
			}
		}
	}()
	return ch, nil
}

func (r *redisStream) processStreamRes(res *redis.XStreamSliceCmd, ch chan events.Event, topic, group string, autoAck bool) error {
	sl, err := res.Result()
	if err != nil {
		logger.Errorf("Error reading from stream %s", err)
		return err
	}
	if sl == nil {
		logger.Errorf("No data received from stream")
		return fmt.Errorf("No data received from stream")
	}

	for _, v := range sl[0].Messages {
		evBytes := v.Values["event"]
		var ev events.Event
		good, ok := evBytes.(string)
		if !ok {
			logger.Warnf("Failed to convert to bytes, discarding %s", v.ID)
			r.redisClient.XAck(context.TODO(), topic, group, v.ID)
			continue
		}
		if err := json.Unmarshal([]byte(good), &ev); err != nil {
			logger.Warnf("Failed to unmarshal event, discarding %s %s", err, v.ID)
			r.redisClient.XAck(context.TODO(), topic, group, v.ID)
			continue
		}
		if !autoAck {
			ev.SetAckFunc(func() error {
				return r.redisClient.XAck(context.TODO(), topic, group, v.ID).Err()
			})
		}
		ch <- ev
		if !autoAck {
			continue
		}
		// TODO check for error
		r.redisClient.XAck(context.TODO(), topic, group, v.ID)
	}
	return nil
}
