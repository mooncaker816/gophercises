package room

import (
	"fmt"
	"strings"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/mooncaker816/gophercises/poker/deck"
)

// configs for Player
const (
	HIDE_FIRST_CARD = 1 << iota
)

// Player is a game player or AI
type Player struct {
	ID     int
	Name   string
	Pos    int
	Hand   *Hand
	Config int
}

// NewPlayer will create a new player
func NewPlayer(id int, name string, pos int, h *Hand, f int) (*Player, error) {
	return &Player{ID: id, Name: name, Pos: pos, Hand: h, Config: f}, nil
}

func (p Player) paint(r *sdl.Renderer, s *Scene) error {
	num := len(*p.Hand)
	leftBase := 400 - (138+30*num-1)/2
	for i, absrank := range p.Hand.C2Rs() {
		rect := &sdl.Rect{X: int32(i*30 + leftBase), Y: int32(300*p.Pos + (300-200)/2), W: 138, H: 200}
		if p.Config&HIDE_FIRST_CARD == 1 && i == 0 && s.Table.result == nil {
			absrank = 0
		}
		if err := r.Copy(s.cards[absrank], nil, rect); err != nil {
			return fmt.Errorf("could not copy card texture: %v", err)
		}
	}
	return nil
}

func (p Player) Hit(r *sdl.Renderer, s *Scene) {
	*p.Hand = append(*p.Hand, deck.DealOneEndless(s.Table.Deck))
}

// Hand holds the current cards of a player
type Hand []deck.Card

// C2Rs convert the cards to abs ranks in slice
func (h Hand) C2Rs() []int {
	var ret []int
	for _, c := range h {
		ret = append(ret, deck.AbsRank1(c))
	}
	return ret
}

// String will format the cards on hand to string
func (h Hand) String() string {
	strs := make([]string, len(h))
	for i, c := range h {
		strs[i] = c.String()
	}
	return strings.Join(strs, ", ")
}

// Score will return the min score and max score according to the cards on hand, J,Q,K = 10, A = 1 or 11
func (h Hand) Score() (int, int) {
	var minscore int
	var ace bool
	for _, c := range h {
		minscore += min(int(c.Rank), 10)
		if c.Rank == deck.Ace {
			ace = true
		}
	}
	if minscore > 11 {
		return minscore, minscore
	}
	if ace {
		return minscore, minscore + 10
	}
	return minscore, minscore
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
