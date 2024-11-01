package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/toledoom/gork/internal/app"
	httpport "github.com/toledoom/gork/internal/ports/http"
	"github.com/toledoom/gork/pkg/gork"
)

func main() {
	a := gork.NewApp(app.SetupUseCases, app.SetupCommandHandlers, app.SetupQueryHandlers)
	a.Start(app.SetupServices)

	httpApi := httpport.NewApi(a)

	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/battle", httpApi.StartBattleHandler)
	r.Put("/battle", httpApi.FinishBattleHandler)
	r.Post("/player", httpApi.CreatePlayerHandler)
	r.Get("/player/{playerID}", httpApi.GetPlayerByIDHandler)
	r.Get("/rank", httpApi.GetRankHandler)
	r.Get("/rank/top_players", httpApi.GetTopPlayersHandler)
	/////////////////////

	err := http.ListenAndServe(":8080", r)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
