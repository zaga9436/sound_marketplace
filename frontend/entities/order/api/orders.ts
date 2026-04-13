import { apiClient } from "@/lib/api/client";
import { Order, OrderStatus } from "@/shared/types/api";

export const ordersApi = {
  list: () => apiClient.get<Order[]>("/orders"),
  getById: (id: string) => apiClient.get<Order>(`/orders/${id}`),
  createFromBid: (bidId: string) => apiClient.post<Order>("/orders/from-bid", { bid_id: bidId }),
  createFromOffer: (cardId: string) => apiClient.post<Order>("/orders/from-offer", { card_id: cardId }),
  updateStatus: (id: string, status: OrderStatus) => apiClient.patch<Order>(`/orders/${id}/status`, { status })
};
