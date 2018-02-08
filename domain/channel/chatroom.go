package channel

import (
	"context"

	"cloud.google.com/go/firestore"
	"gitlab.com/shinofara/alpha/domain/internal"
	"gitlab.com/shinofara/alpha/domain/message"
	"gitlab.com/shinofara/alpha/domain/type"
	"gitlab.com/shinofara/alpha/domain/user"
)

const Collection = "channel"

type Channel struct {
	ID      _type.ChannelID `firestore:"-"`
	Name    string
	OwnerID _type.UserID

	Owner    *user.User         `firestore:"-"`
	Messages []*message.Message `firestore:"-"`
	Members  []*user.User       `firestore:"-"`
}

type Repository struct {
	ctx context.Context
	cli *firestore.Client
}

func New(cli *firestore.Client, ctx context.Context) *Repository {
	return &Repository{
		cli: cli,
		ctx: ctx,
	}
}

func (r *Repository) Find(id _type.ChannelID) (*Channel, error) {
	ref, err := r.cli.Collection(Collection).Doc(string(id)).Get(r.ctx)
	if err != nil {
		return nil, err
	}

	c := new(Channel)
	if err := internal.Convert(ref, &c); err != nil {
		return nil, err
	}

	return c, nil
}

func (r *Repository) Add(c *Channel) (*Channel, error) {
	ref, _, err := r.cli.Collection(Collection).Add(r.ctx, c)
	if err != nil {
		return nil, err
	}
	cc := *c
	cc.ID = _type.ChannelID(ref.ID)
	return &cc, nil
}
