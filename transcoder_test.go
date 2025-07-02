package audio_transcoder

import (
	"fmt"
	"os"
	"testing"
)

func parseADTSFrameSize(data []byte) int {
	if len(data) < 7 {
		return 0
	}

	return (int(data[3]&0x03) << 11) |
		(int(data[4]) << 3) |
		(int(data[5]&0xE0) >> 5)
}

func readADTSFrame(data []byte) ([]byte, []byte) {
	frameSize := parseADTSFrameSize(data)
	if frameSize == 0 {
		return nil, nil
	} else if len(data) < frameSize {
		return nil, nil
	}

	return data[:frameSize], data[frameSize:]
}

func DecodeAAC(path string) (string, int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	aacDecoder := FindDecoder("AAC")
	if aacDecoder == nil {
		panic("aac decoder not found")
	}

	err = aacDecoder.(*AACDecoder).Create(file[:7], nil)
	if err != nil {
		panic(err)
	}
	defer aacDecoder.Destroy()

	// 创建同名文件, 添加.pcm后缀
	pcmPath := path + ".pcm"
	pcmFos, err := os.Create(pcmPath)
	if err != nil {
		panic(err)
	}

	defer pcmFos.Close()

	pcmBuffer := make([]byte, 1024*1024)
	for offset := 0; offset < len(file); {
		readFrame, _ := readADTSFrame(file[offset:])
		if readFrame == nil {
			break
		}

		pcmN, err := aacDecoder.Decode(readFrame, pcmBuffer)
		if err != nil {
			panic(err)
		} else if pcmN > 0 {
			pcmFos.Write(pcmBuffer[:pcmN])
		}

		offset += len(readFrame)
	}

	return pcmPath, aacDecoder.SampleRate(), aacDecoder.Channels()
}

func EncodeAAC(pcmPath string, sampleRate int, channels int) string {
	file, err := os.ReadFile(pcmPath)
	if err != nil {
		panic(err)
	}
	aacEncoder, err := FindEncoder("AAC", sampleRate, channels)
	if aacEncoder == nil {
		panic(err)
	}

	sampleSize, err := aacEncoder.(*AACEncoder).Create(sampleRate, channels, 1)
	if err != nil {
		panic(err)
	}
	defer aacEncoder.Destroy()

	// 创建同名文件, 添加.aac后缀
	aacFos, err := os.Create(pcmPath + ".aac")
	if err != nil {
		panic(err)
	}

	defer aacFos.Close()

	for offset := 0; offset < len(file); {
		size := sampleSize
		if offset+size > len(file) {
			size = len(file) - offset
		}

		_, _ = aacEncoder.Encode(file[offset:offset+size], func(bytes []byte) {
			aacFos.Write(bytes)
		})
		offset += size
	}

	return pcmPath + ".aac"
}

func DecodeOpus(path string, sampleRate int, channels int) (*OpusDecoder, *os.File) {
	opusDecoder := FindDecoder("OPUS")
	if opusDecoder == nil {
		panic("opus decoder not found")
	}

	err := opusDecoder.(*OpusDecoder).Create(sampleRate, channels)
	if err != nil {
		panic(err)
	}

	// 创建同名文件, 添加.pcm后缀
	pcmFos, err := os.Create(path + ".pcm")
	if err != nil {
		panic(err)
	}

	return opusDecoder.(*OpusDecoder), pcmFos
}

func EncodeOpus(pcmPath string, sampleRate int, channels int, pktCb func([]byte)) {
	file, err := os.ReadFile(pcmPath)
	if err != nil {
		panic(err)
	}
	opusEncoder, err := FindEncoder("OPUS", sampleRate, channels)
	if opusEncoder == nil {
		panic(err)
	}

	sampleSize, err := opusEncoder.(*OpusEncoder).Create(sampleRate, channels)
	if err != nil {
		panic(err)
	}
	defer opusEncoder.Destroy()

	// 创建同名文件, 添加.opus后缀
	opusFos, err := os.Create(pcmPath + ".opus")
	if err != nil {
		panic(err)
	}
	defer opusFos.Close()
	for offset := 0; offset < len(file); {
		size := sampleSize
		if offset+size > len(file) {
			size = len(file) - offset
		}

		_, _ = opusEncoder.Encode(file[offset:offset+size], func(bytes []byte) {
			opusFos.Write(bytes)
			if pktCb != nil {
				pktCb(bytes)
			}
		})

		offset += size
	}
}

