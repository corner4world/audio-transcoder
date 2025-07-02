package audio_transcoder

type PCMBuffer struct {
	data      []byte
	size      int
	tmpData   []byte // 引用写入的数据, 超过frame size则不拷贝
	frameSize int
}

func (p *PCMBuffer) Write(data []byte) {
	// 如果待编码的pcm数据大于frame size, 优先不拷贝
	if p.size == 0 && len(data) >= p.frameSize {
		p.tmpData = data
	} else {
		p.Copy(data)
	}
}

func (p *PCMBuffer) Copy(data []byte) {
	if cap(p.data) < p.size+len(data) {
		newData := make([]byte, p.size+len(data))
		copy(newData, p.data[:p.size])
		p.data = newData
	}

	copy(p.data[p.size:], data)
	p.size += len(data)
}

func (p *PCMBuffer) ReadTo(cb func([]byte)) {
	if n := len(p.tmpData); n > 0 {
		for offset := p.frameSize; offset <= n; offset += p.frameSize {
			cb(p.tmpData[offset-p.frameSize : offset])
		}

		if r := n % p.frameSize; r > 0 {
			p.Copy(p.tmpData[n-r:])
		}

		p.tmpData = nil
		return
	}

	if p.size < p.frameSize {
		return
	}

	oldSize := p.size
	for offset := p.frameSize; p.size >= p.frameSize; offset += p.frameSize {
		cb(p.data[offset-p.frameSize : offset])
		p.size -= p.frameSize
	}

	if r := p.size % p.frameSize; r > 0 {
		copy(p.data, p.data[oldSize-r:oldSize])
	}
}

func NewPCMBuffer(frameSize int) *PCMBuffer {
	return &PCMBuffer{
		frameSize: frameSize,
	}
}
