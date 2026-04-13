import { apiClient } from "@/lib/api/client";
import { NotificationListResponse } from "@/shared/types/api";

export const notificationsApi = {
  list: (limit = 20, beforeId?: string) => {
    const search = new URLSearchParams({ limit: String(limit) });
    if (beforeId) search.set("before_id", beforeId);
    return apiClient.get<NotificationListResponse>(`/notifications?${search.toString()}`);
  },
  markRead: (ids?: string[]) => apiClient.post<NotificationListResponse>("/notifications/read", ids?.length ? { ids } : {})
};
