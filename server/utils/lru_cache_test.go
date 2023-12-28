package utils

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestLruCache(t *testing.T) {
	lru := LruCacheConstructor(2, time.Minute, 2*time.Minute)
	results := []string{"a", "a", "a", "a"}

	lru.Put("1", "1")
	lru.Put("2", "2")
	results[0] = lru.Get("1") // 1,2

	lru.Put("3", "3")          // 3,1
	results[1] = lru.Get("1")  // 1,3
	results[2] = lru.Get("-2") // 1,3

	lru.Put("4", "4")         // 4,1
	results[3] = lru.Get("4") // 4,1

	got := ArrayToString(results, ",")
	want := "1,1,,4"
	if got != want {
		t.Errorf("Got %s, wanted %s", got, want)
	}
}

func TestLruCacheTtl1(t *testing.T) {
	lru := LruCacheConstructor(2, time.Second, 2*time.Second)

	results := []string{"a", "a", "a", "a", "a"}

	lru.Put("1", "1")
	lru.Put("2", "2")
	results[0] = lru.Get("2") // 2,1

	lru.Put("3", "3")         // 3,2
	results[1] = lru.Get("1") // 3,2
	results[2] = lru.Get("3") // 3,2

	time.Sleep(3 * time.Second)

	results[3] = lru.Get("3")
	results[4] = lru.Get("2")

	got := ArrayToString(results, ",")
	want := "2,,3,,"
	if got != want {
		t.Errorf("Got %s, wanted %s", got, want)
	}
}

func TestLruCacheTtl2(t *testing.T) {
	lru := LruCacheConstructor(3, 3*time.Second, 2*time.Second)

	results := []string{"a", "a", "a", "a", "a", "a", "a"}

	lru.Put("1", "1")
	lru.Put("2", "2")
	results[0] = lru.Get("1") // 1,2

	lru.Put("3", "3")         // 3,1,2
	results[1] = lru.Get("3") // 3,1,2
	results[2] = lru.Get("2") // 2,3,1

	time.Sleep(3 * time.Second)

	lru.Put("4", "4")         // 4,2,3
	results[3] = lru.Get("2") // 2,4,3
	results[4] = lru.Get("1") // 4,2,3

	time.Sleep(3 * time.Second) // 4

	results[5] = lru.Get("4") // 4
	results[6] = lru.Get("2") // 4

	got := ArrayToString(results, ",")
	want := "1,3,2,2,,4,"
	if got != want {
		t.Errorf("Got %s, wanted %s", got, want)
	}
}

func ArrayToString(a []string, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
