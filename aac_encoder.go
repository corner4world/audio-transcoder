//go:build (amd64 && windows) || (amd64 && linux)

package audio_transcoder

/*
#include "aac_encoder.c"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func init() {
	RegisterEncoder("AAC", &AACEncoder{}, []int{8000, 11025, 12000, 16000, 22050, 24000, 32000, 44100, 48000}, 2)
}

type AACEncoder struct {
	encoder         unsafe.Pointer
	sampleRate      int
	sampleSize      int
	channels        int
	pcmBuf          *PCMBuffer
	frameDurationMS int

	pktData   []byte
	extraData []byte
}

func (e *AACEncoder) Encode(pcm []byte, cb func([]byte)) (int, error) {
	e.pcmBuf.Write(pcm)

	var totalEncoded int
	e.pcmBuf.ReadTo(func(frame []byte) {
		n := C.aac_encoder_encode(e.encoder, (*C.uint8_t)(&frame[0]), C.uint32_t(OpusFrameSize), (*C.uint8_t)(&e.pktData[0]), C.uint32_t(len(e.pktData)))
		if n > 0 {
			cb(e.pktData[:n])
			totalEncoded += int(n)
		}
	})

	return totalEncoded, nil
}

func (e *AACEncoder) Create(sampleRate int, channels int, adts int) (int, error) {
	encoder := C.aac_encoder_create(C.uint32_t(sampleRate), C.uint32_t(channels), C.int(adts), C.LOW)
	if encoder == nil {
		return -1, fmt.Errorf("create aac encoder failed")
	}

	encoderC := (*C.audio_codec_t)(encoder)
	e.encoder = encoder
	e.sampleRate = sampleRate
	e.channels = channels
	e.sampleSize = int(encoderC.sample_size) * 2 // 16位采样
	e.pcmBuf = NewPCMBuffer(e.sampleSize)
	e.frameDurationMS = 1000 * int(encoderC.sample_size) / e.sampleRate /* / channels*/
	e.pktData = make([]byte, sampleRate*channels)

	// 获取extra data
	if encoderC.extra_data_size > 0 {
		e.extraData = make([]byte, encoderC.extra_data_size)
		C.memcpy(unsafe.Pointer(&e.extraData[0]), unsafe.Pointer(encoderC.extra_data), C.size_t(encoderC.extra_data_size))
	}
	return e.sampleSize, nil
}

func (e *AACEncoder) ExtraData() []byte {
	return e.extraData
}

func (e *AACEncoder) Destroy() {
	if e.encoder != nil {
		C.aac_encoder_destroy(e.encoder)
		e.encoder = nil
	}
}

func (e *AACEncoder) PacketDurationMS() int {
	return e.frameDurationMS
}
