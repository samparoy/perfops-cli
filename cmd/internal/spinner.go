package internal

import (
	"fmt"
	"sync"
	"time"
	"unicode/utf8"
)

// Spinner represents the indicator.
type Spinner struct {
	mu       sync.Mutex
	frames   []rune
	length   int
	pos      int
	active   bool
	stopChan chan struct{}
}

// NewSpinner returns a spinner.
func NewSpinner() *Spinner {
	const frames = `|/-\`
	return &Spinner{
		frames:   []rune(frames),
		length:   len([]rune(frames)),
		stopChan: make(chan struct{}, 1),
	}
}

// Start will start the indicator.
func (s *Spinner) Start() {
	if s.active {
		return
	}
	s.active = true
	s.pos = 0
	go func() {
		for {
			for i := 0; i < s.length; i++ {
				select {
				case <-s.stopChan:
					return
				default:
					s.mu.Lock()
					s.erase()
					fmt.Printf("\r%s ", s.next())
					s.mu.Unlock()

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()
}

// Stop will stop the indicator.
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.active {
		s.active = false
		s.erase()
		s.stopChan <- struct{}{}
	}
}

func (s *Spinner) current() string {
	r := s.frames[s.pos%s.length]
	return string(r)
}

func (s *Spinner) next() string {
	r := s.frames[s.pos%s.length]
	s.pos++
	return string(r)
}

func (s *Spinner) erase() {
	n := utf8.RuneCountInString(s.current()) + 1
	for _, c := range []string{"\b", " ", "\b"} {
		for i := 0; i < n; i++ {
			fmt.Printf(c)
		}
	}
}