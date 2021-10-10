package instance

import (
	"bufio"
	"encoding/binary"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"io"
	"os/exec"
	"strconv"
)

const (
	channels  int = 2     // 1 for mono, 2 for stereo
	frameRate int = 48000 // audio sampling rate
	frameSize int = 960   // uint16 size of each audio frame
)

// PlayAudioStream transcodes and plays audio over a voice connection
func PlayAudioStream(vc *discordgo.VoiceConnection, stream io.ReadCloser, stop <-chan bool) error {
	// Command to transcode incoming stream to raw PCM
	transcode := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")

	transcode.Stdin = stream

	transcodeOut, err := transcode.StdoutPipe()

	if err != nil {
		println("Failed to retrieve stdout pipe for command")

		return err
	}

	reader := bufio.NewReaderSize(transcodeOut, 16384)

	err = transcode.Start()

	if err != nil {
		println("Transcode command failed to start")

		return err
	}

	// Cleanup stream and process on end
	defer stream.Close()
	defer transcode.Process.Kill()

	// Kill the process on stop
	go func() {
		<-stop
		err = transcode.Process.Kill()
	}()

	err = vc.Speaking(true)

	if err != nil {
		println("Was unable to start speaking in channel")

		return err
	}

	// Stop speaking once the method ends
	defer func() {
		err := vc.Speaking(false)
		if err != nil {
			println("Was unable to stop speaking in channel")
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	// When the instance has finished
	finished := make(chan bool)

	go func() {
		dgvoice.SendPCM(vc, send)
		finished <- true
	}()

	for {
		// Read from transcode output
		audioFrame := make([]int16, frameSize*channels)

		err = binary.Read(reader, binary.LittleEndian, &audioFrame)

		// EOF or Unexpected EOF is just the end of the stream
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}

		if err != nil {
			println("Encountered an error while reading audio stream", err.Error())

			return err
		}

		select {
		case send <- audioFrame:
		case <-finished:
			return nil
		}
	}
}
