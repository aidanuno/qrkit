package qr

import (
	"bytes"
	"io"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

func CodeToBytes(code *qrcode.QRCode, opts ...standard.ImageOption) ([]byte, error) {
	buf := &bytes.Buffer{}
	w := standard.NewWithWriter(nopWriteCloser{buf}, opts...)
	if err := code.Save(w); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
