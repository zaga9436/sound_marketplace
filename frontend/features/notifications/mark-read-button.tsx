"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { notificationsApi } from "@/entities/notification/api/notifications";
import { Button } from "@/shared/ui/button";

export function MarkNotificationsReadButton({ ids }: { ids?: string[] }) {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: () => notificationsApi.markRead(ids),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["notifications"] });
    }
  });

  return (
    <Button
      variant="outline"
      className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100"
      onClick={() => mutation.mutate()}
      disabled={mutation.isPending}
    >
      {mutation.isPending ? "Обновляем..." : ids?.length ? "Пометить как прочитанное" : "Прочитать все"}
    </Button>
  );
}
