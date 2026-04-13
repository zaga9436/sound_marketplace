"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { Compass, LayoutDashboard, LogOut, Wallet } from "lucide-react";

import { useAuthStore } from "@/lib/auth/session-store";
import { Button } from "@/shared/ui/button";
import { UserAvatar } from "@/shared/ui/user-avatar";

export function PublicHeader() {
  const router = useRouter();
  const hydrated = useAuthStore((state) => state.hydrated);
  const user = useAuthStore((state) => state.user);
  const profile = useAuthStore((state) => state.profile);
  const logout = useAuthStore((state) => state.logout);

  return (
    <header className="sticky top-0 z-40 border-b border-slate-200/70 bg-white/88 backdrop-blur-xl">
      <div className="container flex items-center justify-between gap-4 py-4">
        <div className="flex items-center gap-5">
          <Link href="/" className="inline-flex items-center gap-3">
            <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-950 text-sm font-semibold text-white shadow-[0_16px_36px_-20px_rgba(15,23,42,0.7)]">
              SM
            </div>
            <div className="hidden sm:block">
              <p className="text-sm font-semibold text-slate-950">SoundMarket</p>
              <p className="text-xs text-slate-500">Музыкальный маркетплейс и сервис сделок</p>
            </div>
          </Link>

          <nav className="hidden items-center gap-5 text-sm text-slate-600 md:flex">
            <Link className="transition hover:text-slate-950" href="/">
              Главная
            </Link>
            <Link className="transition hover:text-slate-950" href="/catalog">
              Каталог
            </Link>
            {user ? (
              <>
                <Link className="transition hover:text-slate-950" href="/orders">
                  Заказы
                </Link>
                <Link className="transition hover:text-slate-950" href="/balance">
                  Баланс
                </Link>
              </>
            ) : null}
          </nav>
        </div>

        {!hydrated ? (
          <div className="flex items-center gap-2">
            <div className="h-11 w-28 animate-pulse rounded-2xl bg-slate-100" />
            <div className="h-11 w-36 animate-pulse rounded-2xl bg-slate-100" />
          </div>
        ) : user ? (
          <div className="flex items-center gap-2">
            <Link
              href={user.role === "admin" ? "/admin" : "/dashboard"}
              className="hidden items-center gap-3 rounded-2xl border border-slate-200 bg-white px-3 py-2 shadow-sm transition hover:bg-slate-50 sm:flex"
            >
              <UserAvatar
                avatarUrl={profile?.avatar_url}
                name={profile?.display_name}
                email={user.email}
                className="h-10 w-10 rounded-xl"
              />
              <div className="min-w-0 text-left">
                <p className="truncate text-sm font-semibold text-slate-950">{profile?.display_name || user.email}</p>
                <p className="truncate text-xs text-slate-500">
                  {user.role === "customer" ? "Заказчик" : user.role === "engineer" ? "Исполнитель" : "Администратор"}
                </p>
              </div>
            </Link>

            <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
              <Link href={user.role === "admin" ? "/admin" : "/dashboard"}>
                <LayoutDashboard className="mr-2 h-4 w-4" />
                Кабинет
              </Link>
            </Button>

            <Button asChild variant="outline" className="hidden rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100 md:inline-flex">
              <Link href="/balance">
                <Wallet className="mr-2 h-4 w-4" />
                Баланс
              </Link>
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
            <Button asChild variant="ghost" className="rounded-2xl text-slate-700 hover:bg-slate-100">
              <Link href="/login">
                <Compass className="mr-2 h-4 w-4" />
                Войти
              </Link>
            </Button>
            <Button asChild className="rounded-2xl bg-slate-950 text-white hover:bg-slate-800">
              <Link href="/register">Создать аккаунт</Link>
            </Button>
          </div>
        )}
      </div>
    </header>
  );
}
