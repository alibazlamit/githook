package repository_test

import (
	"context"
	"testing"

	"github.com/ali/githook/internal/repository"
)

func TestPostgresWebhookRepositoryDeliveryExistsStub(t *testing.T) {
	tests := []struct {
		name       string
		deliveryID string
	}{
		{"empty id", ""},
		{"non-empty id", "abc-123"},
	}

	repo := repository.NewPostgresWebhookRepository()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.DeliveryExists(context.Background(), tt.deliveryID)
			if err != nil {
				t.Errorf("err: want nil, got %v", err)
			}
			if got != false {
				t.Errorf("dup: want false, got %v", got)
			}
		})
	}
}
