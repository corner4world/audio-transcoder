package audio_transcoder

import (
	"fmt"
	"reflect"
)

var (
	encoders = make(map[string]struct {
		encoder     Encoder
		sampleRate  []int
		maxChannels int
	})
)

type Encoder interface {
	Encode(pcm []byte, cb func([]byte)) (int, error)

	ExtraData() []byte

	Destroy()

	PacketDurationMS() int
}

func RegisterEncoder(name string, encoder Encoder, sampleRate []int, maxChannels int) {
	if _, ok := encoders[name]; ok {
		panic("encoder already registered with name " + name)
	}

	encoders[name] = struct {
		encoder     Encoder
		sampleRate  []int
		maxChannels int
	}{encoder: encoder, sampleRate: sampleRate, maxChannels: maxChannels}
}

func FindEncoder(name string, sampleRate, channels int) (Encoder, error) {
	encoder, ok := encoders[name]
	if !ok {
		return nil, fmt.Errorf("encoder not found with name %s", name)
	}

	if channels > encoder.maxChannels {
		return nil, fmt.Errorf("encoder %s does not support %d channels", name, channels)
	}

	for _, rate := range encoder.sampleRate {
		if rate == sampleRate {
			t := reflect.TypeOf(encoder.encoder)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}

			newDecoder := reflect.New(t).Interface().(Encoder)
			// 如果是G711Encoder, 则复制encoderType
			_, ok := encoder.encoder.(*G711Encoder)
			if ok {
				newDecoder.(*G711Encoder).encoderType = encoder.encoder.(*G711Encoder).encoderType
			}
			return newDecoder, nil
		}
	}

	return nil, fmt.Errorf("encoder %s does not support %d sample rate", name, sampleRate)
}

func init() {
	RegisterEncoder("PCMA", &G711Encoder{encoderType: PCMA}, []int{8000}, 2)
	RegisterEncoder("PCMU", &G711Encoder{encoderType: PCMU}, []int{8000}, 2)
	RegisterEncoder("G726", &G726Encoder{}, []int{8000}, 1)
}
