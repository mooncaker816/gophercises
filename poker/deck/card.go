package deck

import (
	"math/rand"
	"sort"
	"time"
)

// Suit values
const (
	_ Suit = iota
	Diamond
	Club
	Heart
	Spade
	Joker
)

var suits = [...]Suit{Spade, Heart, Club, Diamond}

// Suit is the type of the card
type Suit uint8

func (s Suit) String() string {
	switch s {
	case Spade:
		return "♠"
	case Diamond:
		return "♦"
	case Club:
		return "♣"
	case Heart:
		return "♥"
	default:
		return ""
	}
}

// Rank values
const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	minRank = Ace
	maxRank = King
)

// Rank is num of the card
type Rank uint8

func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return ""
	}
}

// Card stands for the card
type Card struct {
	Suit
	Rank
}

// BigJoker is a special card of big joker
var BigJoker = Card{Joker, Rank(15)}

// LittleJoker is a special card of little joker
var LittleJoker = Card{Joker, Rank(14)}

func (c Card) String() string {
	if c == BigJoker {
		return "BigJoker"
	}
	if c == LittleJoker {
		return "LittleJoker"
	}
	return c.Suit.String() + c.Rank.String()
}

// Deck is a Set of cards
type Deck []Card

// Option is a function to modify the deck
// 1. Add(card ...card) returns an option to add the card specified to the deck
// 2. SortSuitRank is an option to sort the deck by Suit first then Rank
// 3. SortRankSuit is an option to sort the deck by Rank first then Suit
// 4. Sort(less func(deck Deck) func(i, j int) bool) returns an option allow you to sort the deck in your own way
// 5. Shuffle is an option to shuffle the deck
// 6. Filter(f func(card Card) bool) returns an option allow you to filter specific cards which satisfies f
// 7. Multiple(n int) returns an option allow you to build the deck in double,trible...
type Option func(deck *Deck)

//Add card to deck
func Add(card ...Card) Option {
	return func(deck *Deck) {
		*deck = append(*deck, card...)
	}
}

// New Create a deck of cards without Jokers
func New(opts ...Option) Deck {
	var deck Deck
	for _, s := range suits {
		for i := 1; i <= 13; i++ {
			deck = append(deck, Card{s, Rank(i)})
		}
	}
	for _, opt := range opts {
		opt(&deck)
	}
	return deck
}

func absRank1(c Card) int {
	return int(c.Suit)*13 + int(c.Rank)
}

func absRank2(c Card) int {
	return int(c.Rank)*4 + int(c.Suit)
}

func less1(deck Deck) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank1(deck[i]) > absRank1(deck[j])
	}
}

func less2(deck Deck) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank2(deck[i]) > absRank2(deck[j])
	}
}

// SortSuitRank sorting by Suit first then Rank
func SortSuitRank(deck *Deck) {
	sort.Slice(*deck, less1(*deck))
}

//SortRankSuit sorting by Rank first then Suit
func SortRankSuit(deck *Deck) {
	sort.Slice(*deck, less2(*deck))
}

// Sort use customer less function to do the sorting
func Sort(less func(deck Deck) func(i, j int) bool) Option {
	return func(deck *Deck) {
		sort.Slice(*deck, less(*deck))
	}
}

// Shuffle will shuffle the deck
func Shuffle(deck *Deck) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Shuffle(len(*deck), swap(*deck))
}

func swap(deck Deck) func(i, j int) {
	return func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	}
}

// Filter will remove the cards which satisfies f
func Filter(f func(card Card) bool) Option {
	return func(deck *Deck) {
		for i, c := range *deck {
			if f(c) {
				*deck = append((*deck)[:i], (*deck)[i+1:]...)
			}
		}
	}
}

// Multiple will make the deck = deck * n
func Multiple(n int) Option {
	return func(deck *Deck) {
		base := *deck
		for i := 0; i < n-1; i++ {
			*deck = append(*deck, base...)
		}
	}
}
