package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync/atomic"

	"github.com/miekg/dns"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	oldIPs, _ := net.LookupIP("old-service")
	if len(oldIPs) == 0 {
		oldIPs = []net.IP{net.IPv4(192, 168, 10, 1)}
	}
	newIPs, _ := net.LookupIP("new-service")
	if len(newIPs) == 0 {
		newIPs = []net.IP{net.IPv4(192, 168, 10, 2)}
	}

	oldIP := oldIPs[0]
	newIP := newIPs[0]

	var count uint64
	dns.HandleFunc("nginx.test.", func(w dns.ResponseWriter, msg *dns.Msg) {
		var res dns.Msg
		res.SetReply(msg)
		if msg.Question[0].Qtype != dns.TypeA {
			w.WriteMsg(&res)
			return
		}

		var ip net.IP
		if atomic.AddUint64(&count, 1)%2 == 1 {
			ip = oldIP
		} else {
			ip = newIP
		}
		res.Answer = append(res.Answer, &dns.A{
			Hdr: dns.RR_Header{
				Name:   msg.Question[0].Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    10,
			},
			A: ip,
		})
		w.WriteMsg(&res)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "53"
	}
	err := dns.ListenAndServe(":"+port, "udp", nil)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}
