package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Compress handler
func Compress() gin.HandlerFunc {
	return NewCompressor()
}

// Compressor middleware compressor
type Compressor struct {
	gin.ResponseWriter
	Writer io.Writer
}

// NewCompressor return new middleware compressor
func NewCompressor() gin.HandlerFunc {
	return newCompressHandler().Handle
}

// Write to writer
func (c *Compressor) Write(b []byte) (int, error) {
	return c.Writer.Write(b)
}

// WriteString to writer
func (c *Compressor) WriteString(s string) (int, error) {
	c.Header().Del("Content-Length")
	return c.Writer.Write([]byte(s))
}

// WriteHeader write to header
func (c *Compressor) WriteHeader(code int) {
	c.Header().Del("Content-Length")
	c.ResponseWriter.WriteHeader(code)
}

type compressHandler struct {
	pool sync.Pool
}

func newCompressHandler() *compressHandler {
	handler := &compressHandler{
		pool: sync.Pool{
			New: func() interface{} {
				gz, err := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
				if err != nil {
					panic(err)
				}
				return gz
			},
		},
	}
	return handler
}

// Handle handler of compress
func (c *compressHandler) Handle(ctx *gin.Context) {
	if ctx.Request.Header.Get("Content-Encoding") == "gzip" {
		c.DecompressHandle(ctx)
	}
	if !c.canCompress(ctx.Request) {
		return
	}
	gz := c.pool.Get().(*gzip.Writer)
	defer c.pool.Put(gz)
	defer gz.Reset(io.Discard)
	gz.Reset(ctx.Writer)

	ctx.Header("Content-Encoding", "gzip")
	ctx.Header("Vary", "Accept-Encoding")
	ctx.Writer = &Compressor{ctx.Writer, gz}
	defer func() {
		gz.Close()
		ctx.Header("Content-Length", fmt.Sprint(ctx.Writer.Size()))
	}()
	ctx.Next()
}

// DecompressHandle handler of decompress
func (c *compressHandler) DecompressHandle(ctx *gin.Context) {
	if ctx.Request.Body == nil {
		return
	}
	r, err := gzip.NewReader(ctx.Request.Body)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.Request.Header.Del("Content-Encoding")
	ctx.Request.Header.Del("Content-Length")
	ctx.Request.Body = r
}

func (c *compressHandler) canCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") ||
		strings.Contains(req.Header.Get("Connection"), "Upgrade") ||
		strings.Contains(req.Header.Get("Accept"), "text/event-stream") {
		return false
	}
	if !strings.Contains(req.Header.Get("Content-Type"), "text/html") &&
		!strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		return false
	}
	return true
}