"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { adminApi } from "@/entities/admin/api/admin";
import { profilesApi } from "@/entities/profile/api/profiles";
import { getErrorMessage } from "@/lib/api/errors";
import { AdminActionForm } from "@/widgets/admin/admin-action-form";
import { AdminPageHeader, AdminSection, StatusBadge, formatRole } from "@/widgets/admin/admin-ui";
import { UserAvatar } from "@/shared/ui/user-avatar";

export function AdminUserDetailPage({ id }: { id: string }) {
  const queryClient = useQueryClient();
  const userQuery = useQuery({
    queryKey: ["admin", "user", id],
    queryFn: () => adminApi.getUser(id)
  });
  const profileQuery = useQuery({
    queryKey: ["admin", "user", id, "profile"],
    queryFn: () => profilesApi.getById(id)
  });

  const suspendMutation = useMutation({
    mutationFn: (reason: string) => adminApi.suspendUser(id, reason),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["admin", "user", id] });
      await queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    }
  });
  const unsuspendMutation = useMutation({
    mutationFn: (reason: string) => adminApi.unsuspendUser(id, reason),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["admin", "user", id] });
      await queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    }
  });

  if (userQuery.isLoading) {
    return <div className="surface h-[320px] animate-pulse bg-white/70" />;
  }

  if (userQuery.isError || !userQuery.data) {
    return <p className="text-sm text-destructive">{getErrorMessage(userQuery.error)}</p>;
  }

  const user = userQuery.data;
  const profile = profileQuery.data;

  return (
    <div className="space-y-6">
      <AdminPageHeader title="Карточка пользователя" description="Подробная информация о пользователе и основные moderation actions." />

      <AdminSection className="space-y-6">
        <div className="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
          <div className="flex items-center gap-4">
            <UserAvatar avatarUrl={profile?.avatar_url} name={profile?.display_name} email={user.email} className="h-24 w-24 rounded-[2rem] text-2xl" />
            <div className="space-y-2">
              <h2 className="text-2xl font-semibold text-slate-950">{profile?.display_name || user.email}</h2>
              <p className="text-sm text-slate-500">{user.email}</p>
              <div className="flex flex-wrap gap-2">
                <StatusBadge tone="blue">{formatRole(user.role)}</StatusBadge>
                <StatusBadge tone={user.is_suspended ? "red" : "green"}>{user.is_suspended ? "Заблокирован" : "Активен"}</StatusBadge>
              </div>
            </div>
          </div>

          {user.role !== "admin" ? (
            user.is_suspended ? (
              <AdminActionForm
                actionLabel="Разблокировать пользователя"
                confirmLabel="Разблокировать"
                placeholder="Комментарий к разблокировке (необязательно)"
                optional
                pending={unsuspendMutation.isPending}
                onSubmit={(reason) => unsuspendMutation.mutate(reason)}
              />
            ) : (
              <AdminActionForm
                actionLabel="Заблокировать пользователя"
                confirmLabel="Заблокировать"
                placeholder="Укажите причину блокировки"
                pending={suspendMutation.isPending}
                onSubmit={(reason) => suspendMutation.mutate(reason)}
              />
            )
          ) : null}
        </div>

        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          <div>
            <p className="text-sm text-slate-500">ID пользователя</p>
            <p className="font-mono text-sm text-slate-700">{user.id}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">Дата регистрации</p>
            <p className="text-sm text-slate-700">{user.created_at ? new Date(user.created_at).toLocaleString("ru-RU") : "—"}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">Причина блокировки</p>
            <p className="text-sm text-slate-700">{user.suspension_reason || "Не указана"}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">Описание профиля</p>
            <p className="text-sm text-slate-700">{profile?.bio || "Пусто"}</p>
          </div>
        </div>
      </AdminSection>

      {suspendMutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(suspendMutation.error)}</p> : null}
      {unsuspendMutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(unsuspendMutation.error)}</p> : null}
    </div>
  );
}
