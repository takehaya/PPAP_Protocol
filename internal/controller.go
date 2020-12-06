package internal

import (
	"bytes"
	"net"

	"github.com/google/gopacket/routing"
	"github.com/pkg/errors"
	"github.com/takehaya/PPAP_Protocol/pkg/ppap"
)

type Controller struct {
	Router      routing.Router
	GatewayAddr *net.IPAddr
	SrcAddr     *net.IPAddr
	sendbuffer  bytes.Buffer
	Have1       string
	Have2       string
}

func NewController(dstAddr, bindaddr, h1, h2 string) (*Controller, error) {
	dst := net.IPAddr{IP: net.ParseIP(dstAddr)}
	r, err := routing.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ctl := &Controller{
		Router:      r,
		GatewayAddr: &dst,
		Have1:       h1,
		Have2:       h2,
	}
	if bindaddr == "auto" {
		src, err := ctl.getfibBaseSrcAddr(dst.IP)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ctl.SrcAddr = &net.IPAddr{IP: src}
	} else {
		ctl.SrcAddr = &net.IPAddr{IP: net.ParseIP(bindaddr)}
	}

	return ctl, nil
}

func (c *Controller) getfibBaseSrcAddr(dstAddr net.IP) (net.IP, error) {
	_, _, preferredSrc, err := c.Router.Route(dstAddr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return preferredSrc, nil
}

func (c *Controller) PPAPSender(conn *net.IPConn, p ppap.PPAPLayer) error {
	if _, err := conn.Write(p.Marshal()); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Controller) PPAPSenderToIP(conn *net.IPConn, p ppap.PPAPLayer, dstAddr *net.IPAddr) error {
	// Send PPAP Request
	if _, err := conn.WriteTo(p.Marshal(), dstAddr); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
