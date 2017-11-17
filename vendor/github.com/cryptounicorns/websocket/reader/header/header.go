package header

import (
	"io"
	"reflect"
	"unsafe"
)

// ReadHeader reads a frame header from r.
func Read(r io.Reader) (Header, error) {
	// Make slice of bytes with capacity 12 that could hold any header.
	//
	// The maximum header size is 14, but due to the 2 hop reads,
	// after first hop that reads first 2 constant bytes, we could reuse 2 bytes.
	// So 14 - 2 = 12.
	//
	// We use unsafe to stick bts to stack and avoid allocations.
	//
	// Using stack based slice is safe here, cause golang docs for io.Reader
	// says that "Implementations must not retain p".
	// See https://golang.org/pkg/io/#Reader
	var (
		b   [MaxHeaderSize - 2]byte
		buf = *(*[]byte)(
			unsafe.Pointer(
				&reflect.SliceHeader{
					Data: uintptr(unsafe.Pointer(&b)),
					Len:  2,
					Cap:  MaxHeaderSize - 2,
				},
			),
		)

		n     int
		extra int
		err   error
	)

	// Prepare to hold first 2 bytes to choose size of next read.
	n, err = io.ReadFull(r, buf)
	if err != nil {
		return
	}

	h.Fin = buf[0]&bit0 != 0
	h.Rsv = (buf[0] & 0x70) >> 4
	h.OpCode = OpCode(buf[0] & 0x0f)

	if buf[1]&bit0 != 0 {
		h.Masked = true
		extra += 4
	}

	length := buf[1] & 0x7f
	switch {
	case length < 126:
		h.Length = int64(length)

	case length == 126:
		extra += 2

	case length == 127:
		extra += 8

	default:
		err = ErrHeaderLengthUnexpected
		return
	}

	if extra == 0 {
		return
	}

	// Increase len of buf to extra bytes need to read.
	// Overwrite first 2 bytes read before.
	buf = buf[:extra]
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return
	}

	switch {
	case length == 126:
		h.Length = int64(binary.BigEndian.Uint16(buf[:2]))
		buf = buf[2:]

	case length == 127:
		if buf[0]&0x80 != 0 {
			err = ErrHeaderLengthMSB
			return
		}
		h.Length = int64(binary.BigEndian.Uint64(buf[:8]))
		buf = buf[8:]
	}

	if h.Masked {
		copy(h.Mask[:], buf)
	}

	return
}
