package main

import (
	"fmt"
	"goratelimiter/redis"
)

func main() {
	client, err := redis.Newclient("localhost:6379")
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do([]string{"PING"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
