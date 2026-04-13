"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";

import { adminApi } from "@/entities/admin/api/admin";
import { Button } from "@/shared/ui/button";
import { UserAvatar } from "@/shared/ui/user-avatar";
import { AdminPageHeader, AdminSection, StatusBadge, formatRole } from "@/widgets/admin/admin-ui";

export function AdminOverviewPage() {
  const usersQuery = useQuery({
    queryKey: ["admin", "users", "overview"],
    queryFn: () => adminApi.listUsers()
  });
  const activeCardsQuery = useQuery({
    queryKey: ["admin", "cards", "visible", "overview"],
    queryFn: () => adminApi.listCards({ visibility: "visible", limit: "100", offset: "0" })
  });
  const hiddenCardsQuery = useQuery({
    queryKey: ["admin", "cards", "hidden", "overview"],
    queryFn: () => adminApi.listCards({ visibility: "hidden", limit: "100", offset: "0" })
  });
  const disputesQuery = useQuery({
    queryKey: ["admin", "disputes", "open", "overview"],
    queryFn: () => adminApi.listDisputes("open")
  });
  const actionsQuery = useQuery({
    queryKey: ["admin", "actions", "overview"],
    queryFn: () => adminApi.listActions({ limit: "8" })
  });

  const metrics = [
    { label: "Всего пользователей", value: usersQuery.data?.length ?? "—" },
    { label: "Активные карточки", value: activeCardsQuery.data?.total ?? "—" },
    { label: "Скрытые карточки", value: hiddenCardsQuery.data?.total ?? "—" },
    { label: "Открытые споры", value: disputesQuery.data?.length ?? "—" }
  ];

  return (
    <div className="space-y-6">
      <AdminPageHeader
        title="Админ-панель"
        description="Здесь собрана первая версия служебной части SoundMarket: пользователи, карточки, споры и история действий модерации."
        actions={
          <Button asChild className="rounded-2xl bg-slate-950 text-white hover:bg-slate-800">
            <Link href="/admin/users">Перейти к модерации</Link>
          </Button>
        }
      />

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {metrics.map((metric) => (
          <AdminSection key={metric.label} className="space-y-2">
            <p className="text-sm text-slate-500">{metric.label}</p>
            <p className="text-4xl font-semibold tracking-tight text-slate-950">{metric.value}</p>
          </AdminSection>
        ))}
      </div>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.1fr)_minmax(0,0.9fr)]">
        <AdminSection className="space-y-4">
          <div className="flex items-center justify-between gap-4">
            <h2 className="text-2xl font-semibold text-slate-950">Последние действия модерации</h2>
            <Link href="/admin/actions" className="text-sm font-medium text-primary">
              Открыть журнал
            </Link>
          </div>

          <div className="space-y-3">
            {(actionsQuery.data ?? []).map((action) => (
              <div key={action.id} className="rounded-2xl border border-slate-200 bg-slate-50/80 p-4">
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div className="space-y-1">
                    <p className="text-sm font-semibold text-slate-950">
                      {action.action} • {action.target_type} • {action.target_id}
                    </p>
                    <p className="text-sm text-slate-500">{action.reason || "Без дополнительной причины"}</p>
                  </div>
                  <StatusBadge tone="slate">{new Date(action.created_at).toLocaleString("ru-RU")}</StatusBadge>
                </div>
              </div>
            ))}

            {!actionsQuery.isLoading && (actionsQuery.data?.length ?? 0) === 0 ? (
              <div className="rounded-2xl border border-dashed border-slate-200 bg-white p-6 text-sm text-slate-500">Пока нет действий модерации.</div>
            ) : null}
          </div>
        </AdminSection>

        <AdminSection className="space-y-4">
          <div className="flex items-center justify-between gap-4">
            <h2 className="text-2xl font-semibold text-slate-950">Открытые споры</h2>
            <Link href="/admin/disputes" className="text-sm font-medium text-primary">
              Все споры
            </Link>
          </div>
          <div className="space-y-3">
            {(disputesQuery.data ?? []).slice(0, 5).map((dispute) => (
              <Link key={dispute.id} href={`/admin/disputes/${dispute.id}`} className="block rounded-2xl border border-slate-200 bg-slate-50/80 p-4 transition hover:bg-slate-100">
                <div className="flex items-center justify-between gap-3">
                  <div className="space-y-1">
                    <p className="text-sm font-semibold text-slate-950">Спор по заказу {dispute.order_id}</p>
                    <p className="line-clamp-2 text-sm text-slate-500">{dispute.reason}</p>
                  </div>
                  <StatusBadge tone={dispute.status === "open" ? "yellow" : "green"}>{dispute.status === "open" ? "Открыт" : "Закрыт"}</StatusBadge>
                </div>
              </Link>
            ))}

            {!disputesQuery.isLoading && (disputesQuery.data?.length ?? 0) === 0 ? (
              <div className="rounded-2xl border border-dashed border-slate-200 bg-white p-6 text-sm text-slate-500">Сейчас нет открытых споров.</div>
            ) : null}
          </div>
        </AdminSection>
      </div>

      <AdminSection className="space-y-4">
        <div className="flex items-center justify-between gap-4">
          <h2 className="text-2xl font-semibold text-slate-950">Недавно зарегистрированные пользователи</h2>
          <Link href="/admin/users" className="text-sm font-medium text-primary">
            Все пользователи
          </Link>
        </div>
        <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
          {(usersQuery.data ?? []).slice(0, 6).map((user) => (
            <Link key={user.id} href={`/admin/users/${user.id}`} className="rounded-2xl border border-slate-200 bg-slate-50/80 p-4 transition hover:bg-slate-100">
              <div className="flex items-center gap-3">
                <UserAvatar email={user.email} className="h-12 w-12 rounded-2xl" />
                <div className="min-w-0">
                  <p className="truncate text-sm font-semibold text-slate-950">{user.email}</p>
                  <p className="text-xs text-slate-500">{formatRole(user.role)}</p>
                </div>
              </div>
            </Link>
          ))}
        </div>
      </AdminSection>
    </div>
  );
}
