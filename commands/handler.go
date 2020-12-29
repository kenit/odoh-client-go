package commands

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"net/http"
)

type DnsHandler func(dns.ResponseWriter, *dns.Msg)

func GetHandler(target string, proxy string) (DnsHandler, error) {
	odohConfigs, err := fetchTargetConfigs(target)
	if err != nil {
		return nil, err
	}

	if len(odohConfigs.Configs) == 0 {
		err := errors.New("target provided no odoh configs")
		return nil, err
	}

	odohConfig := odohConfigs.Configs[0]

	return func(w dns.ResponseWriter, r *dns.Msg) {
		switch r.Opcode {
		case dns.OpcodeQuery:

			packedDnsQuery, err := r.Pack()
			if err != nil {
				fmt.Println(err)
				return
			}

			odohQuery, queryContext, err := createOdohQuestion(packedDnsQuery, odohConfig.Contents)
			if err != nil {
				fmt.Println(err)
				return
			}

			client := http.Client{}
			odohMessage, err := resolveObliviousQuery(odohQuery, true, target, proxy, &client)
			if err != nil {
				fmt.Println(err)
				return
			}

			dnsResponse, err := validateEncryptedResponse(odohMessage, queryContext)
			if err != nil {
				fmt.Println(err)
				return
			}

			if err := w.WriteMsg(dnsResponse); err != nil{
				fmt.Println(err)
			}
		}
	}, nil
}