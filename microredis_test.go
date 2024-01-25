package microredis

import (
	"log"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	redis := New[bool](time.Second * 10)

	redis.Set("a", true)
	redis.Set("a", true)
	redis.Set("a", true)

	if redis.Size() != 1 {
		log.Println(redis.Size())
	}
	if !*redis.Get("a") {
		log.Println(redis.Size())
	}
}
