package limiter

import (
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func randomString() []string {
	return []string{"some_string", strconv.Itoa(rand.Int())}
}

func TestNewLimiter(t *testing.T) {
	l, err := NewLimiter(0, 0)
	require.Error(t, err)
	require.Nil(t, l)
	l, err = NewLimiter(10, 0)
	require.Error(t, err)
	require.Nil(t, l)
	l, err = NewLimiter(0, 10)
	require.Error(t, err)
	require.Nil(t, l)
	l, err = NewLimiter(1024, 10)
	require.NoError(t, err)
	require.NotNil(t, l)
}

func TestLimiter_Take(t *testing.T) {
	l, _ := NewLimiter(1024, 10)
	var cnt int
	for i := 0; i < 100000; i++ {
		_, ok := l.Take(randomString()...)
		if !ok {
			cnt++
		}
		_, ok = l.Take()
		if !ok {
			cnt++
		}
	}
	t.Log(cnt)
}

func TestLimiter_Take2(t *testing.T) {
	l, _ := NewLimiter(16, 5)
	wg := sync.WaitGroup{}
	var rejected int32
	const RUN = 1000
	wg.Add(RUN)
	for i := 0; i < RUN; i++ {
		go func() {
			defer wg.Done()
			key := randomString()
			if give, ok := l.Take(key...); ok {
				defer func() {
					give()
					t.Logf("[%v] done", key)
				}()
				t.Logf("[%v] do something", key)
			} else {
				atomic.AddInt32(&rejected, 1)
			}
			if give, ok := l.Take(); ok {
				defer func() {
					give()
					t.Log("[NA] done")
				}()
				t.Log("[NA] do something")
			} else {
				atomic.AddInt32(&rejected, 1)
			}
		}()
	}
	wg.Wait()
	t.Logf("accepted: %d", 2*RUN-rejected)
	t.Logf("rejected: %d", rejected)
}

//BenchmarkLimiter_Take-12    	 2229985	       496.9 ns/op
func BenchmarkLimiter_Take(b *testing.B) {
	l, _ := NewLimiter(1024, 10)
	for i := 0; i < b.N; i++ {
		give, ok := l.Take(randomString()...)
		if ok {
			give()
		}
	}
}
