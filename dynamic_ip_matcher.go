package dynamic_ip_matcher

import (
	"encoding/json"
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(MatchDynamicClientIP{})
}

// MatchDynamicClientIP matchers the requests bu the client IP address.
// The IP ranged are profided by modules to allow for ranges to be dynamic
type MatchDynamicClientIP struct {
	// A module which provides a source of IP ranges, from which
	// requests are matched.
	ProvidersRaw json.RawMessage         `json:"providers,omitempty" caddy:"namespace=http.ip_sources inline_key=source"`
	Providers    caddyhttp.IPRangeSource `json:"-"`

	logger *zap.Logger
}

func (MatchDynamicClientIP) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.matchers.dynamic_client_ip",
		New: func() caddy.Module {
			return new(MatchDynamicClientIP)
		},
	}
}

func (m *MatchDynamicClientIP) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume the directive name

	if !d.NextArg() {
		return d.ArgErr()
	}

	if m.Providers != nil {
		return d.Err("providers already specified")
	}

	dynModule := d.Val()
	modID := "http.ip_sources." + dynModule
	mod, err := caddyfile.UnmarshalModule(d, modID)

	if err != nil {
		return err
	}

	provider, ok := mod.(caddyhttp.IPRangeSource)

	if !ok {
		return d.Errf("module %s (%T) is not an IPRangeSource", modID, mod)
	}

	m.ProvidersRaw = caddyconfig.JSONModuleObject(provider, "source", dynModule, nil)

	return nil
}

func (m *MatchDynamicClientIP) Provision(ctx caddy.Context) error {
	m.logger = ctx.Logger()

	if m.ProvidersRaw != nil {
		val, err := ctx.LoadModule(m, "ProvidersRaw")

		if err != nil {
			return err
		}

		m.Providers = val.(caddyhttp.IPRangeSource)
	}

	return nil
}

func (m MatchDynamicClientIP) Match(r *http.Request) bool {
	address := caddyhttp.GetVar(r.Context(), caddyhttp.ClientIPVarKey).(string)
	clientIP, err := ParseIPZoneFromString(address)

	if err != nil {
		m.logger.Error("getting client IP", zap.Error(err))
		return false
	}

	matches := m.matchIP(r, clientIP)

	return matches
}

func ParseIPZoneFromString(address string) (netip.Addr, error) {
	ipStr, _, err := net.SplitHostPort(address)
	if err != nil {
		ipStr = address // OK; probably didn't have a port
	}

	// Client IP may contain a zone if IPv6, so we need
	// to pull that out before parsing the IP
	ipStr, _, _ = strings.Cut(ipStr, "%")

	ipAddr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return netip.IPv4Unspecified(), err
	}

	return ipAddr, nil
}

func (m *MatchDynamicClientIP) matchIP(r *http.Request, clientIP netip.Addr) bool {
	if m.Providers == nil {
		// We have no provier, So we can't match anything
		return false
	}

	cidrs := m.Providers.GetIPRanges(r)

	for _, ipRange := range cidrs {
		if ipRange.Contains(clientIP) {
			return true
		}
	}
	return false
}

// Interface guards
var (
	_ caddy.Module             = (*MatchDynamicClientIP)(nil)
	_ caddy.Provisioner        = (*MatchDynamicClientIP)(nil)
	_ caddyfile.Unmarshaler    = (*MatchDynamicClientIP)(nil)
	_ caddyhttp.RequestMatcher = (*MatchDynamicClientIP)(nil)
)
