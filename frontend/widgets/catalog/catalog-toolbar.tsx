"use client";

import { useMemo, useState, useTransition } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { Search, SlidersHorizontal } from "lucide-react";

import { Button } from "@/shared/ui/button";
import { Card, CardContent } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";

function setQueryParam(searchParams: URLSearchParams, key: string, value: string) {
  if (!value) searchParams.delete(key);
  else searchParams.set(key, value);
}

const selectClassName =
  "flex h-11 w-full rounded-2xl border border-slate-300 bg-white px-4 py-2 text-sm text-slate-900 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-500";

export function CatalogToolbar() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [isPending, startTransition] = useTransition();
  const [q, setQ] = useState(searchParams.get("q") ?? "");

  const filters = useMemo(
    () => ({
      cardType: searchParams.get("card_type") ?? "",
      kind: searchParams.get("kind") ?? "",
      minPrice: searchParams.get("min_price") ?? "",
      maxPrice: searchParams.get("max_price") ?? "",
      tag: searchParams.get("tag") ?? "",
      sortBy: searchParams.get("sort_by") ?? "created_at",
      sortOrder: searchParams.get("sort_order") ?? "desc",
      limit: searchParams.get("limit") ?? "12"
    }),
    [searchParams]
  );

  function updateSearch(next: Record<string, string>) {
    const params = new URLSearchParams(searchParams.toString());
    Object.entries(next).forEach(([key, value]) => setQueryParam(params, key, value));
    params.set("offset", "0");

    startTransition(() => {
      router.replace(`${pathname}?${params.toString()}`);
    });
  }

  return (
    <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.2)]">
      <CardContent className="space-y-5 pt-6">
        <form
          className="flex flex-col gap-3 md:flex-row"
          onSubmit={(event) => {
            event.preventDefault();
            updateSearch({ q });
          }}
        >
          <div className="relative flex-1">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
            <Input value={q} onChange={(event) => setQ(event.target.value)} className="rounded-2xl border-slate-300 pl-9" placeholder="Поиск по названию, описанию или тегу" />
          </div>
          <Button type="submit" size="lg" className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800" disabled={isPending}>
            Искать
          </Button>
        </form>

        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          <div className="space-y-2">
            <Label htmlFor="card_type">Тип карточки</Label>
            <select id="card_type" value={filters.cardType} onChange={(event) => updateSearch({ card_type: event.target.value })} className={selectClassName}>
              <option value="">Все</option>
              <option value="offer">Предложения</option>
              <option value="request">Запросы</option>
            </select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="kind">Формат</Label>
            <select id="kind" value={filters.kind} onChange={(event) => updateSearch({ kind: event.target.value })} className={selectClassName}>
              <option value="">Все</option>
              <option value="service">Услуга</option>
              <option value="product">Продукт</option>
            </select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="tag">Тег</Label>
            <Input id="tag" defaultValue={filters.tag} className="rounded-2xl border-slate-300" placeholder="mixing" onBlur={(event) => updateSearch({ tag: event.target.value })} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="sort">Сортировка</Label>
            <div className="grid grid-cols-2 gap-2">
              <select id="sort" value={filters.sortBy} onChange={(event) => updateSearch({ sort_by: event.target.value })} className={selectClassName}>
                <option value="created_at">По дате</option>
                <option value="price">По цене</option>
              </select>
              <select value={filters.sortOrder} onChange={(event) => updateSearch({ sort_order: event.target.value })} className={selectClassName}>
                <option value="desc">По убыванию</option>
                <option value="asc">По возрастанию</option>
              </select>
            </div>
          </div>
        </div>

        <div className="grid gap-4 md:grid-cols-[1fr_1fr_220px]">
          <div className="space-y-2">
            <Label htmlFor="min_price">Цена от</Label>
            <Input id="min_price" type="number" defaultValue={filters.minPrice} className="rounded-2xl border-slate-300" onBlur={(event) => updateSearch({ min_price: event.target.value })} placeholder="0" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="max_price">Цена до</Label>
            <Input id="max_price" type="number" defaultValue={filters.maxPrice} className="rounded-2xl border-slate-300" onBlur={(event) => updateSearch({ max_price: event.target.value })} placeholder="50000" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="limit">Карточек на страницу</Label>
            <select id="limit" value={filters.limit} onChange={(event) => updateSearch({ limit: event.target.value })} className={selectClassName}>
              <option value="12">12</option>
              <option value="24">24</option>
              <option value="36">36</option>
            </select>
          </div>
        </div>

        <div className="flex items-center gap-2 text-sm text-slate-500">
          <SlidersHorizontal className="h-4 w-4" />
          Поиск, фильтры, сортировка и пагинация уже работают на реальном backend API.
        </div>
      </CardContent>
    </Card>
  );
}
