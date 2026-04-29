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

function getAvailableActions(order: Order, role?: string | null, readyProduct = false): ActionItem[] {
  if (role === "engineer") {
    if (order.status === "on_hold") {
      return [{ status: "in_progress", label: readyProduct ? "Подготовить выдачу" : "Взять в работу" }];
    }
    if (order.status === "in_progress") {
      return [{ status: "review", label: readyProduct ? "Открыть проверку покупки" : "Отправить на проверку" }];
    }
  }

  if (role === "customer") {
    if (order.status === "review") {
      return [{ status: "completed", label: readyProduct ? "Завершить покупку" : "Подтвердить выполнение" }];
    }
  }

  return [];
}

export function OrderStatusActions({ order, readyProduct = false }: { order: Order; readyProduct?: boolean }) {
  const queryClient = useQueryClient();
  const role = useAuthStore((state) => state.user?.role);
  const actions = getAvailableActions(order, role, readyProduct);

  const mutation = useMutation({
    mutationFn: (status: OrderStatus) => ordersApi.updateStatus(order.id, status),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["orders"] });
      await queryClient.invalidateQueries({ queryKey: ["order", order.id] });
      await queryClient.invalidateQueries({ queryKey: ["balance"] });
    }
  });

  if (actions.length === 0) {
    return (
      <p className="text-sm leading-6 text-slate-500">
        {readyProduct
          ? "На этом этапе дополнительных действий не требуется. Когда статус покупки изменится, здесь появятся доступные шаги."
          : "Сейчас никаких действий выполнять не нужно. Когда заказ перейдет на следующий этап, здесь появятся доступные шаги."}
      </p>
    );
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
