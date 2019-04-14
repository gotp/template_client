// Copyright 2018 doctorwechat
//
// Author: juzhongguoji <juzhongguoji@hotmail.com>
// Date:   2018/12/23

package name_resolver

import (
    "google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(NewLocalResolver())
}

// NewResolver creates a new resolver builder
func NewLocalResolver() *LocalResolver {
	return &LocalResolver{
		scheme: "local",
	}
}

// Resolver is also a resolver builder.
// It's build() function always returns itself.
type LocalResolver struct {
	scheme string
	// Fields actually belong to the resolver.
	clientConn	resolver.ClientConn
}

// Build returns itself for Resolver, because it's both a builder and a resolver.
func (this *LocalResolver) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	this.clientConn = clientConn
	addrs, found := GetResolverConfig().FindAddressByName(target.Endpoint)
	if (found) {
		this.clientConn.NewAddress(buildAddress(addrs))
	} else {
		return nil, nil
	}
	return this, nil
}

// Scheme returns the test scheme.
func (this *LocalResolver) Scheme() string {
	return this.scheme
}

// ResolveNow is a noop for Resolver.
func (*LocalResolver) ResolveNow(o resolver.ResolveNowOption) {}

// Close is a noop for Resolver.
func (*LocalResolver) Close() {}

func buildAddress(strAddrs []string) []resolver.Address {
	var addrs []resolver.Address
	for _, strAddr := range strAddrs {
		addrs = append(addrs, resolver.Address{Addr: strAddr})
	}
	return addrs
}
