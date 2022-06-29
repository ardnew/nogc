package nogc

import "io"

// Buffer is the interface that groups read/write operations on types composed
// of a fixed-length slice of bytes.
type Buffer interface {
	io.ReadWriter // io.Reader, io.Writer

	io.ReaderFrom
	io.WriterTo

	io.ByteScanner // io.ByteReader
	io.ByteWriter
}