func DecodeG711(g711Path string, codeType string) (string, int, int) {
	file, err := os.ReadFile(g711Path)
	if err != nil {
		panic(err)
	}

	g711Decoder := FindDecoder(codeType)
	if g711Decoder == nil {
		panic("g711 decoder not found")
	}

	pcm := make([]byte, len(file)*2)
	n, err := g711Decoder.Decode(file, pcm)
	if err != nil {
		panic(err)
	} else if n != len(file)*2 {
		panic("g711 decode error")
	}
	// 创建同名文件, 添加.pcm后缀
	pcmPath := g711Path + ".pcm"
	pcmFos, err := os.Create(pcmPath)
	if err != nil {
		panic(err)
	}

	defer pcmFos.Close()
	pcmFos.Write(pcm[:n])
	return pcmPath, g711Decoder.SampleRate(), g711Decoder.Channels()
}

func EncodeG711(pcmPath string, codeType string) {
	file, err := os.ReadFile(pcmPath)
	if err != nil {
		panic(err)
	}

	g711Encoder, err := FindEncoder(codeType, 8000, 1)
	if g711Encoder == nil {
		panic(err)
	}

	out, err := os.Create(pcmPath + "." + codeType)
	if err != nil {
		panic(err)
	}

	defer out.Close()
	_, _ = g711Encoder.Encode(file, func(bytes []byte) {
		out.Write(bytes)
	})
}

func testOpusCodec(pcmPath string, sampleRate int, channels int) {
	// 重新解码为pcm, 检查opus解码器是否正常
	opusDecoder, opusPcmFile := DecodeOpus(pcmPath+".opus", sampleRate, channels)
	defer opusPcmFile.Close()
	defer opusDecoder.Destroy()

	pcm := make([]byte, 1024*1024)
	// 编码为opus
	EncodeOpus(pcmPath, sampleRate, channels, func(bytes []byte) {
		n, _ := opusDecoder.Decode(bytes, pcm)
		if n > 0 {
			opusPcmFile.Write(pcm[:n])
		}
	})
}

// 先解码再重新编码
func testAACCodec(aacPath string) (pcmPath string, sampleRate int, channels int) {
	pcmPath, sampleRate, channels = DecodeAAC(aacPath)

	fmt.Println("aac sample rate:", sampleRate)
	fmt.Println("aac channels:", channels)

	// 重新编码为aac, 检查aac编码器是否正常
	EncodeAAC(pcmPath, sampleRate, channels)
	return pcmPath, sampleRate, channels
}

func testG711Codec(g711Path string, codeType string) (string, int, int) {
	pcmPath, sampleRate, channels := DecodeG711(g711Path, codeType)
	// 重新编码为g711, 检查编码器是否正常
	EncodeG711(pcmPath, codeType)
	return pcmPath, sampleRate, channels
}

func TestTranscoder(t *testing.T) {
	t.Run("aac_transcode", func(t *testing.T) {
		aacPath := "../source_files/frxx_48000_2.aac"
		//aacPath := "../source_files/wwzjdy_44100_2.aac"
		pcmPath, sampleRate, channels := testAACCodec(aacPath)
		testOpusCodec(pcmPath, sampleRate, channels)
	})

	t.Run("g711u_transcode", func(t *testing.T) {
		g711aPath := "../source_files/u.g711"
		pcmPath, sampleRate, channels := testG711Codec(g711aPath, "PCMU")

		// 测试aac编解码器, 先编码再重新解码
		aacPath := EncodeAAC(pcmPath, sampleRate, channels)
		DecodeAAC(aacPath)

		// 测试opus编解码器
		testOpusCodec(pcmPath, sampleRate, channels)
	})

	t.Run("g711a_transcode", func(t *testing.T) {
		g711aPath := "../source_files/11111.raw"
		pcmPath, sampleRate, channels := testG711Codec(g711aPath, "PCMA")

		// 测试aac编解码器, 先编码再重新解码
		aacPath := EncodeAAC(pcmPath, sampleRate, channels)
		DecodeAAC(aacPath)

		// 测试opus编解码器
		testOpusCodec(pcmPath, sampleRate, channels)
	})
}
