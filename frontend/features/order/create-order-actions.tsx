"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

import { ordersApi } from "@/entities/order/api/orders";
import { getErrorMessage } from "@/lib/api/errors";
import { Button } from "@/shared/ui/button";

export function CreateOrderFromBidButton({ bidId }: { bidId: string }) {
  const router = useRouter();
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: () => ordersApi.createFromBid(bidId),
    onSuccess: async (order) => {
      await queryClient.invalidateQueries({ queryKey: ["balance"] });
      await queryClient.invalidateQueries({ queryKey: ["orders"] });
      router.push(`/orders/${order.id}`);
    }
  });

  return (
    <div className="space-y-2">
      {mutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(mutation.error)}</p> : null}
      <Button
        onClick={() => mutation.mutate()}
        className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800"
        disabled={mutation.isPending}
      >
        {mutation.isPending ? "Создаем заказ..." : "Выбрать отклик"}
      </Button>
    </div>
  );
}

export function CreateOrderFromOfferButton({
  cardId,
  label = "Создать заказ",
  pendingLabel = "Создаем заказ..."
}: {
  cardId: string;
  label?: string;
  pendingLabel?: string;
}) {
  const router = useRouter();
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: () => ordersApi.createFromOffer(cardId),
    onSuccess: async (order) => {
      await queryClient.invalidateQueries({ queryKey: ["balance"] });
      await queryClient.invalidateQueries({ queryKey: ["orders"] });
      router.push(`/orders/${order.id}`);
    }
  });

  return (
    <div className="space-y-2">
      {mutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(mutation.error)}</p> : null}
      <Button
        onClick={() => mutation.mutate()}
        className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800"
        disabled={mutation.isPending}
      >
        {mutation.isPending ? pendingLabel : label}
      </Button>
    </div>
  );
}
