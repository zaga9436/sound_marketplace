"use client";

import Link from "next/link";

import { AudioCoverPreview } from "@/shared/ui/audio-cover-preview";
import { Badge } from "@/shared/ui/badge";
import { Card as UICard, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";
import { type Card } from "@/shared/types/api";

function formatPrice(value: number) {
  return new Intl.NumberFormat("ru-RU", {
    style: "currency",
    currency: "RUB",
    maximumFractionDigits: 0
  }).format(value);
}

export function CardPreview({ card }: { card: Card }) {
  const preview = card.preview_urls?.[0];

  return (
    <Link href={`/cards/${card.id}`} className="block">
      <UICard className="h-full overflow-hidden border-slate-200/80 bg-white/95 shadow-[0_18px_50px_-28px_rgba(15,23,42,0.28)] transition-transform hover:-translate-y-1">
        <div className="p-3 pb-0">
          <AudioCoverPreview cardId={card.id} audioUrl={preview} title={card.title} compact />
        </div>

        <CardHeader className="space-y-3">
          <div className="flex flex-wrap gap-2">
            <Badge className="bg-slate-900/90 text-white" variant="secondary">
              {card.card_type === "offer" ? "Предложение" : "Запрос"}
            </Badge>
            <Badge variant="outline">{card.kind === "service" ? "Услуга" : "Продукт"}</Badge>
          </div>
          <CardTitle className="line-clamp-2 text-lg text-slate-950">{card.title}</CardTitle>
        </CardHeader>

        <CardContent className="space-y-4">
          <p className="line-clamp-3 text-sm text-slate-600">{card.description}</p>

          {card.tags?.length ? (
            <div className="flex flex-wrap gap-2">
              {card.tags.slice(0, 4).map((tag) => (
                <span key={tag} className="rounded-full bg-slate-100 px-2.5 py-1 text-xs text-slate-700">
                  #{tag}
                </span>
              ))}
            </div>
          ) : null}

          <div className="flex items-center justify-between">
            <span className="text-lg font-semibold text-slate-950">{formatPrice(card.price)}</span>
            <span className="text-xs text-slate-500">{new Date(card.created_at).toLocaleDateString("ru-RU")}</span>
          </div>
        </CardContent>
      </UICard>
    </Link>
  );
}
