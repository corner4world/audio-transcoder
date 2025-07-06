//go:build amd64

package audio_transcoder

/*
#include "aac_decoder.c"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type AACDecoder struct {
	decoder unsafe.Pointer
}

func (d *AACDecoder) Decode(frame []byte, pcm []byte) (int, error) {
	n := C.aac_decoder_decode(d.decoder, (*C.uint8_t)(&frame[0]), C.uint32_t(len(frame)), (*C.uint8_t)(&pcm[0]))
	return int(n), nil
}

func (d *AACDecoder) Create(adts []byte, audioConfig []byte) error {
	if len(adts) == 0 && len(audioConfig) == 0 {
		return fmt.Errorf("expect adts or audio config")
	}

	var adtsPtr *C.uint8_t
	var audioConfigPtr *C.uint8_t
	if adts != nil {
		adtsPtr = (*C.uint8_t)(&adts[0])
	}
	if audioConfig != nil {
		audioConfigPtr = (*C.uint8_t)(&audioConfig[0])
	}

	decoder := C.aac_decoder_create(adtsPtr, C.uint32_t(len(adts)), audioConfigPtr, C.uint32_t(len(audioConfig)))
	if decoder == nil {
		return fmt.Errorf("create aac decoder failed")
	}

	d.decoder = decoder
	return nil
}

func (d *AACDecoder) Destroy() {
	if d.decoder != nil {
		C.aac_decoder_destroy(d.decoder)
		d.decoder = nil
	}
}

func (d *AACDecoder) SampleRate() int {
	decoder := (*C.audio_codec_t)(d.decoder)
	return int(decoder.sample_rate)
}

func (d *AACDecoder) Channels() int {
	decoder := (*C.audio_codec_t)(d.decoder)
	return int(decoder.channels)
}
