package audio

import (
	"io"
	"math"
	"sync"
	"time"

	oto "github.com/ebitengine/oto/v3"
)

type System struct {
	ctx *oto.Context
	mu  sync.Mutex
}

func NewSystem() *System { return &System{} }

func (s *System) Init(sampleRate int) error {
	opts := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 2,
		Format:       oto.FormatFloat32LE,
		BufferSize:   2048,
	}
	ctx, readyChan, err := oto.NewContext(opts)
	if err != nil {
		return err
	}
	<-readyChan
	s.ctx = ctx
	return nil
}

func (s *System) PlaySine(freq float64, seconds float64, gain float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ctx == nil {
		return nil
	}
	dur := time.Duration(seconds * float64(time.Second))
	p := s.ctx.NewPlayer(NewSineReader(freq, 48000, gain, dur))
	p.Play()
	return nil
}

type SineReader struct {
	freq float64
	rate int
	gain float64
	t    int
	endT int
}

func NewSineReader(freq float64, rate int, gain float64, duration time.Duration) *SineReader {
	return &SineReader{freq: freq, rate: rate, gain: gain, t: 0, endT: int(duration.Seconds() * float64(rate))}
}

func (r *SineReader) Read(p []byte) (int, error) {
	// stereo float32
	if r.t >= r.endT {
		return 0, io.EOF
	}
	// number of frames to generate
	frames := len(p) / 8
	for i := 0; i < frames; i++ {
		v := float32(r.gain * math.Sin(2*math.Pi*float64(r.t)/float64(r.rate)*r.freq))
		b := math.Float32bits(v)
		// left
		p[8*i+0] = byte(b)
		p[8*i+1] = byte(b >> 8)
		p[8*i+2] = byte(b >> 16)
		p[8*i+3] = byte(b >> 24)
		// right
		p[8*i+4] = byte(b)
		p[8*i+5] = byte(b >> 8)
		p[8*i+6] = byte(b >> 16)
		p[8*i+7] = byte(b >> 24)
		r.t++
	}
	return frames * 8, nil
}

func (r *SineReader) Close() error { return nil }
