"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { ordersApi } from "@/entities/order/api/orders";
import { getErrorMessage } from "@/lib/api/errors";
import { useAuthStore } from "@/lib/auth/session-store";
import { Button } from "@/shared/ui/button";
import { Order, OrderStatus } from "@/shared/types/api";

type ActionItem = {
  status: OrderStatus;
  label: string;
};

function getAvailableActions(order: Order, role?: string | null): ActionItem[] {
  if (role === "engineer") {
    if (order.status === "on_hold") return [{ status: "in_progress", label: "Взять в работу" }];
    if (order.status === "in_progress") return [{ status: "review", label: "Отправить на проверку" }];
  }

  if (role === "customer") {
    if (order.status === "review") return [{ status: "completed", label: "Подтвердить выполнение" }];
  }

  return [];
}

export function OrderStatusActions({ order }: { order: Order }) {
  const queryClient = useQueryClient();
  const role = useAuthStore((state) => state.user?.role);
  const actions = getAvailableActions(order, role);

  const mutation = useMutation({
    mutationFn: (status: OrderStatus) => ordersApi.updateStatus(order.id, status),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["orders"] });
      await queryClient.invalidateQueries({ queryKey: ["order", order.id] });
    }
  });

  if (actions.length === 0) {
    return <p className="text-sm text-slate-500">Сейчас для вашей роли нет доступных действий по этому заказу.</p>;
  }

  return (
    <div className="space-y-3">
      {mutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(mutation.error)}</p> : null}
      <div className="flex flex-wrap gap-3">
        {actions.map((action) => (
          <Button
            key={action.status}
            onClick={() => mutation.mutate(action.status)}
            disabled={mutation.isPending}
            className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800"
          >
            {mutation.isPending ? "Обновляем..." : action.label}
          </Button>
        ))}
      </div>
    </div>
  );
}
