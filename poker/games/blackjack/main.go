package main

import (
	"fmt"
	"strings"

	"github.com/mooncaker816/gophercises/poker/deck"
)

func main() {

	cards := deck.New(deck.Shuffle)
	var player, dealer Hand
	for i := 0; i < 2; i++ {
		for _, hand := range []*Hand{&player, &dealer} {
			*hand = append(*hand, DealOne(&cards))
		}
	}
	var input string
	for input != "s" {
		input = ""
		fmt.Println("Player", player)
		fmt.Println("Dealer", dealer.DealerString())
		fmt.Println("What's your choice? (h)it (s)tand")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			player = append(player, DealOne(&cards))
		}
	}
	mindScore, dScore := dealer.Score()
	for dScore <= 16 || dScore == 17 && dScore != mindScore {
		dealer = append(dealer, DealOne(&cards))
		mindScore, dScore = dealer.Score()
	}
	_, pScore := player.Score()
	fmt.Println("Final Score:")
	fmt.Println("Player", player, "\nScore:", pScore)
	fmt.Println("Dealer", dealer, "\nScore:", dScore)
	switch {
	case pScore > 21:
		fmt.Println("You lose")
	case dScore > 21:
		fmt.Println("You win")
	case pScore > dScore:
		fmt.Println("You win")
	case dScore >= pScore:
		fmt.Println("You lose")
	}
}

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i, c := range h {
		strs[i] = c.String()
	}
	return strings.Join(strs, ", ")
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **"
}

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

func DealOne(cards *deck.Deck) (card deck.Card) {
	if len(*cards) <= 0 {
		*cards = deck.New(deck.Shuffle)
	}
	card, *cards = (*cards)[0], (*cards)[1:]
	return card
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
