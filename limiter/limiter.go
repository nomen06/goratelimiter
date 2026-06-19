package main

import (
	"fmt"
	"goratelimiter/redis"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Limiter struct {
	client *redis.Client
	limit  int
	window int
}

func NewLimiter(c *redis.Client, limit int, window int) *Limiter {
	return &Limiter{client: c, limit: limit, window: window}
}
func (l *Limiter) Allow(ip string) bool {
	key := ("client:" + ip)
	resp, err := l.client.Do([]string{"INCR", key})
	resp = strings.TrimPrefix(resp, ":")
	resp = strings.TrimSpace(resp)
	count, err2 := strconv.Atoi(resp)
	fmt.Println("DEBUG: raw resp =", resp, "| count =", count, "| err =", err, "| err2 =", err2)
	if err != nil {
		return false
	}
	if count == 1 {
		_, err := l.client.Do([]string{"EXPIRE", key, strconv.Itoa(l.window)})
		if err != nil {
			return false
		}
	}

	if count > l.limit {
		return false
	}
	return true
}

func handler(limiter *Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if !limiter.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		w.Write([]byte("Hello"))
	}
}

func main() {
	client, _ := redis.Newclient("localhost:6379")
	limiter := NewLimiter(client, 1000, 60)
	http.HandleFunc("/", handler(limiter))
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}
