"use client";

import { useMemo } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";

import { cardsApi } from "@/entities/card/api/cards";
import { CardPreview } from "@/entities/card/ui/card-preview";
import { getErrorMessage } from "@/lib/api/errors";
import { Button } from "@/shared/ui/button";
import { Card, CardContent } from "@/shared/ui/card";

export function CatalogResults() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const pathname = usePathname();

  const params = useMemo(
    () => ({
      q: searchParams.get("q") ?? "",
      card_type: searchParams.get("card_type") ?? "",
      kind: searchParams.get("kind") ?? "",
      min_price: searchParams.get("min_price") ?? "",
      max_price: searchParams.get("max_price") ?? "",
      tag: searchParams.get("tag") ?? "",
      sort_by: searchParams.get("sort_by") ?? "created_at",
      sort_order: searchParams.get("sort_order") ?? "desc",
      limit: searchParams.get("limit") ?? "12",
      offset: searchParams.get("offset") ?? "0"
    }),
    [searchParams]
  );

  const query = useQuery({
    queryKey: ["catalog", params],
    queryFn: () => cardsApi.list(params)
  });

  function goToOffset(nextOffset: number) {
    const qs = new URLSearchParams(searchParams.toString());
    qs.set("offset", String(Math.max(0, nextOffset)));
    router.replace(`${pathname}?${qs.toString()}`);
  }

  if (query.isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        {Array.from({ length: 6 }).map((_, index) => (
          <div key={index} className="surface h-[360px] animate-pulse bg-white/70" />
        ))}
      </div>
    );
  }

  if (query.isError) {
    return (
      <Card className="border-destructive/20 bg-white/95">
        <CardContent className="pt-6">
          <p className="text-destructive">{getErrorMessage(query.error)}</p>
        </CardContent>
      </Card>
    );
  }

  if (!query.data) return null;

  const { items, total, limit, offset } = query.data;
  const currentPage = Math.floor(offset / limit) + 1;
  const totalPages = Math.max(1, Math.ceil(total / limit));

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
        <div>
          <p className="text-sm uppercase tracking-[0.18em] text-slate-500">Каталог</p>
          <h2 className="mt-1">Публичные карточки</h2>
        </div>
        <p className="text-sm text-slate-500">
          {total} карточек, страница {currentPage} из {totalPages}
        </p>
      </div>

      {items.length === 0 ? (
        <Card className="border-slate-200/80 bg-white/95">
          <CardContent className="pt-6">
            <p>По текущим фильтрам карточек не найдено.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {items.map((card) => (
            <CardPreview key={card.id} card={card} />
          ))}
        </div>
      )}

      <div className="flex items-center justify-between gap-3">
        <Button variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100" onClick={() => goToOffset(offset - limit)} disabled={offset <= 0}>
          Назад
        </Button>
        <Button variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100" onClick={() => goToOffset(offset + limit)} disabled={offset + limit >= total}>
          Дальше
        </Button>
      </div>
    </div>
  );
}
