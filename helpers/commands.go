package helpers

import (
	"fmt"

	"github.com/godbus/dbus"
)

const (
	destination  = "org.mpris.MediaPlayer2.omxplayer"
	path         = "/org/mpris/MediaPlayer2"
	playerPrefix = "org.mpris.MediaPlayer2.Player"

	methodGetSource = playerPrefix + ".GetSource"
	methodStop      = playerPrefix + ".Stop"
	methodOpenURI   = playerPrefix + ".OpenUri"
	methodMute      = playerPrefix + ".Mute"
	methodUnmute    = playerPrefix + ".Unmute"
	methodPropGet   = "org.freedesktop.DBus.Properties.Get"
	methodPropSet   = "org.freedesktop.DBus.Properties.Set"

	cmdPlayback = "PlaybackStatus"
	cmdVolume   = "Volume"
)

//GetStream returns the url of the stream currently playing
func GetStream(conn *dbus.Conn) (string, error) {
	var stream string
	err := conn.Object(destination, path).Call(methodGetSource, 0).Store(&stream)
	if err != nil {
		return "", fmt.Errorf("Failed to get stream url | %s", err.Error())
	}
	return stream, err
}

//GetPlaybackStatus returns the status of the player
func GetPlaybackStatus(conn *dbus.Conn) (string, error) {
	var playback string
	err := conn.Object(destination, path).Call(methodPropGet, 0, playerPrefix, cmdPlayback).Store(&playback)
	if err != nil {
		return "", fmt.Errorf("Failed to get playback status | %s", err.Error())
	}
	return playback, err
}

//StopStream quits the omxplayer
func StopStream(conn *dbus.Conn) error {
	err := conn.Object(destination, path).Call(methodStop, 0).Err
	if err != nil {
		return fmt.Errorf("Failed to stop stream | %s", err.Error())
	}
	return err
}

//SwitchStream switches player output to a new stream
func SwitchStream(conn *dbus.Conn, streamURL string) error {
	err := conn.Object(destination, path).Call(methodOpenURI, 0, streamURL).Err
	if err != nil {
		return fmt.Errorf("Failed to switch to stream %s | %s", streamURL, err.Error())
	}
	return err
}

//VolumeControl always returns the current volume and optionally can change the volume
func VolumeControl(conn *dbus.Conn, volume ...float64) (currVolume float64, err error) {
	if len(volume) == 0 {
		err := conn.Object(destination, path).Call(methodPropSet, 0, playerPrefix, cmdVolume).Store(&currVolume)
		if err != nil {
			err = fmt.Errorf("Failed to get volume | %s", err.Error())
		}
	} else {
		err = conn.Object(destination, path).Call(methodPropSet, 0, playerPrefix, cmdVolume, volume[0]).Store(&currVolume)
		if err != nil {
			err = fmt.Errorf("Failed to set volume | %s", err.Error())
		}
	}
	return
}

//Mute mutes the current player output
func Mute(conn *dbus.Conn) error {
	err := conn.Object(destination, path).Call(methodMute, 0).Err
	if err != nil {
		return fmt.Errorf("Failed to mute stream | %s", err.Error())
	}
	return err
}

//Unmute unmutes the current player output
func Unmute(conn *dbus.Conn) error {
	err := conn.Object(destination, path).Call(methodUnmute, 0).Err
	if err != nil {
		return fmt.Errorf("Failed to unmute stream | %s", err.Error())
	}
	return err
}
