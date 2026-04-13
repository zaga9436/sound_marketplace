import { apiClient } from "@/lib/api/client";
import { Dispute, DisputeResolution } from "@/shared/types/api";

export const disputesApi = {
  getByOrderId: (orderId: string) => apiClient.get<Dispute>(`/orders/${orderId}/dispute`),
  open: (orderId: string, reason: string) => apiClient.post<Dispute>(`/orders/${orderId}/dispute`, { reason }),
  close: (orderId: string, resolution: DisputeResolution, message?: string) =>
    apiClient.post<Dispute>(`/orders/${orderId}/dispute/close`, message ? { resolution, message } : { resolution })
};
