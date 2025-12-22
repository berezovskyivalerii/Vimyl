package player

import (
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

type AudioPlayer struct {
	streamer beep.StreamSeekCloser
	ctrl     *beep.Ctrl
	format   beep.Format

	IsPlaying bool
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{
		IsPlaying: false,
	}
}

func (p *AudioPlayer) PlayFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		f.Seek(0, 0)
		streamer, format, err = wav.Decode(f)
		if err != nil {
			return err
		}
	}

	if p.format.SampleRate != format.SampleRate {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	}

	if p.streamer != nil {
		p.streamer.Close()
	}

	p.streamer = streamer
	p.format = format

	p.ctrl = &beep.Ctrl{Streamer: p.streamer, Paused: false}

	speaker.Play(p.ctrl)
	p.IsPlaying = true

	return nil
}

func (p *AudioPlayer) TogglePause() {
	if p.ctrl == nil {
		return
	}

	speaker.Lock()
	p.ctrl.Paused = !p.ctrl.Paused
	p.IsPlaying = !p.ctrl.Paused
	speaker.Unlock()
}

// GetProgress returns position and duration in seconds
func (p *AudioPlayer) GetProgress() (position float64, duration float64) {
	if p.streamer == nil || p.format.SampleRate == 0 {
		return 0, 0
	}

	// Lock speaker to prevent race conditions while reading position
	speaker.Lock()
	pos := p.streamer.Position()
	len := p.streamer.Len()
	speaker.Unlock()

	// Convert samples to seconds
	position = p.format.SampleRate.D(pos).Seconds()
	duration = p.format.SampleRate.D(len).Seconds()
	return
}

// Seek moves the song to a specific percentage (0.0 to 1.0)
func (p *AudioPlayer) Seek(percentage float64) {
	if p.streamer == nil {
		return
	}
	speaker.Lock()
	length := p.streamer.Len()
	pos := int(float64(length) * percentage)
	// Clamp position to avoid crashes
	if pos >= length {
		pos = length - 1
	}
	if pos < 0 {
		pos = 0
	}
	p.streamer.Seek(pos)
	speaker.Unlock()
}

