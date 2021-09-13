package notifier

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"text/template"
	"time"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/notifier/compose"
	"github.com/fox-one/holder/pkg/cont/vat"
	"github.com/fox-one/holder/pkg/mtg"
	"github.com/fox-one/holder/pkg/mtg/types"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/holder/service/asset"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/spf13/cast"
)

type Config struct {
	Links map[string]string
}

func New(
	system *core.System,
	assetz core.AssetService,
	messages core.MessageStore,
	pools core.PoolStore,
	vats core.VaultStore,
	users core.UserStore,
	walletz core.WalletService,
	i18n *localizer.Localizer,
	cfg Config,
) core.Notifier {
	links := &template.Template{}
	for name, tpl := range cfg.Links {
		if tpl == "" {
			continue
		}

		links = template.Must(
			links.New(name).Parse(tpl),
		)
	}

	return &notifier{
		system:   system,
		assetz:   asset.Cache(assetz),
		messages: messages,
		pools:    pools,
		vats:     vats,
		users:    users,
		walletz:  walletz,
		i18n:     i18n,
		links:    links,
	}
}

type notifier struct {
	system   *core.System
	assetz   core.AssetService
	messages core.MessageStore
	pools    core.PoolStore
	vats     core.VaultStore
	users    core.UserStore
	walletz  core.WalletService
	i18n     *localizer.Localizer
	links    *template.Template
}

func (n *notifier) executeLink(name string, data interface{}) (string, error) {
	b := bytes.Buffer{}
	if err := n.links.ExecuteTemplate(&b, name, data); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (n *notifier) buildPaymentURL(ctx context.Context, traceID string, args ...interface{}) (string, error) {
	data, err := mtg.Encode(args...)
	if err != nil {
		return "", err
	}

	data, _ = core.TransactionAction{Body: data}.Encode()
	memo := base64.StdEncoding.EncodeToString(data)

	transfer := &core.Transfer{
		TraceID:   traceID,
		AssetID:   n.system.GasAssetID,
		Amount:    n.system.GasAmount,
		Memo:      memo,
		Threshold: n.system.Threshold,
		Opponents: n.system.Members,
	}

	code, err := n.walletz.ReqTransfer(ctx, transfer)
	if err != nil {
		return "", err
	}

	return mixin.URL.Codes(code), nil
}

func (n *notifier) localize(id, lang string, args ...interface{}) string {
	l := n.i18n
	if lang != "" {
		l = localizer.WithLanguage(l, lang)
	}

	s := l.LocalizeOr(id, id, args...)
	return strings.TrimSpace(s)
}

func (n *notifier) fetchLanguage(ctx context.Context, userID string) string {
	user, err := n.users.Find(ctx, userID)
	if err != nil {
		return ""
	}

	return user.Lang
}

func (n *notifier) fetchAssetSymbol(ctx context.Context, id string) string {
	a, err := n.assetz.Find(ctx, id)
	if err != nil {
		return ""
	}

	return a.Symbol
}

func (n *notifier) Auth(ctx context.Context, user *core.User) error {
	msg := n.localize("login_done", user.Lang)
	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(n.system.ClientID, user.MixinID),
		RecipientID:    user.MixinID,
		MessageID:      uuid.Modify(user.MixinID, user.AccessToken),
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(msg)),
	}

	return n.messages.Create(ctx, []*core.Message{core.BuildMessage(req)})
}

func (n *notifier) Snapshot(ctx context.Context, transfer *core.Transfer, signedTx string) error {
	log := logger.FromContext(ctx)

	tx, err := mixin.TransactionFromRaw(signedTx)
	if err != nil {
		log.WithError(err).Debugln("decode transaction from raw tx failed")
		return nil
	}

	hash, err := tx.TransactionHash()
	if err != nil {
		return nil
	}

	traceID := mixinRawTransactionTraceId(hash.String(), 0)

	if len(transfer.Opponents) != 1 {
		log.Debugln("transfer opponents is not 1")
		return nil
	}

	coin, err := n.assetz.Find(ctx, transfer.AssetID)
	if err != nil {
		return err
	}

	card := mixin.AppCardMessage{
		AppID:       n.system.ClientID,
		IconURL:     coin.Logo,
		Title:       transfer.Amount.String(),
		Description: coin.Symbol,
		Action:      mixin.URL.Snapshots("", traceID),
	}
	data, _ := json.Marshal(card)

	recipientID := transfer.Opponents[0]
	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(n.system.ClientID, recipientID),
		RecipientID:    recipientID,
		MessageID:      traceID,
		Category:       mixin.MessageCategoryAppCard,
		Data:           base64.StdEncoding.EncodeToString(data),
	}

	return n.messages.Create(ctx, []*core.Message{core.BuildMessage(req)})
}

