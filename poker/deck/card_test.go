package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Spade, King})
	fmt.Println(Card{Heart, Ace})
	fmt.Println(BigJoker)
	fmt.Println(LittleJoker)
	fmt.Println(Card{Club, Ten})
	fmt.Println(Card{Diamond, Six})

	// Output:
	// ♠K
	// ♥A
	// BigJoker
	// LittleJoker
	// ♣10
	// ♦6
}

func TestNew(t *testing.T) {
	cards := New()
	fmt.Println(cards)
	cards1 := New(Add(BigJoker, LittleJoker))
	fmt.Println(cards1)
	cards2 := New(SortRankSuit)
	fmt.Println(cards2)
	cards3 := New(SortSuitRank)
	fmt.Println(cards3)
	cards4 := New(Add(BigJoker, LittleJoker), SortRankSuit)
	fmt.Println(cards4)
	cards5 := New(Add(BigJoker, LittleJoker), SortSuitRank)
	fmt.Println(cards5)
	cards6 := New(Add(BigJoker, LittleJoker), Sort(less1))
	fmt.Println(cards6)
	cards7 := New(Add(BigJoker, LittleJoker), Shuffle)
	fmt.Println(cards7)
	f := func(card Card) bool {
		if card.Rank == King {
			return true
		}
		return false
	}
	cards8 := New(Add(BigJoker, LittleJoker), Filter(f), Shuffle)
	fmt.Println(cards8)
	cards9 := New(Add(BigJoker, LittleJoker), Multiple(3), Shuffle)
	fmt.Println(cards9, len(cards9))
}
