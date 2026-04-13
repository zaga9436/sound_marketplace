"use client";

import { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";

import { authApi } from "@/entities/auth/api/auth";
import { getStoredToken } from "@/lib/api/auth-token";
import { useAuthStore } from "@/lib/auth/session-store";

export function AuthBootstrap() {
  const token = getStoredToken();
  const setBootstrapResult = useAuthStore((state) => state.setBootstrapResult);
  const markHydrated = useAuthStore((state) => state.markHydrated);

  const query = useQuery({
    queryKey: ["auth", "me"],
    queryFn: () => authApi.me(),
    enabled: Boolean(token),
    retry: 0
  });

  useEffect(() => {
    if (!token) {
      setBootstrapResult(null);
      return;
    }

    if (query.isSuccess) {
      setBootstrapResult(query.data);
    } else if (query.isError) {
      setBootstrapResult(null);
    }
  }, [query.data, query.isError, query.isSuccess, setBootstrapResult, token]);

  useEffect(() => {
    if (!token && !query.isFetching) {
      markHydrated();
    }
  }, [markHydrated, query.isFetching, token]);

  return null;
}
