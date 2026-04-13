import { apiClient } from "@/lib/api/client";
import { Card, CardListResponse, Dispute, ModerationAction, User } from "@/shared/types/api";

type UserFilters = {
  role?: string;
  status?: string;
};

type CardFilters = {
  visibility?: string;
  card_type?: string;
  author_id?: string;
  is_published?: string;
  q?: string;
  limit?: string;
  offset?: string;
};

function toQuery(params: Record<string, string | undefined>) {
  const query = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value) query.set(key, value);
  });
  const serialized = query.toString();
  return serialized ? `?${serialized}` : "";
}

export const adminApi = {
  listUsers: (filters: UserFilters = {}) => apiClient.get<User[]>(`/admin/users${toQuery(filters)}`),
  getUser: (id: string) => apiClient.get<User>(`/admin/users/${id}`),
  suspendUser: (id: string, reason: string) => apiClient.post<User>(`/admin/users/${id}/suspend`, { reason }),
  unsuspendUser: (id: string, reason = "") => apiClient.post<User>(`/admin/users/${id}/unsuspend`, { reason }),

  listCards: (filters: CardFilters = {}) => apiClient.get<CardListResponse>(`/admin/cards${toQuery(filters)}`),
  getCard: (id: string) => apiClient.get<Card>(`/admin/cards/${id}`),
  hideCard: (id: string, reason: string) => apiClient.post<Card>(`/admin/cards/${id}/hide`, { reason }),
  unhideCard: (id: string, reason = "") => apiClient.post<Card>(`/admin/cards/${id}/unhide`, { reason }),

  listDisputes: (status?: string) => apiClient.get<Dispute[]>(`/admin/disputes${toQuery({ status })}`),
  getDispute: (id: string) => apiClient.get<Dispute>(`/admin/disputes/${id}`),
  closeDispute: (id: string, resolution: "complete_order" | "cancel_order", reason = "") =>
    apiClient.post<Dispute>(`/admin/disputes/${id}/close`, { resolution, reason }),

  listActions: (params: { target_type?: string; target_id?: string; limit?: string } = {}) =>
    apiClient.get<ModerationAction[]>(`/admin/actions${toQuery(params)}`)
};
