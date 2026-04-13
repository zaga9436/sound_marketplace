"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";

import { cardsApi } from "@/entities/card/api/cards";
import { CardPreview } from "@/entities/card/ui/card-preview";
import { profilesApi } from "@/entities/profile/api/profiles";
import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";

function createCardHref(role?: string | null) {
  if (role === "customer") return "/cards/new?type=request";
  if (role === "engineer") return "/cards/new?type=offer";
  return "/cards/new";
}

export function DashboardOverview() {
  const user = useAuthStore((state) => state.user);

  const profileQuery = useQuery({
    queryKey: ["profile", "me"],
    queryFn: () => profilesApi.me(),
    enabled: Boolean(user)
  });

  const cardsQuery = useQuery({
    queryKey: ["cards", "mine", user?.id],
    queryFn: () =>
      cardsApi.list({
        author_id: user?.id ?? "",
        limit: "6",
        offset: "0",
        sort_by: "created_at",
        sort_order: "desc"
      }),
    enabled: Boolean(user?.id)
  });

  const profile = profileQuery.data;
  const cards = cardsQuery.data?.items ?? [];
  const canCreateCards = user?.role === "customer" || user?.role === "engineer";

  return (
    <div className="space-y-8">
      <section className="grid gap-4 lg:grid-cols-[minmax(0,1.5fr)_minmax(320px,1fr)]">
        <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
          <CardHeader className="space-y-4">
            <Badge className="bg-slate-900/90 text-white" variant="secondary">
              Личный кабинет
            </Badge>
            <div className="space-y-2">
              <CardTitle className="text-3xl text-slate-950">С возвращением, {profile?.display_name || user?.email || "пользователь"}.</CardTitle>
              <CardDescription className="max-w-2xl text-base leading-7 text-slate-600">
                Отсюда удобно обновлять профиль, создавать карточки и переходить к ключевым сценариям SoundMarket: откликам, заказам и сделкам.
              </CardDescription>
            </div>
          </CardHeader>
          <CardContent className="flex flex-wrap gap-3">
            <Button asChild size="lg" className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800">
              <Link href="/profile">Мой профиль</Link>
            </Button>
            <Button asChild variant="outline" size="lg" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
              <Link href="/profile/edit">Редактировать профиль</Link>
            </Button>
            <Button asChild variant="outline" size="lg" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
              <Link href="/orders">Мои заказы</Link>
            </Button>
            {canCreateCards ? (
              <Button asChild variant="secondary" size="lg" className="rounded-2xl bg-slate-100 text-slate-900 hover:bg-slate-200">
                <Link href={createCardHref(user?.role)}>{user?.role === "customer" ? "Создать запрос" : "Создать предложение"}</Link>
              </Button>
            ) : null}
          </CardContent>
        </Card>

        <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
          <CardHeader>
            <CardTitle className="text-xl">Сводка</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-4 sm:grid-cols-2 lg:grid-cols-1">
            <div>
              <p className="text-sm text-slate-500">Роль</p>
              <p className="text-lg font-semibold capitalize text-slate-950">{user?.role ?? "guest"}</p>
            </div>
            <div>
              <p className="text-sm text-slate-500">Рейтинг</p>
              <p className="text-lg font-semibold text-slate-950">{profile ? profile.rating.toFixed(1) : "0.0"}</p>
            </div>
            <div>
              <p className="text-sm text-slate-500">Отзывы</p>
              <p className="text-lg font-semibold text-slate-950">{profile?.reviews_count ?? 0}</p>
            </div>
            <div>
              <p className="text-sm text-slate-500">Публичные карточки</p>
              <p className="text-lg font-semibold text-slate-950">{cardsQuery.data?.total ?? 0}</p>
            </div>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-4">
        <div className="flex items-center justify-between gap-4">
          <div>
            <h2>Ваши последние карточки</h2>
            <p className="text-sm text-slate-500">Быстрый доступ к карточкам, которые уже опубликованы в маркетплейсе.</p>
          </div>
          {canCreateCards ? (
            <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
              <Link href={createCardHref(user?.role)}>Новая карточка</Link>
            </Button>
          ) : null}
        </div>

        {cardsQuery.isLoading ? (
          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            {Array.from({ length: 3 }).map((_, index) => (
              <div key={index} className="surface h-[280px] animate-pulse bg-white/70" />
            ))}
          </div>
        ) : cards.length === 0 ? (
          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardContent className="flex flex-col items-start gap-4 pt-6">
              <p>У вас пока нет публичных карточек.</p>
              {canCreateCards ? (
                <Button asChild className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800">
                  <Link href={createCardHref(user?.role)}>Создать первую карточку</Link>
                </Button>
              ) : null}
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            {cards.map((card) => (
              <div key={card.id} className="space-y-3">
                <CardPreview card={card} />
                <Button asChild variant="outline" className="w-full rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                  <Link href={`/cards/${card.id}/edit`}>Редактировать карточку</Link>
                </Button>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
