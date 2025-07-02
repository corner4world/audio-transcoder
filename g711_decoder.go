package audio_transcoder

type G711Decoder struct {
	decoderType int
}

func (d *G711Decoder) Decode(data []byte, pcm []byte) (int, error) {
	length := len(data)
	if length == 0 {
		return 0, nil
	}

	if d.decoderType == PCMA {
		DecodeAlawToBuffer(data, pcm)
	} else if d.decoderType == PCMU {
		DecodeUlawToBuffer(data, pcm)
	}

	return length * 2, nil
}

func (d *G711Decoder) Destroy() {

}

func (d *G711Decoder) SampleRate() int {
	return 8000
}

func (d *G711Decoder) Channels() int {
	return 1
}
