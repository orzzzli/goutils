package snowflake

import (
	"strconv"
	"testing"
)

var snow *SnowFlake
var testMap map[int64]int
var testMap2 map[int64]int

func TestSnowFlake_Get(t *testing.T) {
	s, err := getSnowFlake()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s.Get())
}

func TestSnowFlake_GetFromRing(t *testing.T) {
	s, err := getSnowFlake()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s.GetFromRing())
}

func BenchmarkSnowFlake_Get(b *testing.B) {
	s, err := getSnowFlake()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		re := s.Get()
		_, ok := testMap[re]
		if ok {
			b.Fatal("repeat." + strconv.Itoa(i) + " " + strconv.FormatInt(re, 10))
		}
		testMap[re] = i
	}
}

func BenchmarkSnowFlake_GetFromRing(b *testing.B) {
	s, err := getSnowFlake()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		re := s.GetFromRing()
		_, ok := testMap2[re]
		if ok {
			b.Fatal("repeat." + strconv.Itoa(i) + " " + strconv.FormatInt(re, 10))
		}
		testMap2[re] = i
	}
}

func getSnowFlake() (*SnowFlake, error) {
	if snow == nil {
		s, err := NewSnowFlake(1606752000000, 1, 9, 12, 1, 1)
		if err != nil {
			return nil, err
		}
		snow = s
	}
	if testMap == nil {
		testMap = make(map[int64]int)
	}
	if testMap2 == nil {
		testMap2 = make(map[int64]int)
	}
	return snow, nil
}
