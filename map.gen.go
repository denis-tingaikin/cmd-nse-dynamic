package main

import (
	"fmt"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/sdk/pkg/networkservice/common/mechanisms"
	"github.com/networkservicemesh/sdk/pkg/networkservice/common/mechanisms/kernel"
	"github.com/networkservicemesh/sdk/pkg/networkservice/common/mechanisms/recvfd"
	"github.com/networkservicemesh/sdk/pkg/networkservice/common/mechanisms/sendfd"
	"github.com/networkservicemesh/sdk/pkg/networkservice/connectioncontext/dnscontext"
	"github.com/networkservicemesh/sdk/pkg/networkservice/core/chain"
	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"
	"github.com/networkservicemesh/sdk/pkg/networkservice/ipam/point2pointipam"
)

var ChainElements = map[string]any{
	"next":            next.NewNetworkServiceServer,
	"point2pointipam": point2pointipam.NewServer,
	"chain":           chain.NewNetworkServiceServer,
	"sendfd":          sendfd.NewServer,
	"recvfd":          recvfd.NewServer,
	"mechanisms": func(m map[string]networkservice.NetworkServiceServer) networkservice.NetworkServiceServer {
		fmt.Println("used")
		fmt.Println(m["KERNEL"])
		return mechanisms.NewServer(m)
	},
	"dnscontext": dnscontext.NewServer,
	"kernel": func() networkservice.NetworkServiceServer {
		fmt.Println("used2")
		return kernel.NewServer()
	},
}
