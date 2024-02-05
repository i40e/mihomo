//go:build aix || freebsd || linux || netbsd

package net

import (
	"net"
	"syscall"
	"time"

	"github.com/metacubex/mihomo/log"
)

func TCPKeepAlive(c net.Conn) {
	if tcp, ok := c.(*net.TCPConn); ok {
		tcp.SetKeepAlive(true)
		sysConn, err := tcp.SyscallConn()
		if err != nil {
			log.Errorln("[TCPKeepAlive] SyscallConn error: %s", err.Error())
			_ = tcp.SetKeepAlivePeriod(KeepAliveInterval)
			return
		}
		idleSecs := int(roundDurationUp(KeepAliveIdle, time.Second))
		intervalSecs := int(roundDurationUp(KeepAliveInterval, time.Second))
		sysConn.Control(func(fdptr uintptr) {
			fd := int(fdptr)
			err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPIDLE, idleSecs)
			if err != nil {
				log.Errorln("[TCPKeepAlive] SetsockoptInt TCP_KEEPIDLE error: %s", err.Error())
			}
			err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL, intervalSecs)
			if err != nil {
				log.Errorln("[TCPKeepAlive] SetsockoptInt TCP_KEEPINTVL error: %s", err.Error())
			}
		})
	}
}
