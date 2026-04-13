"use client";

import Link from "next/link";
import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { adminApi } from "@/entities/admin/api/admin";
import { CardPreview } from "@/entities/card/ui/card-preview";
import { getErrorMessage } from "@/lib/api/errors";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { AdminActionForm } from "@/widgets/admin/admin-action-form";
import { AdminPageHeader, AdminSection, StatusBadge } from "@/widgets/admin/admin-ui";

export function AdminCardsPage() {
  const queryClient = useQueryClient();
  const [visibility, setVisibility] = useState("");
  const [cardType, setCardType] = useState("");
  const [authorId, setAuthorId] = useState("");

  const cardsQuery = useQuery({
    queryKey: ["admin", "cards", visibility, cardType, authorId],
    queryFn: () =>
      adminApi.listCards({
        visibility: visibility || undefined,
        card_type: cardType || undefined,
        author_id: authorId || undefined,
        limit: "24",
        offset: "0"
      })
  });

  const hideMutation = useMutation({
    mutationFn: ({ id, reason }: { id: string; reason: string }) => adminApi.hideCard(id, reason),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["admin", "cards"] });
    }
  });

  const unhideMutation = useMutation({
    mutationFn: ({ id, reason }: { id: string; reason: string }) => adminApi.unhideCard(id, reason),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["admin", "cards"] });
    }
  });

  const cards = cardsQuery.data?.items ?? [];

  return (
    <div className="space-y-6">
      <AdminPageHeader title="Карточки" description="Публичный каталог с точки зрения модерации: фильтры, статусы и быстрые действия hide/unhide." />

      <AdminSection className="space-y-4">
        <div className="grid gap-3 md:grid-cols-4">
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-700">Видимость</label>
            <select className="h-11 rounded-2xl border border-slate-300 bg-white px-4 text-sm text-slate-900" value={visibility} onChange={(event) => setVisibility(event.target.value)}>
              <option value="">Все</option>
              <option value="visible">Активные</option>
              <option value="hidden">Скрытые</option>
            </select>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-700">Тип карточки</label>
            <select className="h-11 rounded-2xl border border-slate-300 bg-white px-4 text-sm text-slate-900" value={cardType} onChange={(event) => setCardType(event.target.value)}>
              <option value="">Все типы</option>
              <option value="offer">Предложения</option>
              <option value="request">Запросы</option>
            </select>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-700">Автор</label>
            <Input value={authorId} onChange={(event) => setAuthorId(event.target.value)} placeholder="ID автора" className="rounded-2xl border-slate-300" />
          </div>
          <div className="flex items-end">
            <Button variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100" onClick={() => { setVisibility(""); setCardType(""); setAuthorId(""); }}>
              Сбросить фильтры
            </Button>
          </div>
        </div>
      </AdminSection>

      {!cardsQuery.isLoading && cards.length === 0 ? (
        <AdminSection className="text-sm text-slate-500">По текущим фильтрам карточки не найдены.</AdminSection>
      ) : (
        <div className="grid gap-4 xl:grid-cols-2">
          {cards.map((card) => (
            <AdminSection key={card.id} className="space-y-4">
              <div className="grid gap-4 lg:grid-cols-[minmax(0,1.2fr)_minmax(220px,0.8fr)]">
                <CardPreview card={card} />
                <div className="space-y-3">
                  <div className="flex flex-wrap gap-2">
                    <StatusBadge tone={card.is_hidden ? "red" : "green"}>{card.is_hidden ? "Скрыта" : "Активна"}</StatusBadge>
                    <StatusBadge tone="blue">{card.card_type === "offer" ? "Предложение" : "Запрос"}</StatusBadge>
                  </div>
                  <p className="font-mono text-xs text-slate-500">{card.id}</p>
                  <p className="text-sm text-slate-500">Автор: {card.author_id}</p>
                  {card.moderation_reason ? <p className="rounded-2xl bg-rose-50 px-3 py-2 text-sm text-rose-700">Причина скрытия: {card.moderation_reason}</p> : null}
                  <div className="flex flex-wrap gap-2">
                    <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                      <Link href={`/admin/cards/${card.id}`}>Открыть</Link>
                    </Button>
                    {card.is_hidden ? (
                      <AdminActionForm
                        actionLabel="Вернуть в каталог"
                        confirmLabel="Вернуть"
                        placeholder="Комментарий к возврату (необязательно)"
                        optional
                        pending={unhideMutation.isPending}
                        onSubmit={(reason) => unhideMutation.mutate({ id: card.id, reason })}
                      />
                    ) : (
                      <AdminActionForm
                        actionLabel="Скрыть карточку"
                        confirmLabel="Скрыть"
                        placeholder="Укажите причину скрытия"
                        pending={hideMutation.isPending}
                        onSubmit={(reason) => hideMutation.mutate({ id: card.id, reason })}
                      />
                    )}
                  </div>
                </div>
              </div>
            </AdminSection>
          ))}
        </div>
      )}

      {cardsQuery.isError ? <p className="text-sm text-destructive">{getErrorMessage(cardsQuery.error)}</p> : null}
      {hideMutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(hideMutation.error)}</p> : null}
      {unhideMutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(unhideMutation.error)}</p> : null}
    </div>
  );
}
