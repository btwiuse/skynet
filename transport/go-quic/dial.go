package quic

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/webtransport/quic"
)

// 2^60 == 1152921504606846976
var MaxBidiRemoteStreams int64 = 1 << 60

var TLSConfig = &tls.Config{
	InsecureSkipVerify: true,
}

var QUICConfig = &quic.Config{
	TLSConfig:            TLSConfig,
	MaxBidiRemoteStreams: MaxBidiRemoteStreams,
}

func Dial(ctx context.Context, addr string) (*QuicSession, error) {
	l, err := quic.Listen("udp", ":0", QUICConfig)
	if err != nil {
		return nil, err
	}
	session, err := l.Dial(ctx, "udp", addr, QUICConfig)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (go-quic): %w", addr, err)
	}
	return &QuicSession{session}, nil
}
