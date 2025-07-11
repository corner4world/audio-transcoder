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
		// 如果是G711Decoder, 则复制decoderType
		_, ok := decoder.(*G711Decoder)
		if ok {
			newDecoder.(*G711Decoder).decoderType = decoder.(*G711Decoder).decoderType
		}
		return newDecoder
	}

	return nil
}

func init() {
	RegisterDecoder("PCMA", &G711Decoder{decoderType: PCMA})
	RegisterDecoder("PCMU", &G711Decoder{decoderType: PCMU})
	RegisterDecoder("G726", &G726Decoder{})
}
