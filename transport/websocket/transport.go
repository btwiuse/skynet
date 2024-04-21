package websocket

import (
	"context"
	"net"

	"github.com/webteleport/webteleport/transport"
)

type Transport struct{}

var _ transport.Transport = (*Transport)(nil)

func (t *Transport) Dial(ctx context.Context, addr string) (transport.Session, error) {
	return Dial(ctx, addr, nil)
}

func (t *Transport) Listen(ctx context.Context, addr string) (net.Listener, error) {
	return Listen(ctx, addr)
}
