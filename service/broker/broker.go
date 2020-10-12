// Copyright 2020 Asim Aslam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/broker/broker.go

// Package broker is the micro broker
package broker

var (
	// DefaultBroker implementation
	DefaultBroker Broker
)

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Init(...Option) error
	Options() Options
	Address() string
	Connect() error
	Disconnect() error
	Publish(topic string, m *Message, opts ...PublishOption) error
	Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	String() string
}

// Handler is used to process messages via a subscription of a topic.
type Handler func(*Message) error

type ErrorHandler func(*Message, error)

type Message struct {
	Header map[string]string
	Body   []byte
}

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}

// Publish a message to a topic
func Publish(topic string, m *Message) error {
	return DefaultBroker.Publish(topic, m)
}

// Subscribe to a topic
func Subscribe(topic string, h Handler) (Subscriber, error) {
	return DefaultBroker.Subscribe(topic, h)
}
