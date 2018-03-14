package room

import "github.com/veandco/go-sdl2/sdl"

type Game struct {
	Errc   chan error
	Eventc chan sdl.Event
	Status int
}

func NewGame() *Game {
	errc := make(chan error)
	eventc := make(chan sdl.Event)
	return &Game{Errc: errc, Eventc: eventc, Status: 0}
}
