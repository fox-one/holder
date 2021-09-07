package rpc

import (
	"context"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/handler/auth"
	"github.com/fox-one/holder/handler/request"
	"github.com/fox-one/holder/handler/rpc/api"
	"github.com/fox-one/holder/handler/rpc/views"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store"
	"github.com/spf13/cast"
	"github.com/twitchtv/twirp"
)

func New(
	pools core.PoolStore,
	vaults core.VaultStore,
	transactions core.TransactionStore,
) *Server {
	return &Server{
		pools:        pools,
		vaults:       vaults,
		transactions: transactions,
	}
}

type Server struct {
	pools        core.PoolStore
	vaults       core.VaultStore
	transactions core.TransactionStore
}

func (s *Server) TwirpServer() api.TwirpServer {
	opts := []interface{}{
		twirp.WithServerJSONSkipDefaults(false),
		twirp.WithServerInterceptors(func(next twirp.Method) twirp.Method {
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				if _, err := govalidator.ValidateStruct(req); err != nil {
					return nil, twirp.InvalidArgumentError("", err.Error())
				}

				return next(ctx, req)
			}
		}),
	}

	return api.NewHolderServer(s, opts...)
}

func (s *Server) Handle(sessions core.Session) http.Handler {
	return auth.HandleAuthentication(sessions)(s.TwirpServer())
}

// FindTransaction godoc
// @Summary find tx by follow id
// @Description
// @Tags Transactions
// @Accept json
// @Produce json
// @param Authorization header string true "Example: Bearer foo"
// @param follow_id path string true "follow id"
// @Success 200 {object} api.Transaction
// @Router /transactions/{follow_id} [get]
func (s *Server) FindTransaction(ctx context.Context, req *api.Req_FindTransaction) (*api.Transaction, error) {
	user, ok := request.UserFrom(ctx)
	if !ok {
		logger.FromContext(ctx).Debugln("rpc: authentication required")
		return nil, twirp.NewError(twirp.Unauthenticated, "authentication required")
	}

	tx, err := s.transactions.FindFollow(ctx, user.MixinID, req.Id)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil, twirp.NotFoundError("transaction not found")
		}

		return nil, err
	}

	return views.Transaction(tx), nil
}

// ListTransactions godoc
// @Summary list transactions
// @Description
// @Tags Transactions
// @Accept json
// @Produce json
// @param request query api.Req_ListTransactions false "default limit 50"
// @Success 200 {object} api.Resp_ListTransactions
// @Router /transactions [get]
func (s *Server) ListTransactions(ctx context.Context, req *api.Req_ListTransactions) (*api.Resp_ListTransactions, error) {
	fromID := cast.ToInt64(req.Cursor)
	limit := 50
	if l := int(req.Limit); l > 0 && l < limit {
		limit = l
	}

	transactions, err := s.transactions.List(ctx, fromID, limit+1)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: transactions.ListTarget")
		return nil, err
	}

	resp := &api.Resp_ListTransactions{
		Pagination: &api.Pagination{},
	}

	for idx, t := range transactions {
		resp.Transactions = append(resp.Transactions, views.Transaction(t))

		if idx == limit-1 {
			resp.Pagination.NextCursor = cast.ToString(t.ID)
			resp.Pagination.HasNext = true
			break
		}
	}

	return resp, nil
}

// ListPools godoc
// @Summary list all pools
// @Description
// @Tags Pools
// @Accept json
// @Produce json
// @Success 200 {object} api.Resp_ListPools
// @Router /pools [get]
func (s *Server) ListPools(ctx context.Context, _ *api.Req_ListPools) (*api.Resp_ListPools, error) {
	pools, err := s.pools.List(ctx)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("rpc: pools.List")
		return nil, err
	}

	resp := &api.Resp_ListPools{Pools: make([]*api.Pool, 0, len(pools))}
	for _, pool := range pools {
		resp.Pools = append(resp.Pools, views.Pool(pool))
	}

	return resp, nil
}

// FindVault godoc
// @Summary Find Vault By ID
// @Description
// @Tags Vaults
// @Accept json
// @Produce json
// @param id path string true "vault id"
// @Success 200 {object} api.Vault
// @Router /vaults/{id} [get]
func (s *Server) FindVault(ctx context.Context, req *api.Req_FindVault) (*api.Vault, error) {
	vault, err := s.vaults.Find(ctx, req.Id)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("rpc: vaults.Find")
		return nil, err
	}

	pool, err := s.pools.Find(ctx, vault.AssetID)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("rpc: pools.Find")
		return nil, err
	}

	return views.Vault(vault, pool), nil
}

// ListVaults godoc
// @Summary List Vaults
// @Description
// @Tags Vaults
// @Accept json
// @Produce json
// @Success 200 {object} api.Resp_ListVaults
// @Router /vaults [get]
func (s *Server) ListVaults(ctx context.Context, _ *api.Req_ListVaults) (*api.Resp_ListVaults, error) {
	user, ok := request.UserFrom(ctx)
	if !ok {
		logger.FromContext(ctx).Debugln("rpc: authentication required")
		return nil, twirp.NewError(twirp.Unauthenticated, "authentication required")
	}

	vaults, err := s.vaults.ListUser(ctx, user.MixinID)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("vaults.ListUser")
		return nil, err
	}

	resp := &api.Resp_ListVaults{
		Vaults: make([]*api.Vault, 0, len(vaults)),
	}

	if len(vaults) == 0 {
		return resp, nil
	}

	pools, err := s.pools.List(ctx)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("pools.List")
		return nil, err
	}

	poolMap := toPoolMap(pools)
	for _, vault := range vaults {
		resp.Vaults = append(resp.Vaults, views.Vault(vault, poolMap[vault.AssetID]))
	}

	return resp, nil
}

func toPoolMap(pools []*core.Pool) map[string]*core.Pool {
	m := make(map[string]*core.Pool, len(pools))
	for _, pool := range pools {
		m[pool.ID] = pool
	}

	return m
}
