"use client";

import Link from "next/link";
import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { adminApi } from "@/entities/admin/api/admin";
import { getErrorMessage } from "@/lib/api/errors";
import { AdminActionForm } from "@/widgets/admin/admin-action-form";
import { AdminPageHeader, AdminSection, StatusBadge } from "@/widgets/admin/admin-ui";

export function AdminDisputesPage() {
  const queryClient = useQueryClient();
  const [status, setStatus] = useState("");

  const disputesQuery = useQuery({
    queryKey: ["admin", "disputes", status],
    queryFn: () => adminApi.listDisputes(status || undefined)
  });

  const closeMutation = useMutation({
    mutationFn: ({ id, resolution, reason }: { id: string; resolution: "complete_order" | "cancel_order"; reason: string }) =>
      adminApi.closeDispute(id, resolution, reason),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["admin", "disputes"] });
      await queryClient.invalidateQueries({ queryKey: ["admin", "actions"] });
    }
  });

  return (
    <div className="space-y-6">
      <AdminPageHeader title="Споры" description="Рабочая зона для разбора конфликтов по заказам: причины, статусы и финальное решение." />

      <AdminSection className="space-y-4">
        <div className="flex flex-wrap items-end gap-3">
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-700">Статус</label>
            <select className="h-11 rounded-2xl border border-slate-300 bg-white px-4 text-sm text-slate-900" value={status} onChange={(event) => setStatus(event.target.value)}>
              <option value="">Все</option>
              <option value="open">Открытые</option>
              <option value="closed">Закрытые</option>
            </select>
          </div>
        </div>
      </AdminSection>

      <div className="space-y-4">
        {(disputesQuery.data ?? []).map((dispute) => (
          <AdminSection key={dispute.id} className="space-y-4">
            <div className="flex flex-wrap items-start justify-between gap-4">
              <div className="space-y-2">
                <div className="flex flex-wrap gap-2">
                  <StatusBadge tone={dispute.status === "open" ? "yellow" : "green"}>{dispute.status === "open" ? "Открыт" : "Закрыт"}</StatusBadge>
                  {dispute.resolution ? <StatusBadge tone="blue">{dispute.resolution === "complete_order" ? "Завершить заказ" : "Отменить заказ"}</StatusBadge> : null}
                </div>
                <h2 className="text-xl font-semibold text-slate-950">Спор по заказу {dispute.order_id}</h2>
                <p className="max-w-3xl text-sm leading-7 text-slate-600">{dispute.reason}</p>
              </div>
              <Link href={`/admin/disputes/${dispute.id}`} className="text-sm font-medium text-primary">
                Открыть детали
              </Link>
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
                  onSubmit={(reason) => closeMutation.mutate({ id: dispute.id, resolution: "complete_order", reason })}
                />
                <AdminActionForm
                  actionLabel="Закрыть с отменой заказа"
                  confirmLabel="Отменить заказ"
                  placeholder="Причина отмены и возврата средств"
                  optional
                  pending={closeMutation.isPending}
                  onSubmit={(reason) => closeMutation.mutate({ id: dispute.id, resolution: "cancel_order", reason })}
                />
              </div>
            ) : null}
          </AdminSection>
        ))}

        {!disputesQuery.isLoading && (disputesQuery.data?.length ?? 0) === 0 ? (
          <AdminSection className="text-sm text-slate-500">По текущему фильтру споры не найдены.</AdminSection>
        ) : null}
      </div>

      {disputesQuery.isError ? <p className="text-sm text-destructive">{getErrorMessage(disputesQuery.error)}</p> : null}
      {closeMutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(closeMutation.error)}</p> : null}
    </div>
  );
}
