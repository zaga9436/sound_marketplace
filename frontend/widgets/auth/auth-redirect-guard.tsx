"use client";

import { PropsWithChildren, useEffect } from "react";
import { useRouter } from "next/navigation";

import { useAuthStore } from "@/lib/auth/session-store";

export function AuthRedirectGuard({ children }: PropsWithChildren) {
  const router = useRouter();
  const hydrated = useAuthStore((state) => state.hydrated);
  const user = useAuthStore((state) => state.user);

  useEffect(() => {
    if (!hydrated || !user) return;
    router.replace(user.role === "admin" ? "/admin" : "/dashboard");
  }, [hydrated, router, user]);

  if (hydrated && user) {
    return (
      <div className="surface mx-auto flex w-full max-w-md flex-col gap-3 p-8 text-center">
        <h2 className="text-2xl font-semibold text-slate-950">Вы уже авторизованы</h2>
        <p className="text-slate-600">Переводим вас в рабочую область SoundMarket.</p>
      </div>
    );
  }

  return <>{children}</>;
}
