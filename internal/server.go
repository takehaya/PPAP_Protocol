package internal

import (
	"log"
	"net"

	"github.com/pkg/errors"
	"github.com/takehaya/PPAP_Protocol/pkg/ppap"
)

func Server(ctl *Controller) error {
	conn, err := net.ListenIP(ppap.DialParam(4), ctl.SrcAddr)
	if err != nil {
		log.Panic(err)
	}
	buf := make([]byte, 9000)

	for {
		var ppapPyload ppap.PPAPLayer
		n, ra, err := conn.ReadFromIP(buf)
		if err != nil {
			log.Panic(err)
		}
		ppapPyload = ppap.Unmarshal(buf[:n])
		log.Printf("recv: %s\n", ppapPyload.ToPyloadStr())

		n1 := ppap.NewIhave(ppapPyload.Hoplimit, ppapPyload.Origin, ctl.Have1)
		err = ctl.PPAPSenderToIP(conn, n1, ra)
		if err != nil {
			log.Fatal("encode error:", err)
		}
		log.Printf("send: %s\n", n1.ToPyloadStr())

		n2 := ppap.NewIhave(ppapPyload.Hoplimit, ppapPyload.Origin, ctl.Have2)
		err = ctl.PPAPSenderToIP(conn, n2, ra)
		if err != nil {
			log.Fatal("encode error:", err)
		}
		log.Printf("send: %s\n", n2.ToPyloadStr())

		n, ra, err = conn.ReadFromIP(buf)
		if err != nil {
			log.Panic(err)
		}
		ppapPyload = ppap.Unmarshal(buf[:n])
		log.Printf("recv %s\n", ppapPyload.ToPyloadStr())

		n, ra, err = conn.ReadFromIP(buf)
		if err != nil {
			log.Panic(err)
		}
		ppapPyload = ppap.Unmarshal(buf[:n])
		log.Printf("recv %s\n", ppapPyload.ToPyloadStr())

		// todo goroutine
		if 0 < ppapPyload.Hoplimit {
			err = redirect(conn, ctl, &ppapPyload)
			if err != nil {
				log.Panic(err)
			}
		} else {
			cl := ppap.NewClose(ppapPyload.Hoplimit, ppapPyload.Origin)
			log.Printf("recv %s\n", cl.ToPyloadStr())
			err = ctl.PPAPSenderToIP(conn, cl, ra)
			if err != nil {
				log.Fatal("encode error:", err)
			}
			break
		}
	}
	return nil
}

func redirect(orgConn *net.IPConn, ctl *Controller, orgPpapPyload *ppap.PPAPLayer) error {
	conn, err := net.DialIP(ppap.DialParam(4), ctl.SrcAddr, ctl.GatewayAddr)
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	ppapPyload := ppap.NewRequest(orgPpapPyload.Hoplimit-1, orgPpapPyload.Origin)

	err = ctl.PPAPSender(conn, ppapPyload)
	log.Printf("send %s\n", ppapPyload.ToPyloadStr())

	if err != nil {
		log.Fatal("encode error:", err)
	}
	var havelist []string

	for {
		recvbuffer := make([]byte, 9000)
		n, err := conn.Read(recvbuffer)
		if err != nil {
			return errors.WithMessage(err, "Read error:")
		}
		recvbuffer = ppap.GetStripIPv4(recvbuffer[:n])
		ppapPyload = ppap.Unmarshal(recvbuffer)
		log.Printf("recv %s\n", ppapPyload.ToPyloadStr())

		havelist = append(havelist, string(ppapPyload.Pyload))
		if 2 <= len(havelist) {
			break
		}
	}
	ppapPyload = ppap.NewAck(orgPpapPyload.Hoplimit-1, orgPpapPyload.Origin)
	err = ctl.PPAPSender(conn, ppapPyload)
	log.Printf("send %s\n", ppapPyload.ToPyloadStr())

	if err != nil {
		log.Fatal("encode error:", err)
	}
	h1 := ppap.GetHavesizeTrim(havelist[0])
	h2 := ppap.GetHavesizeTrim(havelist[1])
	key := []byte(h1 + "-" + h2 + "-" + orgPpapPyload.ToPyloadStr())
	ppapPyload = ppap.New(orgPpapPyload.Hoplimit-1, orgPpapPyload.Origin, key)
	err = ctl.PPAPSender(conn, ppapPyload)
	log.Printf("send %s\n", ppapPyload.ToPyloadStr())

	if err != nil {
		log.Fatal("encode error:", err)
	}

	// wait
	for {
		recvbuffer := make([]byte, 9000)
		n, err := conn.Read(recvbuffer)
		if err != nil {
			return errors.WithMessage(err, "Read error:")
		}
		recvbuffer = ppap.GetStripIPv4(recvbuffer[:n])
		ppapPyload = ppap.Unmarshal(recvbuffer)
		log.Printf("recv %s\n", ppapPyload.ToPyloadStr())
		if ppapPyload.IsClose() {
			ppapPyload = ppap.NewClose(orgPpapPyload.Hoplimit-1, orgPpapPyload.Origin)
			log.Printf("send %s\n", ppapPyload.ToPyloadStr())
			err = ctl.PPAPSenderToIP(orgConn, ppapPyload, &net.IPAddr{IP: orgPpapPyload.Origin})
			if err != nil {
				log.Fatal("encode error:", err)
			}
			break
		}
	}
	log.Println("recv ppap packet success!")
	return nil
}
