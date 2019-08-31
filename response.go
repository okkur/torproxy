package torproxy

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"sync"

	"golang.org/x/net/proxy"
)

type TorResponse struct {
	headers    http.Header
	body       []byte
	bodyReader bytes.Buffer
	bodyWriter bytes.Buffer
	status     int

	request *http.Request
	dialer  proxy.Dialer
}

var bufferPool = sync.Pool{New: createBuffer}

func createBuffer() interface{} {
	return make([]byte, 0, 32*1024)
}

// Header returns response headers
func (t *TorResponse) Header() http.Header {
	return t.headers
}

func (t *TorResponse) Write(body []byte) (int, error) {
	reader := bytes.NewReader(body)
	pooledIoCopy(&t.bodyReader, reader)
	t.body = body
	return len(body), nil
}

// Body returns response's body. This method should only get called after WriteBody()
func (t *TorResponse) Body() []byte {
	return t.bodyWriter.Bytes()
}

// WriteHeader Writes the given status code to response
func (t *TorResponse) WriteHeader(status int) {
	t.status = status
}

func (t *TorResponse) ReplaceBody(scheme, to, host string) error {
	replacedBody := bytes.Replace(t.bodyWriter.Bytes(), []byte(scheme+"://"+to), []byte(scheme+"://"+host), -1)
	t.bodyWriter.Reset()
	if _, err := t.bodyWriter.Write(replacedBody); err != nil {
		return err
	}
	return nil
}

func (t *TorResponse) WriteBody() error {
	switch t.Header().Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(&t.bodyReader)
		if err != nil {
			return err
		}
		defer reader.Close()
		_, err = io.Copy(&t.bodyWriter, reader)
		if err != nil {
			return err
		}
		t.Header().Del("Content-Encoding")
	default:
		_, err := io.Copy(&t.bodyWriter, &t.bodyReader)
		if err != nil {
			return err
		}
	}
	return nil
}

var skipHeaders = map[string]struct{}{
	"Content-Type":        {},
	"Content-Disposition": {},
	"Accept-Ranges":       {},
	"Set-Cookie":          {},
	"Cache-Control":       {},
	"Expires":             {},
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		if _, ok := dst[k]; ok {
			if _, shouldSkip := skipHeaders[k]; shouldSkip {
				continue
			}
			if k != "Server" {
				dst.Del(k)
			}
		}
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func pooledIoCopy(dst io.Writer, src io.Reader) {
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	bufCap := cap(buf)
	io.CopyBuffer(dst, src, buf[0:bufCap:bufCap])
}
