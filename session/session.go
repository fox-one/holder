package session

import (
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/bluele/gcache"
	"github.com/dgrijalva/jwt-go"
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/internal/request"
	"github.com/fox-one/mixin-sdk-go"
	"golang.org/x/sync/singleflight"
)

type Config struct {
	Capacity int
	Issuers  []string
}

func New(users core.UserStore, userz core.UserService, cfg Config) core.Session {
	var s core.Session = &session{
		users:   users,
		userz:   userz,
		issuers: cfg.Issuers,
		sf:      &singleflight.Group{},
	}

	if cfg.Capacity > 0 {
		s = &cacheSession{
			Session: s,
			tokens:  gcache.New(cfg.Capacity).LRU().Build(),
		}
	}

	return s
}

type session struct {
	userz   core.UserService
	users   core.UserStore
	sf      *singleflight.Group
	issuers []string
}

func (s *session) Login(r *http.Request) (*core.User, error) {
	accessToken := request.ExtractBearerToken(r)
	if accessToken == "" {
		return nil, errors.New("token not found")
	}

	ctx := r.Context()

	user, err, _ := s.sf.Do(accessToken, func() (interface{}, error) {
		var claim struct {
			jwt.StandardClaims
			Scope string `json:"scp,omitempty"`
		}
		_, _ = jwt.ParseWithClaims(accessToken, &claim, nil)

		if claim.Scope != "FULL" && !govalidator.IsIn(claim.Issuer, s.issuers...) {
			return nil, errors.New("invalid issuer")
		}

		if jti := claim.Id; govalidator.IsUUID(jti) {
			ctx = mixin.WithRequestID(ctx, jti)
		}

		user, err := s.userz.Login(ctx, accessToken)
		if err != nil {
			return nil, err
		}

		// handle language
		user.Lang = request.ExtractPreferLanguage(r)
		if err := s.users.Save(ctx, user); err != nil {
			return nil, err
		}

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	return user.(*core.User), nil
}
