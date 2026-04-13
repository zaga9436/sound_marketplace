"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { usePathname, useRouter } from "next/navigation";
import { Bell, LayoutDashboard, LogOut, MessageSquare, Package2, PlusCircle, Shield, ShoppingBag, UserRound } from "lucide-react";

import { notificationsApi } from "@/entities/notification/api/notifications";
import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { cn } from "@/shared/utils/cn";

type NavItem = {
  href: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  roles?: Array<"customer" | "engineer" | "admin">;
};

const items: NavItem[] = [
  { href: "/dashboard", label: "Обзор", icon: LayoutDashboard, roles: ["customer", "engineer", "admin"] },
  { href: "/catalog", label: "Каталог", icon: ShoppingBag, roles: ["customer", "engineer"] },
  { href: "/cards/new", label: "Новая карточка", icon: PlusCircle, roles: ["customer", "engineer"] },
  { href: "/orders", label: "Заказы", icon: Package2, roles: ["customer", "engineer", "admin"] },
  { href: "/chats", label: "Чаты", icon: MessageSquare, roles: ["customer", "engineer", "admin"] },
  { href: "/notifications", label: "Уведомления", icon: Bell, roles: ["customer", "engineer", "admin"] },
  { href: "/profile", label: "Профиль", icon: UserRound, roles: ["customer", "engineer", "admin"] },
  { href: "/admin", label: "Админ", icon: Shield, roles: ["admin"] }
];

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
    <aside className="surface flex h-full flex-col gap-6 p-4">
      <div className="space-y-2">
        <Badge>SoundMarket</Badge>
        <div>
          <p className="text-sm font-medium text-foreground">{profile?.display_name ?? "Пользователь"}</p>
          <p className="text-sm text-muted-foreground">{user?.email ?? "guest@example.com"}</p>
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
                "flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors",
                active ? "bg-primary text-primary-foreground" : "text-foreground/70 hover:bg-secondary hover:text-foreground"
              )}
            >
              <Icon className="h-4 w-4" />
              <span className="flex-1">{item.label}</span>
              {item.href === "/notifications" && (notificationsQuery.data?.unread_count ?? 0) > 0 ? (
                <span className="rounded-full bg-white/20 px-2 py-0.5 text-xs text-current">{notificationsQuery.data?.unread_count}</span>
              ) : null}
            </Link>
          );
        })}
      </nav>

      <Button
        variant="outline"
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
