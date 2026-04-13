import { apiClient } from "@/lib/api/client";
import { Profile, ProfileUpdatePayload, Review } from "@/shared/types/api";

export const profilesApi = {
  me: () => apiClient.get<Profile>("/profiles/me"),
  updateMe: (payload: ProfileUpdatePayload) => apiClient.put<Profile>("/profiles/me", payload),
  getById: (id: string) => apiClient.get<Profile>(`/profiles/${id}`),
  listReviews: (id: string) => apiClient.get<Review[]>(`/profiles/${id}/reviews`)
};
