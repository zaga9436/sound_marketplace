"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { Music2 } from "lucide-react";

import { cardsApi } from "@/entities/card/api/cards";
import { profilesApi } from "@/entities/profile/api/profiles";
import { getErrorMessage } from "@/lib/api/errors";
import { AudioCoverPreview } from "@/shared/ui/audio-cover-preview";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";
import { CardDealPanel } from "@/widgets/cards/card-deal-panel";

function formatPrice(value: number) {
  return new Intl.NumberFormat("ru-RU", {
    style: "currency",
    currency: "RUB",
    maximumFractionDigits: 0
  }).format(value);
}

export function CardDetail({ id }: { id: string }) {
  const cardQuery = useQuery({
    queryKey: ["card", id],
    queryFn: () => cardsApi.getById(id)
  });

  const profileQuery = useQuery({
    queryKey: ["profile", cardQuery.data?.author_id],
    queryFn: () => profilesApi.getById(cardQuery.data!.author_id),
    enabled: Boolean(cardQuery.data?.author_id)
  });

  if (cardQuery.isLoading) {
    return <div className="surface h-[420px] animate-pulse bg-white/70" />;
  }

  if (cardQuery.isError) {
    return (
      <Card className="border-destructive/20 bg-white/95">
        <CardContent className="pt-6">
          <p className="text-destructive">{getErrorMessage(cardQuery.error)}</p>
        </CardContent>
      </Card>
    );
  }

  const card = cardQuery.data;
  if (!card) return null;
  const preview = card.preview_urls?.[0];
  const showAudioPreview = card.card_type === "offer" && card.kind === "product" && Boolean(preview);

  return (
    <div className="space-y-8">
      <div className="grid gap-8 xl:grid-cols-[minmax(0,1.35fr)_380px]">
        <Card className="overflow-hidden border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.32)]">
          <div className="p-4 pb-0">
            {showAudioPreview ? (
              <AudioCoverPreview coverUrl={card.cover_url} audioUrl={preview} title={card.title} className="aspect-[16/9]" />
            ) : (
              <div
                className="relative flex aspect-[16/9] items-end overflow-hidden rounded-[1.75rem] border border-slate-200 bg-[linear-gradient(145deg,rgba(15,23,42,0.92),rgba(51,65,85,0.9))] p-6 text-white"
                style={
                  card.cover_url
                    ? {
                        backgroundImage: `linear-gradient(180deg, rgba(15,23,42,0.12), rgba(15,23,42,0.78)), url(${card.cover_url})`,
                        backgroundSize: "cover",
                        backgroundPosition: "center"
                      }
                    : undefined
                }
              >
                <div className="relative space-y-3">
                  <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-white/14 backdrop-blur">
                    <Music2 className="h-5 w-5" />
                  </div>
                  <p className="max-w-2xl text-2xl font-semibold leading-tight">{card.title}</p>
                </div>
              </div>
            )}
          </div>
          <CardHeader className="space-y-4">
            <div className="flex flex-wrap gap-2">
              <Badge className="bg-slate-900/90 text-white" variant="secondary">
                {card.card_type === "offer" ? "Предложение" : "Запрос"}
              </Badge>
              <Badge variant="outline">{card.kind === "service" ? "Услуга" : "Продукт"}</Badge>
            </div>
            <CardTitle className="text-3xl text-slate-950">{card.title}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-5">
            <p className="text-base text-slate-700">{card.description}</p>
            {card.tags?.length ? (
              <div className="flex flex-wrap gap-2">
                {card.tags.map((tag) => (
                  <span key={tag} className="rounded-full bg-slate-100 px-3 py-1 text-sm text-slate-700">
                    #{tag}
                  </span>
                ))}
              </div>
            ) : null}
          </CardContent>
        </Card>

        <div className="space-y-4">
          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader>
              <CardTitle>Детали карточки</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm text-slate-500">Цена</span>
                <span className="text-xl font-semibold text-slate-950">{formatPrice(card.price)}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-slate-500">Создано</span>
                <span className="text-sm text-slate-700">{new Date(card.created_at).toLocaleDateString("ru-RU")}</span>
              </div>
              <div className="space-y-2">
                <span className="text-sm text-slate-500">Автор</span>
                {profileQuery.data ? (
                  <div className="space-y-1">
                    <p className="font-medium text-slate-950">{profileQuery.data.display_name}</p>
                    <p className="text-sm text-slate-500">
                      Рейтинг {profileQuery.data.rating.toFixed(1)} • {profileQuery.data.reviews_count} отзывов
                    </p>
                  </div>
                ) : (
                  <p className="text-sm text-slate-500">Загружаем профиль...</p>
                )}
                <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                  <Link href={`/profiles/${card.author_id}`}>Открыть профиль</Link>
                </Button>
              </div>
            </CardContent>
          </Card>

          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader>
              <CardTitle>Навигация</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <p>Отсюда удобно перейти в каталог или в кабинет и продолжить работу со своими сценариями.</p>
              <div className="flex flex-wrap gap-3">
                <Button asChild className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800">
                  <Link href="/catalog">Вернуться в каталог</Link>
                </Button>
                <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                  <Link href="/orders">Мои заказы</Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      <CardDealPanel card={card} />
    </div>
  );
}
