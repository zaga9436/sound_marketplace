import { apiClient } from "@/lib/api/client";
import { AuthResponse } from "@/shared/types/api";

export type LoginPayload = {
  email: string;
  password: string;
};

export type RegisterPayload = LoginPayload & {
  role: "customer" | "engineer";
};

export const authApi = {
  login: (payload: LoginPayload) => apiClient.post<AuthResponse>("/auth/login", payload),
  register: (payload: RegisterPayload) => apiClient.post<AuthResponse>("/auth/register", payload),
  me: () => apiClient.get<AuthResponse>("/auth/me")
};
