package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/toledoom/gork/internal/app"
	httpport "github.com/toledoom/gork/internal/ports/http"
	"github.com/toledoom/gork/pkg/application"
)

func main() {
	a := application.New(app.SetupCommandHandlers, app.SetupQueryHandlers)
	a.Start(app.SetupServices, app.SetupDataMapper, app.SetupEventPublisher)

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

	a.HttpListenAndServe(":8080", r)
}
