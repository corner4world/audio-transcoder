package audio_transcoder

import "github.com/lkmio/g726"

type G726Encoder struct {
	state    *g726.G726_state
	pktData  []byte
	duration int
}

func (e *G726Encoder) Encode(pcm []byte, pktCb func(bytes []byte)) (int, error) {
	if len(e.pktData) < len(pcm) {
		e.pktData = make([]byte, len(pcm))
	}

	n, err := e.state.EncodeToBytes(pcm, e.pktData)
	if n > 0 {
		e.duration = len(pcm) / 2 * 1000 / 8000
		pktCb(e.pktData[:n])
	}
	return n, err
}

func (e *G726Encoder) Destroy() {
	e.state = nil
}

func (e *G726Encoder) ExtraData() []byte {
	return nil
}

func (e *G726Encoder) PacketDurationMS() int {
	return e.duration
}

func (e *G726Encoder) Create(rate g726.G726Rate) error {
	e.state = g726.G726_init_state(rate)
	return nil
}
