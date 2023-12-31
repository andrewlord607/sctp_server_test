package main

import (
	"flag"
	"github.com/ishidawataru/sctp"
	"log"
	"net"
	"strings"
)

func serveClient(conn net.Conn, bufsize int) error {
	for {
		//buf := make([]byte, bufsize+128) // add overhead of SCTPSndRcvInfoWrappedConn
		buf := make([]byte, 1024) // add overhead of SCTPSndRcvInfoWrappedConn
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read failed: %v", err)
			return err
		}
		//log.Printf("message %v \n", time.Now().Format(time.StampMicro))
		//log.Printf("read: %d", n)
		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("write failed: %v", err)
			return err
		}
		//log.Printf("write: %d", n)
	}
}

func main() {
	//var sndbuf = flag.Int("sndbuf", 0, "")
	//var rcvbuf = flag.Int("rcvbuf", 0, "")
	var ip = flag.String("ip", "0.0.0.0", "")
	var bufsize = flag.Int("bufsize", 208, "")
	var port = flag.Int("port", 3868, "")

	flag.Parse()

	var ips []net.IPAddr

	for _, i := range strings.Split(*ip, ",") {
		if a, err := net.ResolveIPAddr("ip", i); err == nil {
			log.Printf("Resolved address '%s' to %s", i, a)
			ips = append(ips, *a)
		} else {
			log.Printf("Error resolving address '%s': %v", i, err)
		}
	}

	addr := &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    *port,
	}
	log.Printf("raw addr: %+v\n", addr.ToRawSockAddrBuf())

	ln, err := sctp.ListenSCTPExt(
		"sctp",
		addr,
		sctp.InitMsg{
			NumOstreams:  16,
			MaxInstreams: 16})
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listen on %s", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("failed to accept: %v", err)
		}
		log.Printf("Accepted Connection from RemoteAddr: %s", conn.RemoteAddr())
		/*wconn := sctp.NewSCTPSndRcvInfoWrappedConn(conn.(*sctp.SCTPConn))
		if *sndbuf != 0 {
			err = wconn.SetWriteBuffer(*sndbuf)
			if err != nil {
				log.Fatalf("failed to set write buf: %v", err)
			}
		}
		if *rcvbuf != 0 {
			err = wconn.SetReadBuffer(*rcvbuf)
			if err != nil {
				log.Fatalf("failed to set read buf: %v", err)
			}
		}
		*sndbuf, err = wconn.GetWriteBuffer()
		if err != nil {
			log.Fatalf("failed to get write buf: %v", err)
		}
		*rcvbuf, err = wconn.GetWriteBuffer()
		if err != nil {
			log.Fatalf("failed to get read buf: %v", err)
		}
		log.Printf("SndBufSize: %d, RcvBufSize: %d", *sndbuf, *rcvbuf)*/

		go serveClient(conn, *bufsize)
	}
}
