package audio_transcoder

/*
#include "opus_decoder.c"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type OpusDecoder struct {
	dec        unsafe.Pointer
	sampleRate int
	channels   int
}

func (d *OpusDecoder) Decode(pkt []byte, pcm []byte) (int, error) {
	n := C.opus_decoder_decode(d.dec, (*C.uint8_t)(&pkt[0]), C.int(len(pkt)), (*C.uint8_t)(&pcm[0]), C.int(OpusFrameSize))
	return int(n), nil
}

func (d *OpusDecoder) Destroy() {
	if d.dec != nil {
		C.opus_decoder_destroy2(d.dec)
		d.dec = nil
	}
}

func (d *OpusDecoder) Create(sampleRate, channel int) error {
	decoder := C.opus_decoder_create2(C.int(sampleRate), C.int(channel))
	if decoder == nil {
		return fmt.Errorf("failed to create opus decoder")
	}

	d.dec = decoder
	d.sampleRate = sampleRate
	d.channels = channel
	return nil
}

func (d *OpusDecoder) SampleRate() int {
	return d.sampleRate
}

func (d *OpusDecoder) Channels() int {
	return d.channels
}
