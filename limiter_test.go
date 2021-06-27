package limiter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tjarratt/babble"
)

var babbler = babble.NewBabbler()

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
		ok, _ := l.Take(babbler.Babble(), babbler.Babble(), babbler.Babble())
		if !ok {
			cnt++
		}
	}
	t.Log(cnt)
	l, _ = NewLimiter(1024, 10)
	cnt = 0
	for i := 0; i < 100000; i++ {
		ok, give := l.Take(babbler.Babble(), babbler.Babble(), babbler.Babble())
		if !ok {
			cnt++
		} else {
			give()
		}
	}
	t.Log(cnt)
}

func BenchmarkLimiter_Take(b *testing.B) {
	l, _ := NewLimiter(1024, 10)
	for i := 0; i < b.N; i++ {
		ok, give := l.Take(babbler.Babble(), babbler.Babble(), babbler.Babble())
		if ok {
			give()
		}
	}
}
