package di

import (
	"test-backend-1-kuprinvv/internal/client"
	"test-backend-1-kuprinvv/internal/client/conference"
)

type clientProvider struct {
	conference client.ConferenceClient
}

func (c *Container) ConferenceClient() client.ConferenceClient {
	if c.clients.conference == nil {
		c.clients.conference = conference.NewMock()
	}
	return c.clients.conference
}
