package audio_transcoder

import "github.com/general252/g726"

type G726Decoder struct {
	state *g726.G726_state
}

func (d *G726Decoder) Decode(pkt []byte, pcm []byte) (int, error) {
	return d.state.DecodeToBytes(pkt, pcm)
}

func (d *G726Decoder) Destroy() {
}

func (d *G726Decoder) SampleRate() int {
	return 8000
}

func (d *G726Decoder) Channels() int {
	return 1
}

func (d *G726Decoder) Create(rate g726.G726Rate) error {
	d.state = g726.G726_init_state(rate)
	return nil
}
