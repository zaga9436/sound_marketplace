"use client";

import Link from "next/link";
import { PropsWithChildren } from "react";

import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { UserAvatar } from "@/shared/ui/user-avatar";
import { RoleNavigation } from "@/widgets/app-shell/navigation";

function roleLabel(role?: string | null) {
  if (role === "customer") return "Заказчик";
  if (role === "engineer") return "Инженер";
  if (role === "admin") return "Администратор";
  return "Гость";
}

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
          <header className="flex items-center justify-between rounded-[2rem] border border-slate-200/80 bg-white/95 px-6 py-4 shadow-[0_20px_70px_-46px_rgba(15,23,42,0.42)]">
            <div className="flex min-w-0 items-center gap-3">
              <UserAvatar name={profile?.display_name} email={user?.email} avatarUrl={profile?.avatar_url} className="h-11 w-11 rounded-2xl" />
              <div className="min-w-0 space-y-1">
                <Badge className="bg-slate-950 text-white hover:bg-slate-950" variant="secondary">
                  {roleLabel(user?.role)}
                </Badge>
                <p className="truncate text-sm text-slate-600">{profile?.display_name ?? user?.email ?? "Рабочее пространство SoundMarket"}</p>
              </div>
            </div>
            <Link href="/" className="shrink-0 text-sm font-medium text-slate-700 transition hover:text-slate-950">
              На главную
            </Link>
          </header>

          <main className="min-w-0">{children}</main>
        </div>
      </div>
    </div>
  );
}
