package guard

import (
	"fmt"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

func init() {
	caddy.RegisterPlugin("guard", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}
func setup(c *caddy.Controller) error {
	g := New()

	for c.Next() {
		// shift cursor past alternate
		if !c.Next() {
			return c.ArgErr()
		}

		var (
			original bool
			err      error
		)

		if original, err = getOriginal(c); err != nil {
			return err
		}

		handler, err := initForward(c)
		if err != nil {
			return plugin.Error("alternate", err)
		}
		g.handlers = append(g.handlers, handler)

		if original {
			g.original = true
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		g.Next = next
		return g
	})

	c.OnShutdown(func() error {
		for _, handler := range g.handlers {
			if err := handler.OnShutdown(); err != nil {
				return err
			}
		}
		return nil
	})

	return nil
}

const original = "original"

func getOriginal(c *caddy.Controller) (bool, error) {
	if c.Val() == original {
		// shift cursor past original
		if !c.Next() {
			return false, c.ArgErr()
		}
		return true, nil
	}

	return false, nil
}

func getRCodes(c *caddy.Controller) ([]int, error) {
	in := strings.Split(c.Val(), ",")

	rcodes := make(map[int]interface{}, len(in))

	for _, rcode := range in {
		var rc int
		var ok bool

		if rc, ok = dns.StringToRcode[strings.ToUpper(rcode)]; !ok {
			return nil, fmt.Errorf("%s is not a valid rcode", rcode)
		}

		rcodes[rc] = nil
	}

	results := make([]int, 0, len(rcodes))
	for r := range rcodes {
		results = append(results, r)
	}

	return results, nil
}
