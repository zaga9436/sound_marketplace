"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { adminApi } from "@/entities/admin/api/admin";
import { getErrorMessage } from "@/lib/api/errors";
import { AdminActionForm } from "@/widgets/admin/admin-action-form";
import { AdminPageHeader, AdminSection, StatusBadge } from "@/widgets/admin/admin-ui";

export function AdminDisputeDetailPage({ id }: { id: string }) {
  const queryClient = useQueryClient();
  const disputeQuery = useQuery({
    queryKey: ["admin", "dispute", id],
    queryFn: () => adminApi.getDispute(id)
  });

  const closeMutation = useMutation({
    mutationFn: ({ resolution, reason }: { resolution: "complete_order" | "cancel_order"; reason: string }) =>
      adminApi.closeDispute(id, resolution, reason),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["admin", "dispute", id] });
      await queryClient.invalidateQueries({ queryKey: ["admin", "disputes"] });
      await queryClient.invalidateQueries({ queryKey: ["admin", "actions"] });
    }
  });

  if (disputeQuery.isLoading) {
    return <div className="surface h-[320px] animate-pulse bg-white/70" />;
  }

  if (disputeQuery.isError || !disputeQuery.data) {
    return <p className="text-sm text-destructive">{getErrorMessage(disputeQuery.error)}</p>;
  }

  const dispute = disputeQuery.data;

  return (
    <div className="space-y-6">
      <AdminPageHeader title="Детали спора" description="Решение по конфликту: завершить заказ или отменить его с возвратом средств." />

      <AdminSection className="space-y-5">
        <div className="flex flex-wrap gap-2">
          <StatusBadge tone={dispute.status === "open" ? "yellow" : "green"}>{dispute.status === "open" ? "Открыт" : "Закрыт"}</StatusBadge>
          {dispute.resolution ? <StatusBadge tone="blue">{dispute.resolution === "complete_order" ? "Завершить заказ" : "Отменить заказ"}</StatusBadge> : null}
        </div>
        <div className="space-y-2">
          <h2 className="text-2xl font-semibold text-slate-950">Заказ {dispute.order_id}</h2>
          <p className="max-w-4xl text-sm leading-7 text-slate-600">{dispute.reason}</p>
        </div>
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          <div>
            <p className="text-sm text-slate-500">ID спора</p>
            <p className="font-mono text-sm text-slate-700">{dispute.id}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">Открыл</p>
            <p className="font-mono text-sm text-slate-700">{dispute.opened_by_user_id}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">Создан</p>
            <p className="text-sm text-slate-700">{new Date(dispute.created_at).toLocaleString("ru-RU")}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">Закрыт</p>
            <p className="text-sm text-slate-700">{dispute.closed_at ? new Date(dispute.closed_at).toLocaleString("ru-RU") : "Пока нет"}</p>
          </div>
        </div>

        {dispute.status === "open" ? (
          <div className="grid gap-3 lg:grid-cols-2">
            <AdminActionForm
              actionLabel="Закрыть с завершением заказа"
              confirmLabel="Завершить заказ"
              placeholder="Комментарий администратора к закрытию спора"
              optional
              pending={closeMutation.isPending}
              onSubmit={(reason) => closeMutation.mutate({ resolution: "complete_order", reason })}
            />
            <AdminActionForm
              actionLabel="Закрыть с отменой заказа"
              confirmLabel="Отменить заказ"
              placeholder="Причина отмены и возврата средств"
              optional
              pending={closeMutation.isPending}
              onSubmit={(reason) => closeMutation.mutate({ resolution: "cancel_order", reason })}
            />
          </div>
        ) : null}
      </AdminSection>

      {closeMutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(closeMutation.error)}</p> : null}
    </div>
  );
}
