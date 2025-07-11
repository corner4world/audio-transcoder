//go:build (amd64 && windows) || (amd64 && linux)

package audio_transcoder

/*
#include "opus_encoder.c"
*/
import "C"
import "fmt"

const (
	OpusFrameSize = 960 // 单个声道的样本数
)

func init() {
	RegisterEncoder("OPUS", &OpusEncoder{}, []int{8000, 12000, 16000, 24000, 48000}, 2)
}

type OpusEncoder struct {
	enc             *C.OpusEncoder
	sampleRate      int
	channels        int
	pcmBuf          *PCMBuffer
	sampleSize      int // length is frame_size*channels*sizeof(opus_int16)
	frameDurationMS int
	pktData         []byte
}

func (e *OpusEncoder) Encode(pcm []byte, cb func([]byte)) (int, error) {
	e.pcmBuf.Write(pcm)

	totalEncoded := 0
	e.pcmBuf.ReadTo(func(frame []byte) {
		n := C.opus_encoder_encode(e.enc,
			(*C.uint8_t)(&frame[0]),
			C.int(OpusFrameSize),
			(*C.uint8_t)(&e.pktData[0]),
			C.int(len(e.pktData)))

		if n > 0 {
			cb(e.pktData[:n])
			totalEncoded += int(n)
		}
	})

	return totalEncoded, nil
}

func (e *OpusEncoder) Destroy() {
	if e.enc != nil {
		C.opus_encoder_destroy2(e.enc)
		e.enc = nil
	}
}

func (e *OpusEncoder) ExtraData() []byte {
	return nil
}

func (e *OpusEncoder) Create(sampleRate int, channels int) (int, error) {
	e.enc = C.opus_encoder_create2(C.int(sampleRate), C.int(channels), C.OPUS_APPLICATION_AUDIO)
	if e.enc == nil {
		return -1, fmt.Errorf("failed to create opus encoder")
	}

	e.sampleRate = sampleRate
	e.channels = channels
	e.sampleSize = OpusFrameSize * channels * 2 // 每个样本为 16 位
	e.pcmBuf = NewPCMBuffer(e.sampleSize)
	e.frameDurationMS = 1000 * OpusFrameSize / sampleRate /* / channels*/
	e.pktData = make([]byte, sampleRate*channels)
	return e.sampleSize, nil
}

func (e *OpusEncoder) PacketDurationMS() int {
	return e.frameDurationMS
}
