package qr

import (
	"bytes"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"io"
)

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

func QRCodeToBytes(qrcode *qrcode.QRCode, opts ...standard.ImageOption) ([]byte, error) {
	buf := &bytes.Buffer{}
	w := standard.NewWithWriter(nopWriteCloser{buf}, opts...)
	if err := qrcode.Save(w); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
