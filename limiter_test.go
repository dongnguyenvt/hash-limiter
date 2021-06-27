package limiter

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func randomString() string {
	return "some_string_" + strconv.Itoa(rand.Int())
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
		ok, _ := l.Take(randomString(), randomString(), randomString())
		if !ok {
			cnt++
		}
	}
	t.Log(cnt)
	l, _ = NewLimiter(1024, 10)
	cnt = 0
	for i := 0; i < 100000; i++ {
		ok, give := l.Take(randomString(), "some string 2", "some string 3")
		if !ok {
			cnt++
		} else {
			give()
		}
	}
	t.Log(cnt)
}

//BenchmarkLimiter_Take-12    	 2229985	       496.9 ns/op
func BenchmarkLimiter_Take(b *testing.B) {
	l, _ := NewLimiter(1024, 10)
	for i := 0; i < b.N; i++ {
		ok, give := l.Take(randomString(), randomString(), randomString())
		if ok {
			give()
		}
	}
}
