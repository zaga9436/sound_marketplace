"use client";

import Link from "next/link";
import { PropsWithChildren } from "react";
import { Activity, Gavel, LayoutDashboard, Library, LogOut, ShieldCheck, Users } from "lucide-react";
import { usePathname, useRouter } from "next/navigation";

import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { UserAvatar } from "@/shared/ui/user-avatar";
import { cn } from "@/shared/utils/cn";

const adminNav = [
  { href: "/admin", label: "Обзор", icon: LayoutDashboard },
  { href: "/admin/users", label: "Пользователи", icon: Users },
  { href: "/admin/cards", label: "Карточки", icon: Library },
  { href: "/admin/disputes", label: "Споры", icon: Gavel },
  { href: "/admin/actions", label: "Действия", icon: Activity }
];

export function AdminShell({ children }: PropsWithChildren) {
  const pathname = usePathname();
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const profile = useAuthStore((state) => state.profile);
  const logout = useAuthStore((state) => state.logout);

  return (
    <div className="page-frame">
      <div className="container grid min-h-screen gap-6 py-6 lg:grid-cols-[296px_minmax(0,1fr)]">
        <aside className="hidden lg:block">
          <div className="surface flex h-full flex-col gap-6 p-5">
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-slate-950 text-white shadow-[0_18px_40px_-24px_rgba(15,23,42,0.8)]">
                  <ShieldCheck className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-sm font-semibold text-slate-950">SoundMarket Admin</p>
                  <p className="text-xs text-slate-500">Служебная панель модерации</p>
                </div>
              </div>

              <div className="rounded-[1.5rem] border border-slate-200 bg-slate-50/80 p-4">
                <div className="flex items-center gap-3">
                  <UserAvatar avatarUrl={profile?.avatar_url} name={profile?.display_name} email={user?.email} className="h-12 w-12 rounded-2xl" />
                  <div className="min-w-0">
                    <p className="truncate text-sm font-semibold text-slate-950">{profile?.display_name || user?.email || "Администратор"}</p>
                    <p className="truncate text-xs text-slate-500">Полный доступ к модерации</p>
                  </div>
                </div>
              </div>
            </div>

            <nav className="flex flex-1 flex-col gap-1">
              {adminNav.map((item) => {
                const Icon = item.icon;
                const active = pathname === item.href || (item.href !== "/admin" && pathname.startsWith(`${item.href}/`));

                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={cn(
                      "flex items-center gap-3 rounded-2xl px-4 py-3 text-sm font-medium transition",
                      active ? "bg-slate-950 text-white shadow-sm" : "text-slate-700 hover:bg-slate-100 hover:text-slate-950"
                    )}
                  >
                    <Icon className="h-4 w-4" />
                    {item.label}
                  </Link>
                );
              })}
            </nav>

            <div className="space-y-2">
              <Button asChild variant="outline" className="w-full rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                <Link href="/">На главную</Link>
              </Button>
              <Button
                variant="ghost"
                className="w-full rounded-2xl text-slate-700 hover:bg-slate-100"
                onClick={() => {
                  logout();
                  router.push("/login");
                }}
              >
                <LogOut className="mr-2 h-4 w-4" />
                Выйти
              </Button>
            </div>
          </div>
        </aside>

        <div className="flex min-w-0 flex-col gap-6">
          <header className="surface flex flex-wrap items-center justify-between gap-4 px-6 py-4">
            <div className="space-y-1">
              <Badge variant="secondary" className="bg-slate-950 text-white">
                Администратор
              </Badge>
              <p className="text-sm text-slate-500">Управление пользователями, карточками, спорами и действиями модерации.</p>
            </div>
            <div className="flex items-center gap-3 rounded-2xl border border-slate-200 bg-slate-50 px-3 py-2">
              <UserAvatar avatarUrl={profile?.avatar_url} name={profile?.display_name} email={user?.email} className="h-10 w-10 rounded-xl" />
              <div className="min-w-0">
                <p className="truncate text-sm font-semibold text-slate-950">{profile?.display_name || user?.email || "Администратор"}</p>
                <p className="truncate text-xs text-slate-500">{user?.email}</p>
              </div>
            </div>
          </header>

          <main className="min-w-0">{children}</main>
        </div>
      </div>
    </div>
  );
}
