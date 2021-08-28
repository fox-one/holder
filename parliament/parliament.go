package parliament

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/mtg"
	"github.com/fox-one/holder/pkg/number"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/holder/service/asset"
	"github.com/fox-one/mixin-sdk-go"
)

type Config struct {
	Links map[string]string
}

func New(
	messages core.MessageStore,
	userz core.UserService,
	assetz core.AssetService,
	walletz core.WalletService,
	system *core.System,
	cfg Config,
) core.Parliament {
	links := &template.Template{}
	for name, tpl := range cfg.Links {
		if tpl == "" {
			continue
		}

		links = template.Must(
			links.New(name).Parse(tpl),
		)
	}

	return &parliament{
		messages: messages,
		userz:    userz,
		assetz:   asset.Cache(assetz),
		walletz:  walletz,
		system:   system,
		links:    links,
	}
}

type parliament struct {
	messages core.MessageStore
	userz    core.UserService
	assetz   core.AssetService
	walletz  core.WalletService
	system   *core.System
	links    *template.Template
}

func (s *parliament) executeLink(name string, data interface{}) (string, error) {
	b := bytes.Buffer{}
	if err := s.links.ExecuteTemplate(&b, name, data); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (s *parliament) requestVoteAction(ctx context.Context, proposal *core.Proposal) (string, error) {
	id, _ := uuid.FromString(proposal.TraceID)
	body, err := mtg.Encode(core.ActionProposalVote, id)
	if err != nil {
		return "", err
	}

	data, err := core.TransactionAction{
		Body: body,
	}.Encode()
	if err != nil {
		return "", err
	}

	memo, err := mtg.Encrypt(data, mixin.GenerateEd25519Key(), s.system.PublicKey)
	if err != nil {
		return "", err
	}

	transfer := &core.Transfer{
		TraceID:   uuid.Modify(proposal.TraceID, s.system.ClientID),
		AssetID:   s.system.GasAssetID,
		Amount:    s.system.GasAmount,
		Memo:      base64.StdEncoding.EncodeToString(memo),
		Threshold: s.system.Threshold,
		Opponents: s.system.Members,
	}

	code, err := s.walletz.ReqTransfer(ctx, transfer)
	if err != nil {
		return "", err
	}

	return paymentAction(code), nil
}

func (s *parliament) ProposalCreated(ctx context.Context, p *core.Proposal) error {
	view := Proposal{
		Number: p.ID,
		Action: p.Action.String(),
		Info: []Item{
			{
				Key:   "action",
				Value: p.Action.String(),
			},
			{
				Key:   "id",
				Value: p.TraceID,
			},
			{
				Key:   "date",
				Value: p.CreatedAt.Format(time.RFC3339),
			},
			{
				Key:    "creator",
				Value:  s.fetchUserName(ctx, p.Creator),
				Action: userAction(p.Creator),
			},
			{
				Key:    "pay",
				Value:  fmt.Sprintf("%s %s", number.Humanize(p.Amount), s.fetchAssetSymbol(ctx, p.AssetID)),
				Action: assetAction(p.AssetID),
			},
		},
	}

	data, _ := base64.StdEncoding.DecodeString(p.Data)
	view.Meta = s.renderProposalItems(ctx, p.Action, data)

	items := append(view.Info, view.Meta...)
	voteAction, err := s.requestVoteAction(ctx, p)
	if err != nil {
		return err
	}

	items = append(items, Item{
		Key:    "Vote",
		Value:  "Vote",
		Action: voteAction,
	})

	buttons := generateButtons(items)
	buttonsData, _ := json.Marshal(buttons)
	post := execute("proposal_created", view)

	var messages []*core.Message
	for _, admin := range s.system.Admins {
		// post
		postMsg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(p.TraceID, s.system.ClientID+admin),
			Category:       mixin.MessageCategoryPlainPost,
			Data:           base64.StdEncoding.EncodeToString(post),
		}

		// buttons
		buttonMsg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(postMsg.MessageID, "buttons"),
			Category:       mixin.MessageCategoryAppButtonGroup,
			Data:           base64.StdEncoding.EncodeToString(buttonsData),
		}

		messages = append(messages, core.BuildMessage(postMsg), core.BuildMessage(buttonMsg))
	}

	return s.messages.Create(ctx, messages)
}

func (s *parliament) ProposalApproved(ctx context.Context, p *core.Proposal) error {
	by := p.Votes[len(p.Votes)-1]

	view := Proposal{
		ApprovedCount: len(p.Votes),
		ApprovedBy:    s.fetchUserName(ctx, by),
	}

	post := execute("proposal_approved", view)

	var messages []*core.Message
	for _, admin := range s.system.Admins {
		quote := uuid.Modify(p.TraceID, s.system.ClientID+admin)
		msg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(quote, "ProposalApproved By "+by),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString(post),
			QuoteMessageID: quote,
		}

		messages = append(messages, core.BuildMessage(msg))
	}

	return s.messages.Create(ctx, messages)
}

func (s *parliament) ProposalPassed(ctx context.Context, proposal *core.Proposal) error {
	var messages []*core.Message

	post := execute("proposal_passed", nil)

	for _, admin := range s.system.Admins {
		quote := uuid.Modify(proposal.TraceID, s.system.ClientID+admin)
		msg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(quote, "ProposalPassed"),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString(post),
			QuoteMessageID: quote,
		}

		messages = append(messages, core.BuildMessage(msg))
	}

	return s.messages.Create(ctx, messages)
}
