package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kenit/odoh-client-go/commands"
	"github.com/miekg/dns"
)

func main() {
	var adapters = []Adapter{}

	if config.Redis != nil {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
			Password: "", // no password set
			DB:       config.Redis.DB, // use default DB
		})

		adapters = append(adapters, RedisCache(rdb))
	}

	if handler, err := commands.GetHandler(config.Target, config.Proxy); err == nil {

		dns.HandleFunc(".", Adapt(handler,
			adapters ...
		))

		addr := fmt.Sprintf("%s:%d", config.Listen.Host, config.Listen.Port)
		server := &dns.Server{
			Addr: addr,
			Net:  "udp",
		}

		fmt.Printf("Server will listen on: %s\n ", addr)

		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to start server: %s\n ", err.Error())
		}
	} else {
		fmt.Printf("Failed to start server: %s\n ", err.Error())
	}
}