func (n *notifier) LockDone(ctx context.Context, pool *core.Pool, vault *core.Vault) error {
	b := compose.New(n.system, n.i18n, n.users)

	args := map[string]interface{}{
		"TraceID":   vault.TraceID,
		"Amount":    vault.Amount.String(),
		"Symbol":    n.fetchAssetSymbol(ctx, vault.AssetID),
		"CreatedAt": vault.CreatedAt.Format(time.RFC3339),
		"ExpiredAt": vault.EndAt().Format(time.RFC3339),
	}

	if reward := vat.GetReward(pool, vault); reward.IsPositive() {
		args["Reward"] = reward.String()
	}

	traceID := uuid.Modify(vault.TraceID, "lock_done_notify")
	w := b.Write(traceID, vault.UserID)
	w.Text(w.Localize("lock_done", args))

	if action, err := n.buildPaymentURL(ctx, traceID, core.ActionVaultRelease, types.UUID(vault.TraceID)); err == nil {
		w.Button(w.Localize("withdraw_button"), action)
	}

	messages := b.Messages()
	return n.messages.Create(ctx, messages)
}

func (n *notifier) Transaction(ctx context.Context, tx *core.Transaction) error {
	b := compose.New(n.system, n.i18n, n.users)

	if tx.UserID != "" {
		w := b.Write(tx.TraceID, tx.UserID)
		args := map[string]interface{}{
			"Action":     w.Localize("Action" + tx.Action.String()),
			"Message":    w.Localize(tx.Message),
			"FollowID":   tx.FollowID,
			"Parameters": string(tx.Parameters),
		}

		if tx.Status == core.TransactionStatusOk {
			w.Text(w.Localize("tx_ok", args))
		} else {
			w.Text(w.Localize("tx_abort", args))
		}
	}

	if tx.Status == core.TransactionStatusOk {
		var (
			parameters []interface{}
			err        error
		)

		_ = tx.Parameters.Unmarshal(&parameters)

		switch tx.Action {
		case core.ActionVaultLock:
			err = n.handleVaultLocked(ctx, b, tx, parameters)
		case core.ActionVaultRelease:
			err = n.handleVaultReleased(ctx, b, tx, parameters)
		case core.ActionPoolDonate:
			err = n.handlePoolDonated(ctx, b, tx, parameters)
		}

		if err != nil {
			return err
		}
	}

	messages := b.Messages()
	return n.messages.Create(ctx, messages)
}

func (n *notifier) handleVaultLocked(ctx context.Context, b *compose.Outbox, tx *core.Transaction, _ []interface{}) error {
	vaultID := tx.TraceID
	vault, err := n.vats.Find(ctx, vaultID)
	if err != nil {
		return err
	}

	if vault.ID == 0 {
		return nil
	}

	args := map[string]interface{}{
		"TraceID":   vault.TraceID,
		"Amount":    vault.Amount.String(),
		"Symbol":    n.fetchAssetSymbol(ctx, vault.AssetID),
		"CreatedAt": vault.CreatedAt.Format(time.RFC3339),
		"ExpiredAt": vault.EndAt().Format(time.RFC3339),
	}

	w := b.Write(tx.TraceID, vault.UserID)
	w.Text(w.Localize("vault_locked", args))

	return nil
}

func (n *notifier) handleVaultReleased(ctx context.Context, b *compose.Outbox, tx *core.Transaction, parameters []interface{}) error {
	vaultID := cast.ToString(parameters[0])
	vault, err := n.vats.Find(ctx, vaultID)
	if err != nil {
		return err
	}

	if vault.ID == 0 {
		return nil
	}

	args := map[string]interface{}{
		"TraceID":    vault.TraceID,
		"Amount":     vault.Amount.String(),
		"FillAmount": vault.Amount.Add(vault.Reward).Sub(vault.Penalty).String(),
		"Symbol":     n.fetchAssetSymbol(ctx, vault.AssetID),
		"CreatedAt":  vault.CreatedAt.Format(time.RFC3339),
		"ExpiredAt":  vault.EndAt().Format(time.RFC3339),
		"ReleasedAt": vault.ReleasedAt.Format(time.RFC3339),
	}

	if vault.Reward.IsPositive() {
		args["Reward"] = vault.Reward.String()
	}

	if vault.Penalty.IsPositive() {
		args["Penalty"] = vault.Penalty.String()
	}

	w := b.Write(tx.TraceID, vault.UserID)
	w.Text(w.Localize("vault_released", args))

	return nil
}

func (n *notifier) handlePoolDonated(ctx context.Context, b *compose.Outbox, tx *core.Transaction, _ []interface{}) error {
	args := map[string]interface{}{
		"Amount":   tx.Amount.String(),
		"Symbol":   n.fetchAssetSymbol(ctx, tx.AssetID),
		"DonateAt": tx.CreatedAt.Format(time.RFC3339),
	}

	w := b.Write(tx.TraceID, tx.UserID)
	w.Text(w.Localize("donated", args))

	return nil
}
