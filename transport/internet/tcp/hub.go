package tcp

import (
	"context"

	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/transport/internet"
)

// ListenUnix is the UDS version of ListenTCP
func ListenUnix(ctx context.Context, address net.Address, settings *internet.MemoryStreamConfig, handler internet.ConnHandler) (internet.Listener, error) {
	if settings == nil {
		s, err := internet.ToMemoryStreamConfig(nil)
		if err != nil {
			return nil, errors.New("failed to create default unix stream settings").Base(err)
		}
		settings = s
	}

	protocol := settings.ProtocolName
	listenFunc, err := internet.GetTransportListener(protocol)
	if err != nil {
		return nil, errors.New(protocol, " unix listener not registered.").AtError()
	}
	listener, err := listenFunc(ctx, address, net.Port(0), settings, handler)
	if err != nil {
		return nil, errors.New("failed to listen on unix address: ", address).Base(err)
	}
	return listener, nil
}

func ListenTCP(ctx context.Context, address net.Address, port net.Port, settings *internet.MemoryStreamConfig, handler internet.ConnHandler) (internet.Listener, error) {
	if settings == nil {
		s, err := internet.ToMemoryStreamConfig(nil)
		if err != nil {
			return nil, errors.New("failed to create default stream settings").Base(err)
		}
		settings = s
	}

	if address.Family().IsDomain() && address.Domain() == "localhost" {
		address = net.LocalHostIP
	}

	if address.Family().IsDomain() {
		return nil, errors.New("domain address is not allowed for listening: ", address.Domain())
	}

	protocol := settings.ProtocolName
	listenFunc, err := internet.GetTransportListener(protocol)
	if err != nil {
		return nil, errors.New(protocol, " listener not registered.").AtError()
	}
	listener, err := listenFunc(ctx, address, port, settings, handler)
	if err != nil {
		return nil, errors.New("failed to listen on address: ", address, ":", port).Base(err)
	}
	return listener, nil
}
