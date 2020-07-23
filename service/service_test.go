package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/micro/go-micro/v2/client"
	proto "github.com/micro/go-micro/v2/debug/service/proto"
	"github.com/micro/micro/v2/profile"
)

func testShutdown(wg *sync.WaitGroup, cancel func()) {
	// add 1
	wg.Add(1)
	// shutdown the service
	cancel()
	// wait for stop
	wg.Wait()
}

func testService(ctx context.Context, wg *sync.WaitGroup, name string) *Service {
	profile.Test()

	// add self
	wg.Add(1)

	// create service
	return New(
		Name(name),
		AfterStart(func() error {
			wg.Done()
			return nil
		}),
		AfterStop(func() error {
			wg.Done()
			return nil
		}),
	)
}

func testRequest(ctx context.Context, c client.Client, name string) error {
	// test call debug
	req := c.NewRequest(
		name,
		"Debug.Health",
		new(proto.HealthRequest),
	)

	rsp := new(proto.HealthResponse)

	err := c.Call(context.TODO(), req, rsp)
	if err != nil {
		return err
	}

	if rsp.Status != "ok" {
		return errors.New("service response: " + rsp.Status)
	}

	return nil
}

func benchmarkService(b *testing.B, n int, name string) {
	// stop the timer
	b.StopTimer()

	// waitgroup for server start
	var wg sync.WaitGroup

	// cancellation context
	ctx, cancel := context.WithCancel(context.Background())

	// create test server
	service := testService(ctx, &wg, name)

	// start the server
	go func() {
		if err := service.Run(); err != nil {
			b.Fatal(err)
		}
	}()

	// wait for service to start
	wg.Wait()

	// make a test call to warm the cache
	for i := 0; i < 10; i++ {
		if err := testRequest(ctx, service.Client(), name); err != nil {
			b.Fatal(err)
		}
	}

	// start the timer
	b.StartTimer()

	// number of iterations
	for i := 0; i < b.N; i++ {
		// for concurrency
		for j := 0; j < n; j++ {
			wg.Add(1)

			go func() {
				err := testRequest(ctx, service.Client(), name)
				wg.Done()
				if err != nil {
					b.Fatal(err)
				}
			}()
		}

		// wait for test completion
		wg.Wait()
	}

	// stop the timer
	b.StopTimer()

	// shutdown service
	testShutdown(&wg, cancel)
}

func BenchmarkService1(b *testing.B) {
	benchmarkService(b, 1, "test.service.1")
}

func BenchmarkService8(b *testing.B) {
	benchmarkService(b, 8, "test.service.8")
}

func BenchmarkService16(b *testing.B) {
	benchmarkService(b, 16, "test.service.16")
}

func BenchmarkService32(b *testing.B) {
	benchmarkService(b, 32, "test.service.32")
}

func BenchmarkService64(b *testing.B) {
	benchmarkService(b, 64, "test.service.64")
}
