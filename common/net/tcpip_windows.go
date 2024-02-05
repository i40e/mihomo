package net

import (
	"net"
	"syscall"
	"time"
	"unsafe"
)

func TCPKeepAlive(c net.Conn) {
	if tcp, ok := c.(*net.TCPConn); ok {
		sysConn, err := tcp.SyscallConn()
		if err != nil {
			return
		}
		idleMsecs := uint32(roundDurationUp(KeepAliveIdle, time.Millisecond))
		intervalMsecs := uint32(roundDurationUp(KeepAliveInterval, time.Millisecond))
		sysConn.Control(func(fd uintptr) {
			ka := syscall.TCPKeepalive{
				OnOff:    1,
				Time:     idleMsecs,
				Interval: intervalMsecs,
			}
			ret := uint32(0)
			size := uint32(unsafe.Sizeof(ka))
			_ = syscall.WSAIoctl(syscall.Handle(fd), syscall.SIO_KEEPALIVE_VALS, (*byte)(unsafe.Pointer(&ka)), size, nil, 0, &ret, nil, 0)
		})
	}
}
