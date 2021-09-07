package main

import (
	"github.com/fox-one/holder/cmd/holder-worker/config"
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/notifier"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/google/wire"
)

var notifierSet = wire.NewSet(
	provideNotifierConfig,
	provideNotifier,
)

func provideNotifierConfig(cfg *config.Config) notifier.Config {
	return notifier.Config{Links: map[string]string{}}
}

func provideNotifier(
	system *core.System,
	assetz core.AssetService,
	messages core.MessageStore,
	pools core.PoolStore,
	vats core.VaultStore,
	users core.UserStore,
	i18n *localizer.Localizer,
	cfg notifier.Config,
) core.Notifier {
	if _flag.notify {
		return notifier.New(
			system,
			assetz,
			messages,
			pools,
			vats,
			users,
			i18n,
			cfg,
		)
	}

	return notifier.Mute()
}
