package packetsender

import (
	"errors"
	"github.com/colek42/ffgopeg/avcodec"
	"github.com/colek42/ffgopeg/avformat"
	"github.com/colek42/ffgopeg/avutil"
	"log"
	"time"
)

func init() {
	avformat.RegisterAll()
	avformat.NetworkInit()
}

func OpenStream(uri string, packetChan chan tspacket) {
	formatCtx, code1 := avformat.OpenInput(uri, nil, nil)
	if !code1.Ok() {
		log.Printf("%v", code1.Error())
	}
	defer formatCtx.Close()

	//formatCtx.FindStreamInfo(nil)

	videoSteamIndex, err := findFirstVideoStream(formatCtx)
	if err != nil {
		log.Printf("Error")
	}
	log.Printf("VideoStream Index: %v", videoSteamIndex)

	codec := avcodec.FindDecoder(formatCtx.Streams()[videoSteamIndex].CodecPar().CodecID())
	codecCtx := avcodec.NewCodecContext(codec)

	frame := avutil.NewFrame()
	defer frame.Free()

	var packet avcodec.Packet
	packet.Init()

	for {
		err := formatCtx.ReadFrame(&packet)
		if err.IsOneOf(avutil.AVERROR_EOF()) {
			break
		}

		//for now we only care about video packets
		if packet.StreamIndex() != videoSteamIndex {
			packet.Unref()
			continue
		}

		code := codecCtx.SendPacket(&packet)
		if code.Ok() {
			//frame.Unref()
		}

		if len(packetChan) < 5 {
			packetChan <- tspacket{
				data:      packet.GetData(),
				pts:       packet.Pts(),
				dts:       packet.Dts(),
				timeStamp: time.Now().UnixNano(),
			}
		}
		packet.Unref()
	}

}

func findFirstVideoStream(ctx *avformat.FormatContext) (int, error) {
	for i, s := range ctx.Streams() {
		if s.CodecPar().CodecType() == avutil.AVMEDIA_TYPE_VIDEO {
			return i, nil
		}
	}

	return -1, errors.New("Could Not Find Video Stream")
}
