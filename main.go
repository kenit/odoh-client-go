package main

import (
	"fmt"
	"github.com/kenit/odoh-client-go/commands"
	"github.com/miekg/dns"
)

func main() {
	if handler, err := commands.GetHandler(config.Target, config.Proxy); err == nil {

		dns.HandleFunc(".", Adapt(
			handler,
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
