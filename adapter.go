package main

import (
	"github.com/kenit/odoh-client-go/commands"
	"github.com/miekg/dns"
	"net"
)

var adapters = []Adapter{}

type Adapter func(handler commands.DnsHandler) commands.DnsHandler

func Adapt(h commands.DnsHandler, adapters ...Adapter) commands.DnsHandler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func AddAdapter(a Adapter){
	adapters = append(adapters, a)
}

type WrappedResponseWriter struct{
	originResponseWriter dns.ResponseWriter
	beforeWriteMsg func(*dns.Msg) error
}

func(w *WrappedResponseWriter) LocalAddr() net.Addr {
	return w.originResponseWriter.LocalAddr()
}

func(w *WrappedResponseWriter) RemoteAddr() net.Addr{
	return w.originResponseWriter.RemoteAddr()
}

func (w *WrappedResponseWriter) WriteMsg(m *dns.Msg) error {
	if err := w.beforeWriteMsg(m); err != nil{
		return err
	}
	return w.originResponseWriter.WriteMsg(m)
}

func (w *WrappedResponseWriter) Write(m []byte) (int, error){
	return w.originResponseWriter.Write(m)
}

func (w *WrappedResponseWriter) Close() error{
	return w.originResponseWriter.Close()
}

func (w *WrappedResponseWriter) TsigStatus() error{
	return w.originResponseWriter.TsigStatus()
}

func (w *WrappedResponseWriter) TsigTimersOnly(b bool){
	w.originResponseWriter.TsigTimersOnly(b)
}

func (w *WrappedResponseWriter) Hijack(){
	w.originResponseWriter.Hijack()
}
