package audio_transcoder

const (
	PCMA = iota + 1
	PCMU
)

type G711Encoder struct {
	encoderType int
	duration    int
	pktData     []byte
}

func (e *G711Encoder) Encode(pcm []byte, cb func([]byte)) (int, error) {
	length := len(pcm) / 2
	if length == 0 {
		return 0, nil
	} else if length > len(e.pktData) {
		e.pktData = make([]byte, length*2)
	}

	if e.encoderType == PCMA {
		EncodeAlawToBuffer(pcm, e.pktData)
	} else if e.encoderType == PCMU {
		EncodeUlawToBuffer(pcm, e.pktData)
	}

	cb(e.pktData[:length])
	e.duration = length * 1000 / 8000
	return length, nil
}

func (e *G711Encoder) ExtraData() []byte {
	return nil
}

func (e *G711Encoder) Destroy() {

}

func (e *G711Encoder) PacketDurationMS() int {
	return e.duration
}
