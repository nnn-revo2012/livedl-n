package message

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

func (bs *BinaryStream) decodeVarint(offset int) (*varintResult) {
	value := 0
	shift := 0
	length := len(bs.buffer)

	for {
		if length <= offset {
			return nil
		}

		byteValue := bs.buffer[offset]
		more := (byteValue & 128) != 0
		value |= int(byteValue&127) << shift

		if more {
			offset++
			shift += 7
		} else {
			break
		}
	}

	//return value, *offset, nil
	return &varintResult{value: value, offset: offset}
}

func (bs *BinaryStream) FromBinary(data []byte) []byte {
	// 必要な処理を追加
	return data
}

type varintResult struct {
	value  int
	offset int
}

func (bs *BinaryStream) Read() <- chan []byte {
	results := make(chan []byte)

	go func() {
		defer close(results)
		offset := 0

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

			//if offset > 0 {
			//	bs.buffer = bs.buffer[offset:]
			//}
		}
	}()

	return results
}
