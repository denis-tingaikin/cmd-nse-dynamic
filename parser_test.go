package main_test

import (
	"context"
	"net"
	"testing"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
	main "github.com/networkservicemesh/cmd-nse-dynamic"
	"github.com/stretchr/testify/require"
)

func Test_Build(t *testing.T) {
	var conf main.Config
	var _, ipnet, _ = net.ParseCIDR("169.254.0.0/16")
	conf.CidrPrefix = []*net.IPNet{ipnet}

	conf.ChainElements = main.ChainElements
	var source = `chain {
		point2pointipam
		mechanisms {
			'KERNEL': kernel
		}
	}`

	var server = main.Build(source, &conf)

	conn, err := server.Request(context.Background(), &networkservice.NetworkServiceRequest{
		Connection: &networkservice.Connection{
			Context: &networkservice.ConnectionContext{
				IpContext: &networkservice.IPContext{},
			},
		},
		MechanismPreferences: []*networkservice.Mechanism{
			{
				Type: "KERNEL",
				Cls:  "LOCAL",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, conn)
	require.NotNil(t, conn.Mechanism)
	require.NotEmpty(t, conn.GetContext().GetIpContext().GetSrcIpAddrs())
	require.NotEmpty(t, conn.GetContext().GetIpContext().GetDstIpAddrs())
}

func TestParse(t *testing.T) {
	var nodes = main.Parse(`
	chain {
		a
		b
		c
		e {
			f
			'Option1': "Option2": b
			c		
		}
	}
	`)

	require.Len(t, nodes, 1)
	var elements = nodes[0]

	require.Equal(t, "chain", elements.Name)
	require.Len(t, elements.Body, 4)

	require.Contains(t, elements.Body[0].Name, "a")
	require.Contains(t, elements.Body[1].Name, "b")
	require.Contains(t, elements.Body[2].Name, "c")
	require.Contains(t, elements.Body[3].Name, "e")

	elements = elements.Body[3]

	require.Equal(t, "f", elements.Body[0].Name)
	require.Equal(t, "b", elements.Body[1].Name)
	require.Equal(t, "c", elements.Body[2].Name)

	require.Equal(t, []string{"Option1", "Option2"}, elements.Body[1].Metadata)

}
