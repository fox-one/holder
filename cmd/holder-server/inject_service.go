package main

import (
	"github.com/fox-one/holder/cmd/holder-server/config"
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/service/asset"
	"github.com/fox-one/holder/service/message"
	"github.com/fox-one/holder/service/user"
	"github.com/fox-one/holder/service/wallet"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/google/wire"
	"golang.org/x/text/language"
)

var serviceSet = wire.NewSet(
	provideMixinClient,
	asset.New,
	message.New,
	wire.Value(user.Config{}),
	user.New,
	provideSystem,
	provideWalletServiceConfig,
	wallet.New,
	provideLocalizer,
)

func provideMixinClient(cfg *config.Config) (*mixin.Client, error) {
	return mixin.NewFromKeystore(&cfg.Dapp.Keystore)
}

func provideWalletServiceConfig(cfg *config.Config) wallet.Config {
	return wallet.Config{
		Pin:       cfg.Dapp.Pin,
		Members:   cfg.Group.Members,
		Threshold: cfg.Group.Threshold,
	}
}

func provideSystem(cfg *config.Config) *core.System {
	return &core.System{
		Admins:       cfg.Group.Admins,
		ClientID:     cfg.Dapp.ClientID,
		ClientSecret: cfg.Dapp.ClientSecret,
		Members:      cfg.Group.Members,
		Threshold:    cfg.Group.Threshold,
		Version:      version,
	}
}

func provideLocalizer(cfg *config.Config) (*localizer.Localizer, error) {
	files, err := localizer.FindMessageFiles(cfg.I18n.Path)
	if err != nil {
		return nil, err
	}

	lang, err := language.Parse(cfg.I18n.Language)
	if err != nil {
		return nil, err
	}

	return localizer.New(lang, files...), nil
}
