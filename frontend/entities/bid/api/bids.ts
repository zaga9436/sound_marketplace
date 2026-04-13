import { apiClient } from "@/lib/api/client";
import { Bid, BidCreatePayload } from "@/shared/types/api";

export const bidsApi = {
  listByRequest: (requestId: string) => apiClient.get<Bid[]>(`/requests/${requestId}/bids`),
  create: (requestId: string, payload: BidCreatePayload) => apiClient.post<Bid>(`/requests/${requestId}/bids`, payload)
};
