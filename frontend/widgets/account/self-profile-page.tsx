"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";

import { cardsApi } from "@/entities/card/api/cards";
import { CardPreview } from "@/entities/card/ui/card-preview";
import { profilesApi } from "@/entities/profile/api/profiles";
import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";
import { UserAvatar } from "@/shared/ui/user-avatar";

function createCardHref(role?: string | null) {
  if (role === "customer") return "/cards/new?type=request";
  if (role === "engineer") return "/cards/new?type=offer";
  return "/cards/new";
}

function formatRole(role?: string | null) {
  if (role === "customer") return "Заказчик";
  if (role === "engineer") return "Исполнитель";
  if (role === "admin") return "Администратор";
  return "Гость";
}

export function SelfProfilePage() {
  const user = useAuthStore((state) => state.user);

  const profileQuery = useQuery({
    queryKey: ["profile", "me"],
    queryFn: () => profilesApi.me(),
    enabled: Boolean(user)
  });

  const cardsQuery = useQuery({
    queryKey: ["cards", "mine", user?.id, "profile-page"],
    queryFn: () =>
      cardsApi.list({
        author_id: user?.id ?? "",
        limit: "12",
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
      <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
        <CardHeader className="space-y-6">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div className="flex items-center gap-4">
              <UserAvatar
                avatarUrl={profile?.avatar_url}
                name={profile?.display_name}
                email={user?.email}
                className="h-24 w-24 rounded-[2rem] text-2xl"
              />
              <div className="space-y-2">
                <Badge className="bg-slate-950 text-white" variant="secondary">
                  Мой профиль
                </Badge>
                <CardTitle className="text-3xl text-slate-950">{profile?.display_name || user?.email || "Профиль"}</CardTitle>
                <p className="text-sm text-slate-500">{formatRole(user?.role)}</p>
              </div>
            </div>
            <div className="flex flex-wrap gap-2">
              <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                <Link href={`/profiles/${user?.id ?? ""}`}>Публичный профиль</Link>
              </Button>
              <Button asChild className="rounded-2xl bg-slate-950 text-white hover:bg-slate-800">
                <Link href="/profile/edit">Редактировать профиль</Link>
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="grid gap-6 lg:grid-cols-[minmax(0,2fr)_minmax(260px,1fr)]">
          <div className="space-y-3">
            <p className="text-sm font-medium text-slate-500">О себе</p>
            <p className="leading-7 text-slate-700">
              {profile?.bio || "Добавьте короткое описание, чтобы заказчики и исполнители быстрее понимали ваш стиль работы, опыт и специализацию."}
            </p>
          </div>
          <div className="grid gap-4 sm:grid-cols-3 lg:grid-cols-1">
            <div>
              <p className="text-sm text-slate-500">Рейтинг</p>
              <p className="text-lg font-semibold text-slate-950">{profile ? profile.rating.toFixed(1) : "0.0"}</p>
            </div>
            <div>
              <p className="text-sm text-slate-500">Отзывы</p>
              <p className="text-lg font-semibold text-slate-950">{profile?.reviews_count ?? 0}</p>
            </div>
            <div>
              <p className="text-sm text-slate-500">ID профиля</p>
              <p className="font-mono text-sm text-slate-700">{profile?.user_id ?? user?.id ?? "—"}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <section className="space-y-4">
        <div className="flex items-center justify-between gap-4">
          <div>
            <h2 className="text-2xl font-semibold text-slate-950">Мои карточки</h2>
            <p className="text-sm text-slate-500">Управляйте своими предложениями и запросами из одного места.</p>
          </div>
          {canCreateCards ? (
            <Button asChild className="rounded-2xl bg-slate-950 text-white hover:bg-slate-800">
              <Link href={createCardHref(user?.role)}>Создать карточку</Link>
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
          <Card className="border-slate-200/80 bg-white/95">
            <CardContent className="flex flex-col items-start gap-4 pt-6">
              <p className="text-slate-700">Пока нет публичных карточек.</p>
              {canCreateCards ? (
                <Button asChild className="rounded-2xl bg-slate-950 text-white hover:bg-slate-800">
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
