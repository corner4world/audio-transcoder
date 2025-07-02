package audio_transcoder

import (
	"reflect"
)

var (
	decoders = make(map[string]Decoder)
)

type Decoder interface {
	Decode(pkt []byte, pcm []byte) (int, error)

	Destroy()

	SampleRate() int

	Channels() int
}

func RegisterDecoder(name string, decoder Decoder) {
	if _, ok := decoders[name]; ok {
		panic("decoder already registered with name " + name)
	}

	decoders[name] = decoder
}

func FindDecoder(name string) Decoder {
	if decoder, ok := decoders[name]; ok {
		t := reflect.TypeOf(decoder)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		newDecoder := reflect.New(t).Interface().(Decoder)
		return newDecoder
	}

	return nil
}

func init() {
	RegisterDecoder("AAC", &AACDecoder{})
	RegisterDecoder("OPUS", &OpusDecoder{})
}
