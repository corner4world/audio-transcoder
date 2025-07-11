//go:build (amd64 && windows) || (amd64 && linux)

package audio_transcoder

/*
#cgo CFLAGS: -I./3rd/wrapper -I./3rd/include
#cgo windows,amd64 LDFLAGS: -L./3rd/lib/win64 -lfaad -lfaac -lopus
#cgo linux,amd64 LDFLAGS: -L./3rd/lib/linux64 -lfaad -lfaac -lopus
*/
import "C"
