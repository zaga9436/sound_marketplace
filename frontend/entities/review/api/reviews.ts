import { apiClient } from "@/lib/api/client";
import { Review } from "@/shared/types/api";

export const reviewsApi = {
  create: (orderId: string, rating: number, text: string) => apiClient.post<Review>(`/orders/${orderId}/reviews`, { rating, text })
};
