import { apiClient } from "@/lib/api/client";
import { Card, CardListResponse, CardPayload, DownloadResponse, MediaFile } from "@/shared/types/api";

export type CatalogQuery = {
  q?: string;
  card_type?: string;
  kind?: string;
  min_price?: string;
  max_price?: string;
  tag?: string;
  sort_by?: string;
  sort_order?: string;
  limit?: string;
  offset?: string;
  author_id?: string;
};

function buildQueryString(params: CatalogQuery) {
  const search = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value != null && value !== "") {
      search.set(key, value);
    }
  });

  const query = search.toString();
  return query ? `?${query}` : "";
}

export const cardsApi = {
  list: (params: CatalogQuery) => apiClient.get<CardListResponse>(`/cards${buildQueryString(params)}`),
  getById: (id: string) => apiClient.get<Card>(`/cards/${id}`),
  create: (payload: CardPayload) => apiClient.post<Card>("/cards", payload),
  update: (id: string, payload: Omit<CardPayload, "card_type">) => apiClient.put<Card>(`/cards/${id}`, payload),
  uploadPreview: (id: string, formData: FormData) => apiClient.upload<MediaFile>(`/cards/${id}/media/preview`, formData),
  uploadFull: (id: string, formData: FormData) => apiClient.upload<MediaFile>(`/cards/${id}/media/full`, formData),
  getFullDownloadUrl: (id: string) => apiClient.get<DownloadResponse>(`/cards/${id}/download`)
};
