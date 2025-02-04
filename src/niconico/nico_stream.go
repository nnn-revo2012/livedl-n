package niconico

type BinaryStream struct {
	buffer []byte
	offset int
}

func NewBinaryStream() *BinaryStream {
	return &BinaryStream{
		buffer: make([]byte, 0),
		offset: 0,
	}
}

func (bs *BinaryStream) AddBuffer(buf *BinaryStream, data []byte) {
	buf.buffer = append(buf.buffer, data...)
}

func (bs *BinaryStream) ClearBuffer(buf *BinaryStream) {
	buf.buffer = buf.buffer[:0]
	buf.offset = 0
}

func (bs *BinaryStream) decodeVarint(buf *BinaryStream, t int) (*varintResult) {
	a := 0
	o := len(buf.buffer)
	i := 0

	for {
		if o <= t {
			return nil
		}

		b := buf.buffer[t]
		r := (b & 128) != 0
		a |= int(b & 127) << i

		if r {
			t++
			i += 7
		} else {
			break
		}
	}

	//return value, *offset, nil
	return &varintResult{value: a, offset: t}
}

type varintResult struct {
	value  int
	offset int
}

func (bs *BinaryStream) Read(buf *BinaryStream) <- chan []byte {
	results := make(chan []byte)

	go func() {
		defer close(results)
		offset := buf.offset

		for {
			e := bs.decodeVarint(buf, offset)
			if e == nil {
				break
			}

			value, newOffset := e.value, e.offset
			start := newOffset + 1
			rEnd := start + value

			if len(buf.buffer) < rEnd {
				break
			}

			offset = rEnd
			buf.offset = rEnd
			if rEnd - start > 0 {
				binaryData := make([]byte, rEnd-start)
				binaryData = buf.buffer[start:rEnd]
				results <- binaryData
			} else {
				break
			}
		}
	}()

	return results
}
