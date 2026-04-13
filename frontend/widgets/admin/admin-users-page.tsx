"use client";

import Link from "next/link";
import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { adminApi } from "@/entities/admin/api/admin";
import { getErrorMessage } from "@/lib/api/errors";
import { Button } from "@/shared/ui/button";
import { AdminActionForm } from "@/widgets/admin/admin-action-form";
import { AdminPageHeader, AdminSection, StatusBadge, formatRole } from "@/widgets/admin/admin-ui";

export function AdminUsersPage() {
  const queryClient = useQueryClient();
  const [role, setRole] = useState("");
  const [status, setStatus] = useState("");

  const usersQuery = useQuery({
    queryKey: ["admin", "users", role, status],
    queryFn: () => adminApi.listUsers({ role: role || undefined, status: status || undefined })
  });

  const suspendMutation = useMutation({
    mutationFn: ({ id, reason }: { id: string; reason: string }) => adminApi.suspendUser(id, reason),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    }
  });

  const unsuspendMutation = useMutation({
    mutationFn: ({ id, reason }: { id: string; reason: string }) => adminApi.unsuspendUser(id, reason),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    }
  });

  return (
    <div className="space-y-6">
      <AdminPageHeader title="Пользователи" description="Список участников платформы с основными moderation actions для быстрой проверки и блокировки." />

      <AdminSection className="space-y-4">
        <div className="grid gap-3 md:grid-cols-3">
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-700">Роль</label>
            <select className="h-11 rounded-2xl border border-slate-300 bg-white px-4 text-sm text-slate-900" value={role} onChange={(event) => setRole(event.target.value)}>
              <option value="">Все роли</option>
              <option value="customer">Заказчики</option>
              <option value="engineer">Исполнители</option>
              <option value="admin">Администраторы</option>
            </select>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-700">Статус</label>
            <select className="h-11 rounded-2xl border border-slate-300 bg-white px-4 text-sm text-slate-900" value={status} onChange={(event) => setStatus(event.target.value)}>
              <option value="">Все статусы</option>
              <option value="active">Активные</option>
              <option value="suspended">Заблокированные</option>
            </select>
          </div>
          <div className="flex items-end">
            <Button variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100" onClick={() => { setRole(""); setStatus(""); }}>
              Сбросить фильтры
            </Button>
          </div>
        </div>
      </AdminSection>

      <AdminSection className="overflow-hidden p-0">
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-slate-50 text-left text-slate-500">
              <tr>
                <th className="px-5 py-4 font-medium">Email</th>
                <th className="px-5 py-4 font-medium">Роль</th>
                <th className="px-5 py-4 font-medium">Статус</th>
                <th className="px-5 py-4 font-medium">Дата регистрации</th>
                <th className="px-5 py-4 font-medium">Действия</th>
              </tr>
            </thead>
            <tbody>
              {(usersQuery.data ?? []).map((user) => (
                <tr key={user.id} className="border-t border-slate-200 align-top">
                  <td className="px-5 py-4">
                    <div className="space-y-1">
                      <Link href={`/admin/users/${user.id}`} className="font-semibold text-slate-950 hover:text-primary">
                        {user.email}
                      </Link>
                      <p className="font-mono text-xs text-slate-500">{user.id}</p>
                    </div>
                  </td>
                  <td className="px-5 py-4 text-slate-700">{formatRole(user.role)}</td>
                  <td className="px-5 py-4">
                    <StatusBadge tone={user.is_suspended ? "red" : "green"}>{user.is_suspended ? "Заблокирован" : "Активен"}</StatusBadge>
                  </td>
                  <td className="px-5 py-4 text-slate-700">{user.created_at ? new Date(user.created_at).toLocaleDateString("ru-RU") : "—"}</td>
                  <td className="px-5 py-4">
                    <div className="space-y-2">
                      <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                        <Link href={`/admin/users/${user.id}`}>Открыть</Link>
                      </Button>
                      {user.role !== "admin" ? (
                        user.is_suspended ? (
                          <AdminActionForm
                            actionLabel="Разблокировать"
                            confirmLabel="Подтвердить"
                            placeholder="Комментарий к разблокировке (необязательно)"
                            optional
                            pending={unsuspendMutation.isPending}
                            onSubmit={(reason) => unsuspendMutation.mutate({ id: user.id, reason })}
                          />
                        ) : (
                          <AdminActionForm
                            actionLabel="Заблокировать"
                            confirmLabel="Заблокировать"
                            placeholder="Укажите причину блокировки"
                            pending={suspendMutation.isPending}
                            onSubmit={(reason) => suspendMutation.mutate({ id: user.id, reason })}
                          />
                        )
                      ) : (
                        <p className="text-xs text-slate-500">Администраторов нельзя блокировать.</p>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {!usersQuery.isLoading && (usersQuery.data?.length ?? 0) === 0 ? (
          <div className="p-6 text-sm text-slate-500">Пользователи по текущим фильтрам не найдены.</div>
        ) : null}
      </AdminSection>

      {usersQuery.isError ? <p className="text-sm text-destructive">{getErrorMessage(usersQuery.error)}</p> : null}
      {suspendMutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(suspendMutation.error)}</p> : null}
      {unsuspendMutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(unsuspendMutation.error)}</p> : null}
    </div>
  );
}
