package user

import (
	"context"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/mixin-sdk-go"
)

type Config struct {
	ClientSecret string
}

func New(client *mixin.Client, cfg Config) core.UserService {
	return &userService{
		client: client,
		secret: cfg.ClientSecret,
	}
}

type userService struct {
	client *mixin.Client
	secret string
}

func (s *userService) Find(ctx context.Context, mixinID string) (*core.User, error) {
	profile, err := s.client.ReadUser(ctx, mixinID)
	if err != nil {
		return nil, err
	}

	user := &core.User{
		MixinID: profile.UserID,
		Name:    profile.FullName,
		Avatar:  profile.AvatarURL,
	}
	return user, nil
}

func (s *userService) Login(ctx context.Context, token string) (*core.User, error) {
	profile, err := mixin.UserMe(ctx, token)
	if err != nil {
		return nil, err
	}

	user := &core.User{
		MixinID:     profile.UserID,
		Name:        profile.FullName,
		Avatar:      profile.AvatarURL,
		AccessToken: token,
	}
	return user, nil
}

func (s *userService) Auth(ctx context.Context, code string) (string, error) {
	token, _, err := mixin.AuthorizeToken(ctx, s.client.ClientID, s.secret, code, "")
	return token, err
}
