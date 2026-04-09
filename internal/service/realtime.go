package service

import (
	"context"
	"time"

	"github.com/soundmarket/backend/internal/storage"
	"github.com/soundmarket/backend/internal/worker"
)

type RealtimeService struct {
	queue   worker.Queue
	storage storage.Adapter
}

func NewRealtimeService(queue worker.Queue, storage storage.Adapter) *RealtimeService {
	return &RealtimeService{queue: queue, storage: storage}
}

func (s *RealtimeService) ChatInfo(orderID string) map[string]string {
	deliveryURL, _ := s.storage.GenerateSignedURL(context.Background(), "orders/"+orderID+"/preview.m3u8", 15*time.Minute)
	return map[string]string{
		"order_id":     orderID,
		"channel":      "order:" + orderID,
		"delivery_url": deliveryURL,
	}
}
