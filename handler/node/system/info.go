package system

import (
	"net/http"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/handler/render"
)

func HandleInfo(system *core.System) http.HandlerFunc {
	view := render.H{
		"version":   system.Version,
		"members":   system.Members,
		"threshold": system.Threshold,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, view)
	}
}
