package limiter

import (
	"errors"
	"sync/atomic"

	"github.com/segmentio/fasthash/fnv1a"
)

type Give func()

type Limiter interface {
	Take(key ...string) (give Give, ok bool)
}

func NewLimiter(bucket, limit int) (Limiter, error) {
	if bucket <= 0 || limit <= 0 {
		return nil, errors.New("invalid limit")
	}
	return &limiter{
		counters: make([]int32, bucket),
		limit:    int32(limit),
		bucket:   uint32(bucket),
	}, nil
}

type limiter struct {
	bucket   uint32
	counters []int32
	counter  int32
	limit    int32
}

func (l *limiter) Take(key ...string) (Give, bool) {
	if len(key) == 0 {
		if atomic.LoadInt32(&l.counter) > l.limit {
			return nil, false
		}
		atomic.AddInt32(&l.counter, 1)
		return func() { atomic.AddInt32(&l.counter, -1) }, true
	}
	h := fnv1a.Init32
	for _, k := range key {
		h = fnv1a.AddString32(h, k)
	}
	i := h % l.bucket
	if atomic.LoadInt32(&l.counters[i]) > l.limit {
		return nil, false
	}
	atomic.AddInt32(&l.counters[i], 1)
	return func() {
		atomic.AddInt32(&l.counters[i], -1)
	}, true
}
