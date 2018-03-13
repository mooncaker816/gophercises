package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/ttf"

	"github.com/mooncaker816/gophercises/poker/deck"
	"github.com/mooncaker816/gophercises/poker/room"
	"github.com/veandco/go-sdl2/sdl"
)

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not init sdl: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not init TTF: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window & renderer: %v", err)
	}
	defer w.Destroy()

	if err := drawTitle(r); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(1 * time.Second)
	s, err := room.NewScene(r)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.Destroy()

	if err := start(s, r); err != nil {
		return fmt.Errorf("could not start the game: %v", err)
	}

	time.Sleep(3 * time.Second)
	return nil
}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()
	//font
	f, err := ttf.OpenFont("../../../res/fonts/test.ttf", 800)
	if err != nil {
		return fmt.Errorf("could not open font: %v", err)
	}
	defer f.Close()

	// surface for text
	c := sdl.Color{R: 0, G: 255, B: 0, A: 255}
	s, err := f.RenderUTF8Solid("BlackJack", c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	// texture
	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	// copy to renderer
	if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("could not copy title texture: %v", err)
	}
	r.Present()
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func start(s *room.Scene, r *sdl.Renderer) error {
	//prepare a deck of cards
	cards := deck.New(deck.Shuffle)
	dealer := room.NewPlayer(0, "Dealer", 0, new(room.Hand), room.HIDE_FIRST_CARD)
	player := room.NewPlayer(1, "Player", 1, new(room.Hand), 0)
	s.AddPlayer(player, dealer)

	for i := 0; i < 2; i++ {
		for _, p := range s.Players {
			*p.Hand = append(*p.Hand, deck.DealOneEndless(&cards))
			if err := s.Paint(r); err != nil {
				return err
			}
		}
	}

	var input string
	//input = "s"
	for input != "s" {
		input = ""
		fmt.Println("Player", player.Hand)
		fmt.Println("Dealer", DealerString(dealer))
		fmt.Println("What's your choice? (h)it (s)tand")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			*player.Hand = append(*player.Hand, deck.DealOneEndless(&cards))
			if err := s.Paint(r); err != nil {
				return err
			}
		}
	}
	mindScore, dScore := dealer.Hand.Score()
	for dScore <= 16 || dScore == 17 && dScore != mindScore {
		*dealer.Hand = append(*dealer.Hand, deck.DealOneEndless(&cards))
		if err := s.Paint(r); err != nil {
			return err
		}
		mindScore, dScore = dealer.Hand.Score()
	}
	_, pScore := player.Hand.Score()
	fmt.Println("Final Score:")
	fmt.Println("Player", player.Hand, "\nScore:", pScore)
	fmt.Println("Dealer", dealer.Hand, "\nScore:", dScore)
	switch {
	case pScore > 21 || dScore >= pScore && dScore <= 21:
		fmt.Println("You lose")
		if err := s.UpdateResult(r, false); err != nil {
			return fmt.Errorf("could not update result: %v", err)
		}
		if err := s.Paint(r); err != nil {
			return err
		}
	case dScore > 21 || pScore > dScore && pScore <= 21:
		fmt.Println("You win")
		if err := s.UpdateResult(r, true); err != nil {
			return fmt.Errorf("could not update result: %v", err)
		}
		if err := s.Paint(r); err != nil {
			return err
		}
	}
	return nil
}

// DealerString will only show the first card of dealer for player, others will be masked
func DealerString(p *room.Player) string {
	return (*p.Hand)[0].String() + ", **"
}
