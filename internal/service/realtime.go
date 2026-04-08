package service

import (
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
	return map[string]string{
		"order_id":     orderID,
		"channel":      "order:" + orderID,
		"delivery_url": s.storage.SignedURL("orders/" + orderID + "/preview.m3u8"),
	}
}
