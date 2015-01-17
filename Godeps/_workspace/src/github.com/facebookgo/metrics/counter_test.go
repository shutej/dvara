package metrics_test

import (
	"testing"

	"github.com/facebookgo/metrics"
)

func TestCounterZero(t *testing.T) {
	c := metrics.NewCounter()
	if count := c.Count(); 0 != count {
		t.Errorf("c.Count(): 0 != %v\n", count)
	}
}

func TestCounterInc1(t *testing.T) {
	c := metrics.NewCounter()
	c.Inc(1)
	if count := c.Count(); 1 != count {
		t.Errorf("c.Count(): 1 != %v\n", count)
	}
}

func TestCounterInc12(t *testing.T) {
	c := metrics.NewCounter()
	c.Inc(12)
	if count := c.Count(); 12 != count {
		t.Errorf("c.Count(): 12 != %v\n", count)
	}
}

func TestCounterDec1(t *testing.T) {
	c := metrics.NewCounter()
	c.Dec(1)
	if count := c.Count(); -1 != count {
		t.Errorf("c.Count(): -1 != %v\n", count)
	}
}

func TestCounterDec12(t *testing.T) {
	c := metrics.NewCounter()
	c.Dec(12)
	if count := c.Count(); -12 != count {
		t.Errorf("c.Count(): -12 != %v\n", count)
	}
}

func TestCounterClear(t *testing.T) {
	c := metrics.NewCounter()
	c.Inc(3)
	c.Clear()
	if count := c.Count(); 0 != count {
		t.Errorf("c.Count(): 0 != %v\n", count)
	}
}

func TestCounterIncDec(t *testing.T) {
	c := metrics.NewCounter()
	func() {
		defer c.Inc(2).Dec(1)
	}()
	if count := c.Count(); 1 != count {
		t.Errorf("c.Count(): 1 != %v\n", count)
	}
}
