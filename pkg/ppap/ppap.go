package ppap

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

type PPAPLayer struct {
	Type     uint8
	Hoplimit uint8
	Length   uint16
	Padding  uint32
	Origin   net.IP
	Pyload   []byte
}

const (
	PPAPLayerSize = 1 + 1 + 2 + 4 + 4
	PPAPVersion   = 254
	IHaveStr      = "I have a "
	PPAPReq       = "PPAP"
	PPAPAck       = "Ah!"
	PPAPClose     = "Pico"
)

func NewRequest(hoplimit uint8, origin net.IP) PPAPLayer {
	return New(hoplimit, origin, []byte(PPAPReq))
}

func NewAck(hoplimit uint8, origin net.IP) PPAPLayer {
	return New(hoplimit, origin, []byte(PPAPAck))
}

func NewClose(hoplimit uint8, origin net.IP) PPAPLayer {
	return New(hoplimit, origin, []byte(PPAPClose))
}

func NewIhave(hoplimit uint8, origin net.IP, obj string) PPAPLayer {
	return New(hoplimit, origin, []byte(IHaveStr+obj))
}

func New(hoplimit uint8, origin net.IP, pyload []byte) PPAPLayer {
	var ppap PPAPLayer
	ppap.Hoplimit = hoplimit
	ppap.Origin = origin
	ppap.Pyload = pyload
	return ppap
}

func (p *PPAPLayer) addTermNull() {
	p.Pyload = append(p.Pyload, byte(0))
}

func DialParam(v int) string {
	proto := strconv.Itoa(PPAPVersion)
	if v == 6 {
		return "ip6:" + proto
	}
	return "ip4:" + proto
}

func (p *PPAPLayer) Marshal() []byte {
	buf := make([]byte, PPAPLayerSize)
	buf[0] = p.Type
	buf[1] = p.Hoplimit
	binary.BigEndian.PutUint16(buf[2:], p.Length)
	binary.BigEndian.PutUint32(buf[4:], p.Padding)
	binary.BigEndian.PutUint32(buf[8:], ip2int(p.Origin))
	buf = append(buf, p.Pyload...)

	return buf
}

func Unmarshal(buf []byte) PPAPLayer {
	var p PPAPLayer
	p.Type = buf[0]
	p.Hoplimit = buf[1]
	p.Length = binary.BigEndian.Uint16(buf[2:])
	p.Padding = binary.BigEndian.Uint32(buf[4:])
	p.Origin = int2ip(binary.BigEndian.Uint32(buf[8:]))
	p.Pyload = buf[12:]

	return p
}

func GetStripIPv4(buf []byte) []byte {
	var size uint8
	size = (buf[0] & 0b00001111) * 4
	return buf[size:]
}

func GetHavesizeTrim(str string) string {
	str = strings.TrimLeft(str, IHaveStr)
	str = strings.TrimRight(str, "!")
	return str
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func (p *PPAPLayer) IsClose() bool {
	return p.ToPyloadStr() == "Pico"
}

func (p *PPAPLayer) ToPyloadStr() string {
	return string(p.Pyload)
}
