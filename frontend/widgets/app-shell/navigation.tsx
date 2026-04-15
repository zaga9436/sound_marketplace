"use client";

import type { ComponentType } from "react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { usePathname, useRouter } from "next/navigation";
import { Bell, CreditCard, LayoutDashboard, LogOut, MessageSquare, Package2, PlusCircle, Shield, ShoppingBag, UserRound } from "lucide-react";

import { notificationsApi } from "@/entities/notification/api/notifications";
import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { UserAvatar } from "@/shared/ui/user-avatar";
import { cn } from "@/shared/utils/cn";

type NavItem = {
  href: string;
  label: string;
  icon: ComponentType<{ className?: string }>;
  roles?: Array<"customer" | "engineer" | "admin">;
};

const items: NavItem[] = [
  { href: "/dashboard", label: "Обзор", icon: LayoutDashboard, roles: ["customer", "engineer", "admin"] },
  { href: "/catalog", label: "Каталог", icon: ShoppingBag, roles: ["customer", "engineer"] },
  { href: "/cards/new", label: "Новая карточка", icon: PlusCircle, roles: ["customer", "engineer"] },
  { href: "/orders", label: "Заказы", icon: Package2, roles: ["customer", "engineer", "admin"] },
  { href: "/balance", label: "Баланс", icon: CreditCard, roles: ["customer", "engineer", "admin"] },
  { href: "/chats", label: "Чаты", icon: MessageSquare, roles: ["customer", "engineer", "admin"] },
  { href: "/notifications", label: "Уведомления", icon: Bell, roles: ["customer", "engineer", "admin"] },
  { href: "/profile", label: "Профиль", icon: UserRound, roles: ["customer", "engineer", "admin"] },
  { href: "/admin", label: "Админ", icon: Shield, roles: ["admin"] }
];

function roleLabel(role?: string | null) {
  if (role === "customer") return "Заказчик";
  if (role === "engineer") return "Инженер";
  if (role === "admin") return "Администратор";
  return "Гость";
}

export function RoleNavigation() {
  const pathname = usePathname();
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const profile = useAuthStore((state) => state.profile);
  const logout = useAuthStore((state) => state.logout);

  const notificationsQuery = useQuery({
    queryKey: ["notifications"],
    queryFn: () => notificationsApi.list(20),
    enabled: Boolean(user),
    refetchInterval: 5000
  });

  const visibleItems = items.filter((item) => !item.roles || (user?.role ? item.roles.includes(user.role) : false));

  return (
    <aside className="flex h-full flex-col gap-6 rounded-[2rem] border border-slate-200/80 bg-white/95 p-4 shadow-[0_24px_80px_-48px_rgba(15,23,42,0.42)]">
      <div className="space-y-4 rounded-[1.5rem] bg-slate-50 p-4">
        <Badge className="bg-slate-950 text-white hover:bg-slate-950">SoundMarket</Badge>
        <div className="flex items-center gap-3">
          <UserAvatar name={profile?.display_name} email={user?.email} avatarUrl={profile?.avatar_url} className="h-11 w-11 rounded-2xl" />
          <div className="min-w-0">
            <p className="truncate text-sm font-semibold text-slate-950">{profile?.display_name ?? "Пользователь"}</p>
            <p className="truncate text-xs text-slate-500">{user?.email ?? "guest@example.com"}</p>
            <p className="mt-1 text-xs font-medium text-slate-600">{roleLabel(user?.role)}</p>
          </div>
        </div>
      </div>

      <nav className="flex flex-1 flex-col gap-1">
        {visibleItems.map((item) => {
          const Icon = item.icon;
          const active = pathname === item.href || pathname.startsWith(`${item.href}/`);

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center gap-3 rounded-2xl px-3 py-2.5 text-sm font-medium transition-colors",
                active ? "bg-slate-950 text-white shadow-sm" : "text-slate-700 hover:bg-slate-100 hover:text-slate-950"
              )}
            >
              <Icon className="h-4 w-4" />
              <span className="flex-1">{item.label}</span>
              {item.href === "/notifications" && (notificationsQuery.data?.unread_count ?? 0) > 0 ? (
                <span className={cn("rounded-full px-2 py-0.5 text-xs", active ? "bg-white/20 text-white" : "bg-slate-950 text-white")}>
                  {notificationsQuery.data?.unread_count}
                </span>
              ) : null}
            </Link>
          );
        })}
      </nav>

      <Button
        variant="outline"
        className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100"
        onClick={() => {
          logout();
          router.push("/login");
        }}
      >
        <LogOut className="mr-2 h-4 w-4" />
        Выйти
      </Button>
    </aside>
  );
}
