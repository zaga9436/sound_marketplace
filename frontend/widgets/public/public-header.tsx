"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { LogOut, UserRound } from "lucide-react";

import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";

export function PublicHeader() {
  const router = useRouter();
  const hydrated = useAuthStore((state) => state.hydrated);
  const user = useAuthStore((state) => state.user);
  const profile = useAuthStore((state) => state.profile);
  const logout = useAuthStore((state) => state.logout);

  return (
    <header className="border-b border-border/70 bg-white/85 backdrop-blur">
      <div className="container flex items-center justify-between py-4">
        <div className="flex items-center gap-4">
          <Badge>SoundMarket</Badge>
          <nav className="hidden gap-5 text-sm text-muted-foreground md:flex">
            <Link href="/">Главная</Link>
            <Link href="/catalog">Каталог</Link>
            {user ? <Link href="/orders">Заказы</Link> : null}
          </nav>
        </div>

        {!hydrated ? (
          <div className="flex items-center gap-2">
            <div className="h-10 w-28 animate-pulse rounded-2xl bg-slate-100" />
            <div className="h-10 w-28 animate-pulse rounded-2xl bg-slate-100" />
          </div>
        ) : user ? (
          <div className="flex items-center gap-3">
            <Link href="/dashboard" className="flex items-center gap-3 rounded-2xl border border-slate-200 bg-white px-3 py-2 shadow-sm transition hover:bg-slate-50">
              <div className="flex h-9 w-9 items-center justify-center rounded-full bg-slate-900 text-white">
                <UserRound className="h-4 w-4" />
              </div>
              <div className="hidden text-left sm:block">
                <p className="text-sm font-medium text-slate-900">{profile?.display_name || user.email}</p>
                <p className="text-xs text-slate-500 capitalize">{user.role}</p>
              </div>
            </Link>
            <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
              <Link href="/profile">Профиль</Link>
            </Button>
            <Button
              variant="ghost"
              className="rounded-2xl text-slate-700 hover:bg-slate-100"
              onClick={() => {
                logout();
                router.push("/");
              }}
            >
              <LogOut className="mr-2 h-4 w-4" />
              Выйти
            </Button>
          </div>
        ) : (
          <div className="flex items-center gap-2">
            <Button asChild variant="ghost">
              <Link href="/login">Войти</Link>
            </Button>
            <Button asChild variant="secondary">
              <Link href="/register">Создать аккаунт</Link>
            </Button>
          </div>
        )}
      </div>
    </header>
  );
}
