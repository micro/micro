package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	attempts    map[string]int
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
		attempts:    map[string]int{},
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
		Stream: fmt.Sprintf("stream-%s", event.Topic),
		Values: map[string]interface{}{"event": string(bytes), "attempt": 1},
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
	topic = fmt.Sprintf("stream-%s", topic)
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

			if err := r.processStreamRes(res, ch, topic, group, options.AutoAck, options.RetryLimit); err != nil {
				close(ch)
				return
			}
		}
	}()
	return ch, nil
}

func (r *redisStream) processStreamRes(res *redis.XStreamSliceCmd, ch chan events.Event, topic, group string, autoAck bool, retryLimit int) error {
	topic = fmt.Sprintf("stream-%s", topic)
	sl, err := res.Result()
	if err != nil {
		logger.Errorf("Error reading from stream %s", err)
		return err
	}
	if sl == nil {
		logger.Errorf("No data received from stream")
		return fmt.Errorf("no data received from stream")
	}

	for _, v := range sl[0].Messages {
		evBytes := v.Values["event"]
		var ev events.Event
		bStr, ok := evBytes.(string)
		if !ok {
			logger.Warnf("Failed to convert to bytes, discarding %s", v.ID)
			r.redisClient.XAck(context.Background(), topic, group, v.ID)
			continue
		}
		if err := json.Unmarshal([]byte(bStr), &ev); err != nil {
			logger.Warnf("Failed to unmarshal event, discarding %s %s", err, v.ID)
			r.redisClient.XAck(context.Background(), topic, group, v.ID)
			continue
		}
		r.attempts[v.ID], _ = strconv.Atoi(v.Values["attempt"].(string))

		if !autoAck {
			ev.SetAckFunc(func() error {
				delete(r.attempts, v.ID)
				return r.redisClient.XAck(context.Background(), topic, group, v.ID).Err()
			})
			ev.SetNackFunc(func() error {
				// no way to nack a message. Best you can do is to ack and readd
				if err := r.redisClient.XAck(context.Background(), topic, group, v.ID).Err(); err != nil {
					return err
				}
				attempt := r.attempts[v.ID]
				if retryLimit > 0 && attempt > retryLimit {
					// don't readd
					delete(r.attempts, v.ID)
					return nil
				}
				bytes, err := json.Marshal(ev)
				if err != nil {
					return errors.Wrap(err, "Error encoding event")
				}
				return r.redisClient.XAdd(context.Background(), &redis.XAddArgs{
					Stream: fmt.Sprintf("stream-%s", ev.Topic),
					Values: map[string]interface{}{"event": string(bytes), "attempt": attempt + 1},
				}).Err()
			})
		}
		ch <- ev
		if !autoAck {
			continue
		}
		// TODO check for error
		r.redisClient.XAck(context.Background(), topic, group, v.ID)
	}
	return nil
}
