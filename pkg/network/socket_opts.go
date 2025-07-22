package network

import (
	"net"
	"syscall"
	"time"
)

func SetSocketOptions(conn net.Conn) {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return
	}

	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(30 * time.Second)
	tcpConn.SetNoDelay(true) // disable Nagle's algorithm for low-latency

	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return
	}

	_ = rawConn.Control(func(fd uintptr) {
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, 1<<20) // 1MB recv buffer
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF, 1<<20) // 1MB send buffer
	})
}
