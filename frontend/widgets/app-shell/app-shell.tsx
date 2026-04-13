"use client";

import Link from "next/link";
import { PropsWithChildren } from "react";

import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { RoleNavigation } from "@/widgets/app-shell/navigation";

export function AppShell({ children }: PropsWithChildren) {
  const user = useAuthStore((state) => state.user);
  const profile = useAuthStore((state) => state.profile);

  return (
    <div className="page-frame">
      <div className="container grid min-h-screen gap-6 py-6 lg:grid-cols-[280px_minmax(0,1fr)]">
        <div className="hidden lg:block">
          <RoleNavigation />
        </div>

        <div className="flex min-w-0 flex-col gap-6">
          <header className="surface flex items-center justify-between px-6 py-4">
            <div className="space-y-1">
              <Badge variant="secondary">{user?.role ?? "guest"}</Badge>
              <p className="text-sm text-muted-foreground">{profile?.display_name ?? user?.email ?? "Рабочее пространство SoundMarket"}</p>
            </div>
            <Link href="/" className="text-sm font-medium text-primary">
              На главную
            </Link>
          </header>

          <main className="min-w-0">{children}</main>
        </div>
      </div>
    </div>
  );
}
