//go:build !android
// +build !android

package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	kcp "github.com/xtaci/kcp-go/v5"
	"github.com/xtaci/kcptun/generic"
	"github.com/xtaci/tcpraw"
)

func dial(config *Config, block kcp.BlockCrypt) (*kcp.UDPSession, error) {
	// Parse the multiPort definition only once.
	multiPortOnce.Do(func() {
		multiPort, multiPortParseError = generic.ParseMultiPort(config.RemoteAddr)
	})

	// Abort when the multiPort definition is invalid.
	if multiPortParseError != nil {
		return nil, multiPortParseError
	}

	// Pick a random destination port within the configured range.
	var randport uint64
	err := binary.Read(rand.Reader, binary.LittleEndian, &randport)
	if err != nil {
		return nil, err
	}

	remoteAddr := fmt.Sprintf("%v:%v", multiPort.Host, uint64(multiPort.MinPort)+randport%uint64(multiPort.MaxPort-multiPort.MinPort+1))

	if config.TCP {
		conn, err := tcpraw.Dial("tcp", remoteAddr)
		if err != nil {
			return nil, errors.Wrap(err, "tcpraw.Dial()")
		}
		return kcp.NewConn(remoteAddr, block, config.DataShard, config.ParityShard, conn)
	}
	return kcp.DialWithOptions(remoteAddr, block, config.DataShard, config.ParityShard)

}
