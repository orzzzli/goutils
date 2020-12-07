package logger

import (
	"testing"
	"time"
)

var logger *Logger

func TestLogger_Log(t *testing.T) {
	l, err := getLogger()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10000; i++ {
		l.Log("aaaaa", "bbbb", "cccc", "dddd", "eeee", "fffff")
	}
	time.Sleep(3 * time.Second)
}

func BenchmarkLogger_Log(b *testing.B) {
	l, err := getLogger()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		l.Log("aaaaa", "bbbb", "cccc", "dddd", "eeee", "fffff")
	}
}

func getLogger() (*Logger, error) {
	if logger == nil {
		l, err := NewLogger("./", "test_aaa", ',', 1000, 5)
		if err != nil {
			return nil, err
		}
		logger = l
	}
	return logger, nil
}
