package guard

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/miekg/dns"
)

type Guard struct {
	Next     plugin.Handler
	original bool // At least one rule has "original" flag
	handlers []HandlerWithCallbacks
}

// HandlerWithCallbacks interface is made for handling the requests
type HandlerWithCallbacks interface {
	plugin.Handler
	OnStartup() error
	OnShutdown() error
}

// New initializes Guard plugin
func New() (f *Guard) {
	return &Guard{}
}

// ServeDNS implements the plugin.Handler interface.
func (g Guard) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	//	var originalRequest *dns.Msg
	//if g.original {
	//originalRequest = r.Copy()
	//}
	nw := nonwriter.New(w)
	rcode, err := plugin.NextOrFailure(g.Name(), g.Next, ctx, nw, r)
	if nw.Msg != nil {
		w.WriteMsg(nw.Msg)
	}
	return rcode, err
}

// Name implements the Handler interface.
func (f Guard) Name() string { return "guard" }
