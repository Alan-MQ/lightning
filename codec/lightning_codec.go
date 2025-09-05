package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type LightningCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

func NewLightningCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (c *LightningCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *LightningCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *LightningCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc: gob error encoding header:", err)
		return
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc: gob error encoding body:", err)
		return
	}

	return
}

func (c *LightningCodec) Close() error {
	return c.conn.Close()
}
