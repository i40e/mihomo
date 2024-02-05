package net

import (
	"net"
	"syscall"
	"time"
)

// syscall.TCP_KEEPINTVL is missing on some darwin architectures.
const sysTCP_KEEPINTVL = 0x101

func TCPKeepAlive(c net.Conn) {
	if tcp, ok := c.(*net.TCPConn); ok {
		tcp.SetKeepAlive(true)
		sysConn, err := tcp.SyscallConn()
		if err != nil {
			return
		}
		idleSecs := int(roundDurationUp(KeepAliveIdle, time.Second))
		intervalSecs := int(roundDurationUp(KeepAliveInterval, time.Second))
		sysConn.Control(func(fd uintptr) {
			_ = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_KEEPALIVE, idleSecs)
			_ = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, sysTCP_KEEPINTVL, intervalSecs)
		})
	}
}
