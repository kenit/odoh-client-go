package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
	"github.com/go-redis/redis/v8"
	"github.com/kenit/odoh-client-go/commands"
	"github.com/miekg/dns"
)

func RedisCache(rdb *redis.Client) Adapter {

	ctx := context.Background()
	return func(next commands.DnsHandler) commands.DnsHandler {
		if rdb == nil{
			return next
		}

		return func(w dns.ResponseWriter, r *dns.Msg) {
				fmt.Println(len(r.Question))
				question := r.Question[0].String()
				key := fmt.Sprintf("%x", sha256.Sum256([]byte(question)))
				fmt.Println(key)
				answer, err := rdb.Get(ctx, key).Result()
				if err == nil {
					fmt.Println("HIT")
					msg := &dns.Msg{}
					if msg.Unpack([]byte(answer)) == nil {
						msg.Id = r.Id
						if err := w.WriteMsg(msg); err == nil {
							return
						}
					}
				} else {
					fmt.Println("NO HIT")
					fmt.Println(err)
					wrappedRW := &WrappedResponseWriter{
						originResponseWriter: w,
						beforeWriteMsg: func(m *dns.Msg) error {
							if msg, err := m.Pack(); err != nil {
								return err
							} else {
								rdb.Set(ctx, key, msg, 24 * time.Hour)
							}
							return nil
						},
					}
					next(wrappedRW, r)
					return
				}
				next(w, r)
		}
	}
}
