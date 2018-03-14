package room

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var suits = [...]string{"♦", "♣", "♥", "♠"}
var ranks = [...]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

// Scene holds background,cards...
type Scene struct {
	bg    *sdl.Texture
	cards []*sdl.Texture
	Table *Table
	Game  *Game
}

// NewScene creates a scene with background and basic 54 cards' textures
func NewScene(r *sdl.Renderer) (*Scene, error) {
	bg, err := img.LoadTexture(r, "../../../res/img/bg.jpg")
	if err != nil {
		return nil, err
	}
	var cards []*sdl.Texture
	card, err := img.LoadTexture(r, "../../../res/img/cards/back.png")
	if err != nil {
		return nil, err
	}
	cards = append(cards, card)
	for i := 0; i < 4; i++ {
		for j := 0; j < 13; j++ {
			path := "../../../res/img/cards/" + suits[i] + ranks[j] + ".png"
			card, err := img.LoadTexture(r, path)
			if err != nil {
				return nil, err
			}
			cards = append(cards, card)
		}
	}
	card, err = img.LoadTexture(r, "../../../res/img/cards/joker1.png")
	if err != nil {
		return nil, err
	}
	cards = append(cards, card)
	card, err = img.LoadTexture(r, "../../../res/img/cards/joker2.png")
	if err != nil {
		return nil, err
	}
	cards = append(cards, card)
	t := new(Table)
	g := NewGame()
	return &Scene{bg: bg, cards: cards, Table: t, Game: g}, nil
}

// Paint will copy the texture to renderer
func (s *Scene) Paint(r *sdl.Renderer) error {
	r.Clear()
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background texture: %v", err)
	}

	if err := s.Table.paint(r, s); err != nil {
		return fmt.Errorf("could not paint table: %v", err)
	}

	r.Present()
	time.Sleep(800 * time.Millisecond)
	return nil
}

// Destroy will delete all the textures
func (s *Scene) Destroy() {
	s.bg.Destroy()
	for i := range s.cards {
		s.cards[i].Destroy()
	}
	s.Table.destroy()
}
