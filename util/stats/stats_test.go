package stats

import (
	"testing"
)

func TestStats(t *testing.T) {
	testCounters := []struct {
		c string
		i []int
	}{
		{
			c: "test",
			i: []int{1, 10, 100},
		},
	}

	s := New()

	if err := s.Start(); err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCounters {
		for _, i := range tc.i {
			s.Record(tc.c, i)
		}
	}

	if err := s.Stop(); err != nil {
		t.Fatal(err)
	}

	if len(s.Counters) == 0 {
		t.Fatalf("stats not recorded, counters are %+v", s.Counters)
	}

	for _, tc := range testCounters {
		if _, ok := s.Counters[0].Status[tc.c]; !ok {
			t.Fatalf("%s counter not found", tc.c)
		}
	}
}
