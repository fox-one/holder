package api

import (
	"net/http"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/handler/api/actions"
	"github.com/fox-one/holder/handler/api/system"
	"github.com/fox-one/holder/handler/api/user"
	"github.com/fox-one/holder/handler/auth"
	"github.com/fox-one/holder/handler/render"
	"github.com/fox-one/holder/handler/rpc"
	"github.com/fox-one/holder/pkg/reversetwirp"
	"github.com/go-chi/chi"
	"github.com/twitchtv/twirp"
)

func New(
	sessions core.Session,
	userz core.UserService,
	gems core.GemStore,
	vaults core.VaultStore,
	transactions core.TransactionStore,
	walletz core.WalletService,
	notifier core.Notifier,
	system *core.System,
) *Server {
	return &Server{
		sessions:     sessions,
		userz:        userz,
		gems:         gems,
		vaults:       vaults,
		transactions: transactions,
		walletz:      walletz,
		notifier:     notifier,
		system:       system,
	}
}

type Server struct {
	sessions     core.Session
	userz        core.UserService
	gems         core.GemStore
	vaults       core.VaultStore
	transactions core.TransactionStore
	walletz      core.WalletService
	notifier     core.Notifier
	system       *core.System
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(auth.HandleAuthentication(s.sessions))
	r.Use(render.WrapResponse(true))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Error(w, twirp.NotFoundError("not found"))
	})

	r.Get("/time", system.HandleTime())
	r.Get("/info", system.HandleInfo(s.system))

	r.Post("/login", user.HandleOauth(s.userz, s.sessions, s.notifier))

	svr := rpc.New(s.gems, s.vaults, s.transactions).TwirpServer()
	rt := reversetwirp.NewSingleTwirpServerProxy(svr)

	r.Route("/gems", func(r chi.Router) {
		r.Get("/", rt.Handle("ListGems"))
	})

	r.Route("/vaults", func(r chi.Router) {
		r.Get("/{id}", rt.Handle("FindVault"))
		r.Get("/", rt.Handle("ListVaults"))
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Get("/{id}", rt.Handle("FindTransaction"))
		r.Get("/", rt.Handle("ListTransactions"))
	})

	r.Route("/actions", func(r chi.Router) {
		r.Post("/", actions.HandleCreate(s.walletz, s.system))
	})

	return r
}
