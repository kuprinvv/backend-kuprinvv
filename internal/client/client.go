package client

import "context"

type ConferenceClient interface {
	CreateLink(ctx context.Context) (string, error)
}
