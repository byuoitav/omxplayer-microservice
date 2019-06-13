package helpers

import "github.com/omxplayer"

var player *omxplayer.Player

const (
	Option   = "--avdict"
	Protocol = "rtsp_transport:tcp"
)

func init() {
	omxplayer.SetUser("pi", "/home/pi")
}

func SwitchStream(streamURL string) {
	if player == nil {
		player, err := omxplayer.New(streamURL, Option, Protocol)
	} else {
		//player.OpenUri(streamURL)
	}
}
