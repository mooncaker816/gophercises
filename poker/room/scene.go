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
	bg      *sdl.Texture
	cards   []*sdl.Texture
	result  *sdl.Texture
	Players []*Player
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
	var players []*Player
	return &Scene{bg: bg, cards: cards, Players: players}, nil
}

// Paint will copy the texture to renderer
func (s *Scene) Paint(r *sdl.Renderer) error {
	r.Clear()
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background texture: %v", err)
	}

	for _, p := range s.Players {
		num := len(*p.Hand)
		leftBase := 400 - (138+30*num-1)/2
		for i, absrank := range p.Hand.C2Rs() {
			rect := &sdl.Rect{X: int32(i*30 + leftBase), Y: int32(300*p.Pos + (300-200)/2), W: 138, H: 200}
			if p.Config&HIDE_FIRST_CARD == 1 && i == 0 && s.result == nil {
				absrank = 0
			}
			if err := r.Copy(s.cards[absrank], nil, rect); err != nil {
				return fmt.Errorf("could not copy card texture: %v", err)
			}
		}
	}

	if s.result != nil {
		rect := &sdl.Rect{X: 200, Y: 150, W: 400, H: 300}
		if err := r.Copy(s.result, nil, rect); err != nil {
			return fmt.Errorf("could not copy card texture: %v", err)
		}
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
}

// AddPlayer will add new palyer to scene
func (s *Scene) AddPlayer(ps ...*Player) {
	for _, p := range ps {
		s.Players = append(s.Players, p)
	}
}

// UpdateResult will update the result for final rendering
func (s *Scene) UpdateResult(r *sdl.Renderer, win bool) error {
	var path string
	if win {
		path = "../../../res/img/youwin.jpg"
	} else {
		path = "../../../res/img/youlose.jpg"
	}
	res, err := img.LoadTexture(r, path)
	if err != nil {
		return fmt.Errorf("could not load result: %v", err)
	}
	s.result = res
	return nil
}
