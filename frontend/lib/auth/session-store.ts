"use client";

import { create } from "zustand";
import { persist } from "zustand/middleware";

import { clearStoredToken, setStoredToken } from "@/lib/api/auth-token";
import { AuthResponse, Profile, Role, User } from "@/shared/types/api";

const ROLE_COOKIE = "sm_role";
const TOKEN_COOKIE = "sm_token";

function writeCookie(name: string, value: string) {
  if (typeof document === "undefined") return;
  document.cookie = `${name}=${value}; path=/; max-age=2592000; samesite=lax`;
}

function clearCookie(name: string) {
  if (typeof document === "undefined") return;
  document.cookie = `${name}=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; samesite=lax`;
}

type AuthStore = {
  token: string | null;
  user: User | null;
  profile: Profile | null;
  hydrated: boolean;
  setSession: (payload: AuthResponse) => void;
  setBootstrapResult: (payload: AuthResponse | null) => void;
  updateProfile: (profile: Profile) => void;
  logout: () => void;
  markHydrated: () => void;
  role: () => Role | null;
};

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      profile: null,
      hydrated: false,
      setSession: (payload) => {
        if (payload.token) {
          setStoredToken(payload.token);
          writeCookie(TOKEN_COOKIE, payload.token);
        }
        writeCookie(ROLE_COOKIE, payload.user.role);
        set({
          token: payload.token ?? get().token,
          user: payload.user,
          profile: payload.profile,
          hydrated: true
        });
      },
      setBootstrapResult: (payload) => {
        if (!payload) {
          set({ token: null, user: null, profile: null, hydrated: true });
          clearStoredToken();
          clearCookie(TOKEN_COOKIE);
          clearCookie(ROLE_COOKIE);
          return;
        }
        writeCookie(ROLE_COOKIE, payload.user.role);
        set({
          token: get().token,
          user: payload.user,
          profile: payload.profile,
          hydrated: true
        });
      },
      updateProfile: (profile) => set({ profile }),
      logout: () => {
        clearStoredToken();
        clearCookie(TOKEN_COOKIE);
        clearCookie(ROLE_COOKIE);
        set({ token: null, user: null, profile: null, hydrated: true });
      },
      markHydrated: () => set({ hydrated: true }),
      role: () => get().user?.role ?? null
    }),
    {
      name: "soundmarket-auth",
      partialize: (state) => ({
        token: state.token,
        user: state.user,
        profile: state.profile
      })
    }
  )
);
