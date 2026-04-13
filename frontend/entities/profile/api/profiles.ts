import { apiClient } from "@/lib/api/client";
import { MediaFile, Profile, ProfileUpdatePayload, Review } from "@/shared/types/api";

export const profilesApi = {
  me: () => apiClient.get<Profile>("/profiles/me"),
  updateMe: (payload: ProfileUpdatePayload) => apiClient.put<Profile>("/profiles/me", payload),
  uploadAvatar: (formData: FormData) => apiClient.upload<MediaFile>("/profiles/me/avatar", formData),
  getById: (id: string) => apiClient.get<Profile>(`/profiles/${id}`),
  listReviews: (id: string) => apiClient.get<Review[]>(`/profiles/${id}/reviews`)
};
