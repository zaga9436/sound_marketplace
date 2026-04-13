"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { adminApi } from "@/entities/admin/api/admin";
import { CardPreview } from "@/entities/card/ui/card-preview";
import { getErrorMessage } from "@/lib/api/errors";
import { AdminActionForm } from "@/widgets/admin/admin-action-form";
import { AdminPageHeader, AdminSection, StatusBadge } from "@/widgets/admin/admin-ui";

export function AdminCardDetailPage({ id }: { id: string }) {
  const queryClient = useQueryClient();
  const cardQuery = useQuery({
    queryKey: ["admin", "card", id],
    queryFn: () => adminApi.getCard(id)
  });

  const hideMutation = useMutation({
    mutationFn: (reason: string) => adminApi.hideCard(id, reason),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["admin", "card", id] });
      await queryClient.invalidateQueries({ queryKey: ["admin", "cards"] });
    }
  });
  const unhideMutation = useMutation({
    mutationFn: (reason: string) => adminApi.unhideCard(id, reason),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["admin", "card", id] });
      await queryClient.invalidateQueries({ queryKey: ["admin", "cards"] });
    }
  });

  if (cardQuery.isLoading) {
    return <div className="surface h-[320px] animate-pulse bg-white/70" />;
  }

  if (cardQuery.isError || !cardQuery.data) {
    return <p className="text-sm text-destructive">{getErrorMessage(cardQuery.error)}</p>;
  }

  const card = cardQuery.data;

  return (
    <div className="space-y-6">
      <AdminPageHeader title="Карточка в модерации" description="Полный просмотр карточки с актуальным статусом, медиа и действием hide/unhide." />

      <AdminSection className="space-y-5">
        <div className="grid gap-6 lg:grid-cols-[minmax(0,1.2fr)_minmax(280px,0.8fr)]">
          <CardPreview card={card} />
          <div className="space-y-4">
            <div className="flex flex-wrap gap-2">
              <StatusBadge tone={card.is_hidden ? "red" : "green"}>{card.is_hidden ? "Скрыта" : "Активна"}</StatusBadge>
              <StatusBadge tone="blue">{card.card_type === "offer" ? "Предложение" : "Запрос"}</StatusBadge>
              <StatusBadge tone="slate">{card.kind === "service" ? "Услуга" : "Продукт"}</StatusBadge>
            </div>
            <div className="space-y-2">
              <p className="text-sm text-slate-500">Автор</p>
              <p className="font-mono text-sm text-slate-700">{card.author_id}</p>
            </div>
            <div className="space-y-2">
              <p className="text-sm text-slate-500">ID карточки</p>
              <p className="font-mono text-sm text-slate-700">{card.id}</p>
            </div>
            {card.moderation_reason ? (
              <div className="rounded-2xl bg-rose-50 px-4 py-3 text-sm text-rose-700">Причина скрытия: {card.moderation_reason}</div>
            ) : null}

            {card.is_hidden ? (
              <AdminActionForm
                actionLabel="Вернуть в каталог"
                confirmLabel="Вернуть карточку"
                placeholder="Комментарий к возврату (необязательно)"
                optional
                pending={unhideMutation.isPending}
                onSubmit={(reason) => unhideMutation.mutate(reason)}
              />
            ) : (
              <AdminActionForm
                actionLabel="Скрыть карточку"
                confirmLabel="Скрыть карточку"
                placeholder="Укажите причину скрытия"
                pending={hideMutation.isPending}
                onSubmit={(reason) => hideMutation.mutate(reason)}
              />
            )}
          </div>
        </div>
      </AdminSection>

      {hideMutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(hideMutation.error)}</p> : null}
      {unhideMutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(unhideMutation.error)}</p> : null}
    </div>
  );
}
