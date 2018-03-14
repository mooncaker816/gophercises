package room

import (
	"fmt"

	"github.com/mooncaker816/gophercises/poker/deck"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Table struct {
	result  *sdl.Texture
	Players []*Player
	Deck    *deck.Deck
}

// AddDeck create new deck as desired
func (t *Table) AddDeck(n int) deck.Deck {
	d := deck.New(deck.Multiple(n), deck.Shuffle)
	t.Deck = &d
	return d
}

// AddPlayer will add new palyer to scene
func (t *Table) AddPlayer(ps ...*Player) {
	for _, p := range ps {
		t.Players = append(t.Players, p)
	}
}

// UpdateResult will update the result for final rendering
func (t *Table) UpdateResult(r *sdl.Renderer, win bool) error {
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
	t.result = res
	return nil
}

func (t *Table) paint(r *sdl.Renderer, s *Scene) error {
	for _, p := range s.Table.Players {
		if err := p.paint(r, s); err != nil {
			return fmt.Errorf("could not paint player: %v", err)
		}
	}
	if s.Table.result != nil {
		rect := &sdl.Rect{X: 200, Y: 150, W: 400, H: 300}
		if err := r.Copy(s.Table.result, nil, rect); err != nil {
			return fmt.Errorf("could not copy card texture: %v", err)
		}
	}
	return nil
}

func (t *Table) destroy() {
	t.result.Destroy()
}
