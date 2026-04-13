"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { disputesApi } from "@/entities/dispute/api/disputes";
import { getErrorMessage } from "@/lib/api/errors";
import { Button } from "@/shared/ui/button";
import { Label } from "@/shared/ui/label";
import { Textarea } from "@/shared/ui/textarea";
import { Dispute, Order, User } from "@/shared/types/api";

const openSchema = z.object({
  reason: z.string().min(3, "Опишите причину спора.")
});

type OpenValues = z.infer<typeof openSchema>;

function canOpenDispute(order: Order, user?: User | null) {
  return !!user && (user.id === order.customer_id || user.id === order.engineer_id) && !["completed", "cancelled", "dispute"].includes(order.status);
}

function canCloseDispute(dispute: Dispute, user?: User | null) {
  return !!user && dispute.status === "open" && (user.role === "admin" || user.id === dispute.opened_by_user_id);
}

function resolutionLabel(resolution?: string) {
  if (resolution === "cancel_order") return "Отменить заказ и вернуть деньги";
  if (resolution === "complete_order") return "Завершить заказ";
  return "Без решения";
}

export function DisputeActions({
  order,
  dispute,
  user
}: {
  order: Order;
  dispute?: Dispute | null;
  user?: User | null;
}) {
  const queryClient = useQueryClient();

  const form = useForm<OpenValues>({
    resolver: zodResolver(openSchema),
    defaultValues: { reason: "" }
  });

  const openMutation = useMutation({
    mutationFn: (values: OpenValues) => disputesApi.open(order.id, values.reason),
    onSuccess: async () => {
      form.reset();
      await queryClient.invalidateQueries({ queryKey: ["order", order.id] });
      await queryClient.invalidateQueries({ queryKey: ["dispute", order.id] });
    }
  });

  const closeMutation = useMutation({
    mutationFn: (resolution: "complete_order" | "cancel_order") => disputesApi.close(order.id, resolution),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["order", order.id] });
      await queryClient.invalidateQueries({ queryKey: ["dispute", order.id] });
      await queryClient.invalidateQueries({ queryKey: ["balance"] });
    }
  });

  const onSubmit = form.handleSubmit((values) => openMutation.mutate(values));

  if (dispute) {
    return (
      <div className="space-y-4">
        <div className="rounded-2xl border border-rose-200 bg-rose-50 p-4">
          <p className="text-sm font-medium text-rose-900">Статус спора: {dispute.status === "open" ? "открыт" : "закрыт"}</p>
          <p className="mt-2 text-sm text-rose-800">{dispute.reason}</p>
          {dispute.resolution ? <p className="mt-2 text-sm text-rose-700">Решение: {resolutionLabel(dispute.resolution)}</p> : null}
        </div>

        {canCloseDispute(dispute, user) ? (
          <div className="space-y-3">
            {closeMutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(closeMutation.error)}</p> : null}
            <div className="flex flex-wrap gap-3">
              <Button onClick={() => closeMutation.mutate("cancel_order")} variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                Вернуть деньги и отменить
              </Button>
              <Button onClick={() => closeMutation.mutate("complete_order")} className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800">
                Завершить заказ
              </Button>
            </div>
          </div>
        ) : null}
      </div>
    );
  }

  if (!canOpenDispute(order, user)) {
    return <p className="text-sm text-slate-500">Спор можно открыть только участнику активного заказа.</p>;
  }

  return (
    <form onSubmit={onSubmit} className="space-y-3">
      <div className="space-y-2">
        <Label htmlFor="dispute-reason">Причина спора</Label>
        <Textarea id="dispute-reason" className="rounded-2xl border-slate-300" placeholder="Кратко опишите проблему: сроки, качество, несоответствие договоренностям." {...form.register("reason")} />
        {form.formState.errors.reason ? <p className="text-sm text-red-600">{form.formState.errors.reason.message}</p> : null}
      </div>
      {openMutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(openMutation.error)}</p> : null}
      <Button type="submit" className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800" disabled={openMutation.isPending}>
        {openMutation.isPending ? "Открываем спор..." : "Открыть спор"}
      </Button>
    </form>
  );
}
