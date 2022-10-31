package main

import (
	"fmt"
	"reflect"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"
)

/*
	syntax

	chain_element

# MEMIF


	chain {
		point2pointipam
		mechanisms {
			'MEMIF' chain {
				sendfd
				up
				connectioncontext
				tag
				memif
			}
		}
	}

	chain {
		point2pointipam
		recvfd
		mechanisms {
			'KERNEL' kernel
		}
		dnscontext
		sendfd
	}

*/

func Build(source string, conf *Config) networkservice.NetworkServiceServer {

	var chainElements []networkservice.NetworkServiceServer
	var nodes = Parse(source)

	for _, n := range nodes {
		chainElements = append(chainElements, build(n, conf))
	}

	return next.NewNetworkServiceServer(chainElements...)
}

func build(n *Node, conf *Config) networkservice.NetworkServiceServer {
	var parameters = make(map[string]*reflect.Value)

	for i := 0; i < reflect.ValueOf(*conf).NumField(); i++ {
		var param = reflect.ValueOf(*conf).Field(i)
		parameters[param.Type().String()] = &param
	}

	var fun = conf.ChainElements[n.Name]

	if fun != nil {
		if p := extraParameters(n.Body, conf); p != nil {
			v := reflect.ValueOf(p)
			parameters[v.Type().String()] = &v
		}

		var args []reflect.Value
		for i := 0; i < reflect.ValueOf(fun).Type().NumIn(); i++ {
			var want = reflect.ValueOf(fun).Type().In(i)
			if reflect.ValueOf(fun).Type().NumIn() == 1 && reflect.ValueOf(fun).Type().IsVariadic() && parameters[want.String()] == nil {
				continue
			}
			args = append(args, *parameters[want.String()])
		}

		if reflect.ValueOf(fun).Type().IsVariadic() && len(args) > 0 {
			var last = args[len(args)-1]
			args[len(args)-1] = last.Index(0)
			for i := 1; i < last.Len(); i++ {
				args = append(args, last.Index(i))
			}
		}

		return reflect.ValueOf(fun).Call(args)[0].Interface().(networkservice.NetworkServiceServer)
	}

	return next.NewNetworkServiceServer()
}

func extraParameters(body []*Node, conf *Config) interface{} {
	if len(body) == 0 {
		return nil
	}

	if len(body[0].Metadata) > 0 {
		var r = make(map[string]networkservice.NetworkServiceServer)

		for _, n := range body {
			fmt.Println(n.Name)
			var s = build(n, conf)
			for _, k := range n.Metadata {
				r[k] = s
			}
		}

		return r
	}

	var result []networkservice.NetworkServiceServer
	for _, n := range body {

		result = append(result, build(n, conf))
	}

	return result
}
