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

func (bs *BinaryStream) AddBuffer(data []byte) {
	bs.buffer = append(bs.buffer, data...)
}

func (bs *BinaryStream) CheckClearBuffer() {
	if len(bs.buffer) == bs.offset {
		bs.buffer = bs.buffer[:0]
		bs.offset = 0
	}
}

func (bs *BinaryStream) ClearBuffer() {
	bs.buffer = bs.buffer[:0]
	bs.offset = 0
}

func (bs *BinaryStream) decodeVarint(t int) (*varintResult) {
	a := 0
	o := len(bs.buffer)
	i := 0

	for {
		if o <= t {
			return nil
		}

		b := bs.buffer[t]
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

func (bs *BinaryStream) Read() <- chan []byte {
	results := make(chan []byte)

	go func() {
		defer close(results)
		//offset := 0
		offset := bs.offset

		for {
			e := bs.decodeVarint(offset)
			if e == nil {
				break
			}

			value, newOffset := e.value, e.offset
			start := newOffset + 1
			rEnd := start + value

			if len(bs.buffer) < rEnd {
				break
			}

			offset = rEnd
			bs.offset = rEnd
			binaryData := make([]byte, rEnd-start)
			binaryData = bs.buffer[start:rEnd]
			results <- binaryData
		}
	}()

	return results
}
