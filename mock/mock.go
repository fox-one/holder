package mock

//go:generate mockgen -package=mock -destination=mock_gen.go github.com/fox-one/holder/core AssetService,UserStore,UserService,MessageStore,MessageService,Notifier,PoolStore,VaultStore,WalletStore
