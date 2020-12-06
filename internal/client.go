package internal

import (
	"log"
	"net"

	"github.com/pkg/errors"
	"github.com/takehaya/PPAP_Protocol/pkg/ppap"
)

const ClientHopLimit = 1

func Client(ctl *Controller) error {
	conn, err := net.DialIP(ppap.DialParam(4), ctl.SrcAddr, ctl.GatewayAddr)
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	var (
		ppapPyload ppap.PPAPLayer
		havelist   []string
	)

	ppapPyload = ppap.NewRequest(ClientHopLimit, ctl.SrcAddr.IP)
	log.Println("send PPAP")

	err = ctl.PPAPSender(conn, ppapPyload)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	for {
		recvbuffer := make([]byte, 9000)
		n, err := conn.Read(recvbuffer)
		if err != nil {
			return errors.WithMessage(err, "Read error:")
		}
		recvbuffer = ppap.GetStripIPv4(recvbuffer[:n])
		ppapPyload = ppap.Unmarshal(recvbuffer)

		log.Println("recv: ", string(ppapPyload.Pyload))

		havelist = append(havelist, string(ppapPyload.Pyload))
		if 2 <= len(havelist) {
			break
		}
	}

	ppapPyload = ppap.NewAck(ClientHopLimit, ctl.SrcAddr.IP)
	err = ctl.PPAPSender(conn, ppapPyload)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	h1 := ppap.GetHavesizeTrim(havelist[0])
	h2 := ppap.GetHavesizeTrim(havelist[1])

	key := []byte(h1 + "-" + h2 + "!")
	ppapPyload = ppap.New(ClientHopLimit, ctl.SrcAddr.IP, key)

	log.Println("send: ", ppapPyload.ToPyloadStr())
	err = ctl.PPAPSender(conn, ppapPyload)
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
		log.Printf("Pyload %s\n", ppapPyload.ToPyloadStr())
		if ppapPyload.IsClose() {
			break
		}
	}
	log.Println("send ppap packet success!")

	return nil
}
