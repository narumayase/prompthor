package repository

import (
	"anyompt/config"
	"anyompt/internal/domain"
	"context"
	anysherhttp "github.com/narumayase/anysher/http"
)

// AnywayRepository implements ProducerRepository using anyway API
type AnywayRepository struct {
	apiKey     string
	model      string
	httpClient *anysherhttp.Client
	baseURL    string
}

// NewAnywayRepository creates a new instance of the Anyway repository
func NewAnywayRepository(config config.Config, httpClient *anysherhttp.Client) domain.ProducerRepository {
	return &AnywayRepository{
		model:      config.ChatModel,
		httpClient: httpClient,
		baseURL:    config.GatewayAPIUrl,
	}
}

func (r *AnywayRepository) Send(ctx context.Context, message []byte) error {
	correlationID := ctx.Value("correlation_id").(string)
	routingID := ctx.Value("routing_id").(string)

	// send to anyway
	resp, err := r.httpClient.Post(ctx, anysherhttp.Payload{
		URL:   r.baseURL,
		Token: r.apiKey,
		Headers: map[string]string{
			"Content-Type":     "application/json",
			"X-Correlation-ID": correlationID,
			"X-Routing-ID":     routingID},
		Content: message,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
