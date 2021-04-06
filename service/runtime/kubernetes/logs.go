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
// Original source: github.com/micro/go-micro/v3/runtime/kubernetes/logs.go

package kubernetes

import (
	"bufio"
	"strconv"
	"sync"
	"time"

	"github.com/micro/micro/v3/internal/kubernetes/client"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
)

type klog struct {
	client      client.Client
	serviceName string
	options     runtime.LogsOptions
}

func (k *klog) podLogs(podName string, stream *kubeStream) error {
	p := make(map[string]string)
	p["follow"] = "true"

	opts := []client.LogOption{
		client.LogParams(p),
		client.LogNamespace(k.options.Namespace),
	}

	// get the logs for the pod
	body, err := k.client.Log(&client.Resource{
		Name: podName,
		Kind: "pod",
	}, opts...)

	if err != nil {
		stream.err = err
		stream.Stop()
		return err
	}

	s := bufio.NewScanner(body)
	defer body.Close()

	for {
		select {
		case <-stream.stop:
			return stream.Error()
		default:
			if s.Scan() {
				record := runtime.Log{
					Message: s.Text(),
				}

				// send the records to the stream
				// there can be multiple pods doing this
				select {
				case stream.stream <- record:
				case <-stream.stop:
					return stream.Error()
				}
			} else {
				// TODO: is there a blocking call
				// rather than a sleep loop?
				time.Sleep(time.Second)
			}
		}
	}
}

func (k *klog) getMatchingPods() ([]string, error) {
	r := &client.Resource{
		Kind:  "pod",
		Value: new(client.PodList),
	}

	l := make(map[string]string)

	l["name"] = client.Format(k.serviceName)
	// TODO: specify micro:service
	// l["micro"] = "service"

	opts := []client.GetOption{
		client.GetLabels(l),
		client.GetNamespace(k.options.Namespace),
	}

	if err := k.client.Get(r, opts...); err != nil {
		return nil, err
	}

	var matches []string

	for _, p := range r.Value.(*client.PodList).Items {
		// find labels that match the name
		if p.Metadata.Labels["name"] == client.Format(k.serviceName) {
			matches = append(matches, p.Metadata.Name)
		}
	}

	return matches, nil
}

func (k *klog) Read() ([]runtime.Log, error) {
	pods, err := k.getMatchingPods()
	if err != nil {
		return nil, err
	}
	if len(pods) == 0 {
		return nil, errors.NotFound("runtime.logs", "no such service")
	}

	var records []runtime.Log

	for _, pod := range pods {
		logParams := make(map[string]string)

		//if !opts.Since.Equal(time.Time{}) {
		//	logParams["sinceSeconds"] = strconv.Itoa(int(time.Since(opts.Since).Seconds()))
		//}

		if k.options.Count != 0 {
			logParams["tailLines"] = strconv.Itoa(int(k.options.Count))
		}

		if k.options.Stream == true {
			logParams["follow"] = "true"
		}

		opts := []client.LogOption{
			client.LogParams(logParams),
			client.LogNamespace(k.options.Namespace),
		}

		logs, err := k.client.Log(&client.Resource{
			Name: pod,
			Kind: "pod",
		}, opts...)

		if err != nil {
			return nil, err
		}
		defer logs.Close()

		s := bufio.NewScanner(logs)

		for s.Scan() {
			record := runtime.Log{
				Message: s.Text(),
			}
			// record.Metadata["pod"] = pod
			records = append(records, record)
		}
	}

	// sort the records
	// sort.Slice(records, func(i, j int) bool { return records[i].Timestamp.Before(records[j].Timestamp) })

	return records, nil
}

func (k *klog) Stream() (runtime.LogStream, error) {
	// find the matching pods
	pods, err := k.getMatchingPods()
	if err != nil {
		return nil, err
	}

	if len(pods) == 0 {
		return nil, errors.NotFound("runtime.logs", "no such service")
	}

	stream := &kubeStream{
		stream: make(chan runtime.Log),
		stop:   make(chan bool),
	}

	var wg sync.WaitGroup

	// stream from the individual pods
	for _, pod := range pods {
		wg.Add(1)

		go func(podName string) {
			if err := k.podLogs(podName, stream); err != nil {
				logger.Errorf("Error streaming from pod: %v", err)
			}

			wg.Done()
		}(pod)
	}

	go func() {
		// wait until all pod log watchers are done
		wg.Wait()

		// do any cleanup
		stream.Stop()

		// close the stream
		close(stream.stream)
	}()

	return stream, nil
}

// NewLog returns a configured Kubernetes logger
func newLog(c client.Client, serviceName string, opts ...runtime.LogsOption) *klog {
	options := runtime.LogsOptions{
		Namespace: client.DefaultNamespace,
	}
	for _, o := range opts {
		o(&options)
	}

	klog := &klog{
		serviceName: serviceName,
		client:      c,
		options:     options,
	}

	return klog
}
