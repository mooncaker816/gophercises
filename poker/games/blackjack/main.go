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

const (
	SDL_LEFT_BUTTON  = 1
	SDL_RIGHT_BUTTON = 3
)

func run() error {
	var w *sdl.Window
	var r *sdl.Renderer
	var s *room.Scene
	var err error
	sdl.Do(func() {
		err = sdl.Init(sdl.INIT_EVERYTHING)
	})
	if err != nil {
		return fmt.Errorf("could not init sdl: %v", err)
	}

	defer func() {
		sdl.Do(func() {
			sdl.Quit()
		})
	}()

	sdl.Do(func() {
		err = ttf.Init()
	})
	if err != nil {
		return fmt.Errorf("could not init TTF: %v", err)
	}

	defer func() {
		sdl.Do(func() {
			ttf.Quit()
		})
	}()

	sdl.Do(func() {
		w, r, err = sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	})

	if err != nil {
		return fmt.Errorf("could not create window & renderer: %v", err)
	}

	defer func() {
		sdl.Do(func() {
			w.Destroy()
		})
	}()

	sdl.Do(func() {
		err = drawTitle(r)
	})
	if err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(1 * time.Second)

	sdl.Do(func() {
		s, err = room.NewScene(r)
	})

	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}

	defer func() {
		sdl.Do(func() {
			s.Destroy()
		})
	}()

	g := room.NewGame()
	// start the game in a go routine and close the start channel once it's done
	go start(r, s, g)

	//events := make(chan sdl.Event)
	// reveive events and handle them
	go func() {
		defer close(g.Errc)
		log.Println("start receive event")
		for {
			select {
			case e := <-g.Eventc:
				if quit := handleEvent(e, r, s, g); quit {
					return
				}
			}
		}
	}()
	// listening user event in main thread to send to events channel
	sdl.Do(func() {
		log.Println("start listen event")
		for {
			select {
			case g.Eventc <- sdl.WaitEvent():
			case err = <-g.Errc:
				return
			}
		}
	})
	//time.Sleep(3 * time.Second)
	return err
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
	var err error
	sdl.Main(func() { err = run() })
	if err != nil {
		log.Fatal(err)
	}
}

func start(r *sdl.Renderer, s *room.Scene, g *room.Game) {
	//prepare a deck of cards
	cards := s.Table.AddDeck(1)
	// prepare players
	dealer, err := room.NewPlayer(0, "Dealer", 0, new(room.Hand), room.HIDE_FIRST_CARD)
	if err != nil {
		g.Errc <- fmt.Errorf("could not create dealer: %v", err)
		return
	}
	player, err := room.NewPlayer(1, "Player", 1, new(room.Hand), 0)
	if err != nil {
		g.Errc <- fmt.Errorf("could not create player: %v", err)
		return
	}
	s.Table.AddPlayer(player, dealer)

	for i := 0; i < 2; i++ {
		for _, p := range s.Table.Players {
			*p.Hand = append(*p.Hand, deck.DealOneEndless(&cards))
			if err := s.Paint(r); err != nil {
				g.Errc <- err
				return
			}
		}
	}
	fmt.Println("Player", player.Hand)
	fmt.Println("Dealer", DealerString(dealer))
	g.Status = 1
}

// DealerString will only show the first card of dealer for player, others will be masked
func DealerString(p *room.Player) string {
	return "**, " + (*p.Hand)[1:].String()
}

func handleEvent(event sdl.Event, r *sdl.Renderer, s *room.Scene, g *room.Game) bool {

	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		fmt.Printf("button %+v\n", e)
		if g.Status == 1 {
			switch e.Button {
			case SDL_LEFT_BUTTON:
				if e.State == sdl.PRESSED { // player's turn
					s.Table.Players[0].Hit(r, s)
					if err := s.Paint(r); err != nil {
						g.Errc <- err
						return true
					}
					if _, tmp := s.Table.Players[0].Hand.Score(); tmp > 21 {
						fmt.Println("You lose")
						if err := gameover(r, s, g, false); err != nil {
							g.Errc <- err
							return true
						}
					}
					fmt.Println("Player", s.Table.Players[0].Hand)
					fmt.Println("Dealer", DealerString(s.Table.Players[1]))
					fmt.Println("What's your choice? (h)it (s)tand")
				}
			case SDL_RIGHT_BUTTON:
				if e.State == sdl.PRESSED { //dealer's turn
					mindScore, dScore := s.Table.Players[1].Hand.Score()
					for dScore <= 16 || dScore == 17 && dScore != mindScore {
						s.Table.Players[1].Hit(r, s)
						if err := s.Paint(r); err != nil {
							g.Errc <- err
							return true
						}
						mindScore, dScore = s.Table.Players[1].Hand.Score()
					}
					_, pScore := s.Table.Players[0].Hand.Score()
					fmt.Println("Final Score:")
					fmt.Println("Player", s.Table.Players[0].Hand, "\nScore:", pScore)
					fmt.Println("Dealer", s.Table.Players[1].Hand, "\nScore:", dScore)
					switch {
					case dScore >= pScore && dScore <= 21:
						fmt.Println("You lose")
						if err := gameover(r, s, g, false); err != nil {
							g.Errc <- err
							return true
						}
					case dScore > 21 || pScore > dScore:
						fmt.Println("You win")
						if err := gameover(r, s, g, true); err != nil {
							g.Errc <- err
							return true
						}
					}
					g.Status = 0
				}
			default:
				log.Printf("unknown button %v", e.Button)
			}
		}
	// case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent, *sdl.CommonEvent:
	default:
		//log.Printf("unknown event %T\n%+v\n", event, e)
	}
	return false

}

func gameover(r *sdl.Renderer, s *room.Scene, g *room.Game, win bool) error {
	if err := s.Table.UpdateResult(r, win); err != nil {
		return fmt.Errorf("could not update result: %v", err)
	}
	if err := s.Paint(r); err != nil {
		return err
	}
	g.Status = 0
	return nil
}
