package opcode

// OpCode represents operation code.
type OpCode byte

// Operation codes defined by specification.
// See https://tools.ietf.org/html/rfc6455#section-5.2
const (
	Continuation OpCode = 0x0
	Text                = 0x1
	Binary              = 0x2
	Close               = 0x8
	Ping                = 0x9
	Pong                = 0xa
)

// IsControl checks wheter the c is control operation code.
// See https://tools.ietf.org/html/rfc6455#section-5.5
func (c OpCode) IsControl() bool {
	// RFC6455: Control frames are identified by opcodes where
	// the most significant bit of the opcode is 1.
	//
	// Note that OpCode is only 4 bit length.
	return c&0x8 != 0
}

// IsData checks wheter the c is data operation code.
// See https://tools.ietf.org/html/rfc6455#section-5.6
func (c OpCode) IsData() bool {
	// RFC6455: Data frames (e.g., non-control frames) are identified by opcodes
	// where the most significant bit of the opcode is 0.
	//
	// Note that OpCode is only 4 bit length.
	return c&0x8 == 0
}

// IsReserved checks wheter the c is reserved operation code.
// See https://tools.ietf.org/html/rfc6455#section-5.2
func (c OpCode) IsReserved() bool {
	// RFC6455:
	// %x3-7 are reserved for further non-control frames
	// %xB-F are reserved for further control frames
	return (0x3 <= c && c <= 0x7) || (0xb <= c && c <= 0xf)
}
