package main

import (
	"strconv"
	"time"

	"google.golang.org/grpc/resolver"
)

// NewBuilderWithScheme creates a new test resolver builder with the given scheme.
func NewBuilderWithScheme(scheme string) *Resolver {
	return &Resolver{
		scheme: scheme,
	}
}

// Resolver is also a resolver builder.
// It's build() function always returns itself.
type Resolver struct {
	scheme string

	// Fields actually belong to the resolver.
	cc             resolver.ClientConn
	bootstrapAddrs []resolver.Address
}

// InitialAddrs adds resolved addresses to the resolver so that
// NewAddress doesn't need to be explicitly called after Dial.
func (r *Resolver) InitialAddrs(addrs []resolver.Address) {
	for _, addr := range addrs {
		num := addr.Metadata.(int)
		for ; num > 0; num-- {
			addr.Metadata = num
			r.bootstrapAddrs = append(r.bootstrapAddrs, addr)
		}
	}
}

// Build returns itself for Resolver, because it's both a builder and a resolver.
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r.cc = cc
	r.cc.NewAddress(r.bootstrapAddrs)
	return r, nil
}

// Scheme returns the test scheme.
func (r *Resolver) Scheme() string {
	return r.scheme
}

// ResolveNow is a noop for Resolver.
func (*Resolver) ResolveNow(o resolver.ResolveNowOption) {}

// Close is a noop for Resolver.
func (*Resolver) Close() {}

// GenerateAndRegisterManualResolver generates a random scheme and a Resolver
// with it. It also registers this Resolver.
// It returns the Resolver and a cleanup function to unregister it.
func GenerateAndRegisterManualResolver() (*Resolver, func()) {
	scheme := strconv.FormatInt(time.Now().UnixNano(), 36)
	r := NewBuilderWithScheme(scheme)
	resolver.Register(r)
	return r, func() { resolver.UnregisterForTesting(scheme) }
}
