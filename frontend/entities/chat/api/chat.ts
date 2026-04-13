import { apiClient } from "@/lib/api/client";
import { ChatMessage, Conversation } from "@/shared/types/api";

export const chatApi = {
  listConversations: (limit = 20) => apiClient.get<Conversation[]>(`/chats?limit=${limit}`),
  listMessages: (orderId: string, limit = 50, beforeId?: string) => {
    const search = new URLSearchParams({ limit: String(limit) });
    if (beforeId) search.set("before_id", beforeId);
    return apiClient.get<ChatMessage[]>(`/orders/${orderId}/messages?${search.toString()}`);
  },
  sendMessage: (orderId: string, body: string) => apiClient.post<ChatMessage>(`/orders/${orderId}/messages`, { body }),
  markRead: (orderId: string) => apiClient.post<{ ok: boolean }>(`/orders/${orderId}/messages/read`, {})
};
