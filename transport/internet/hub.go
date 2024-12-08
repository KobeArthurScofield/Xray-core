package internet

import (
	"context"

	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/transport/internet/stat"
)

var transportListenerCache = make(map[string]ListenFunc)

func RegisterTransportListener(protocol string, listener ListenFunc) error {
	if _, found := transportListenerCache[protocol]; found {
		return errors.New(protocol, " listener already registered.").AtError()
	}
	transportListenerCache[protocol] = listener
	return nil
}

// For inbound listener to get the transport
func GetTransportListener(protocol string) (ListenFunc, error) {
	listenFunc := transportListenerCache[protocol]
	if listenFunc == nil {
		return nil, errors.New(protocol, " transport listener is not registered.").AtError()
	}
	return listenFunc, nil
}

type ConnHandler func(stat.Connection)

type ListenFunc func(ctx context.Context, address net.Address, port net.Port, settings *MemoryStreamConfig, handler ConnHandler) (Listener, error)

type Listener interface {
	Close() error
	Addr() net.Addr
}
// ListenSystem listens on a local address for incoming TCP connections.
//
// xray:api:beta
func ListenSystem(ctx context.Context, addr net.Addr, sockopt *SocketConfig) (net.Listener, error) {
	return effectiveListener.Listen(ctx, addr, sockopt)
}

// ListenSystemPacket listens on a local address for incoming UDP connections.
//
// xray:api:beta
func ListenSystemPacket(ctx context.Context, addr net.Addr, sockopt *SocketConfig) (net.PacketConn, error) {
	return effectiveListener.ListenPacket(ctx, addr, sockopt)
}
