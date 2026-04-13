"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";

import { cardsApi } from "@/entities/card/api/cards";
import { profilesApi } from "@/entities/profile/api/profiles";
import { CardPreview } from "@/entities/card/ui/card-preview";
import { getErrorMessage } from "@/lib/api/errors";
import { Badge } from "@/shared/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";

export function PublicProfile({ id }: { id: string }) {
  const profileQuery = useQuery({
    queryKey: ["profile", id],
    queryFn: () => profilesApi.getById(id)
  });

  const cardsQuery = useQuery({
    queryKey: ["profile", id, "cards"],
    queryFn: () => cardsApi.list({ author_id: id, limit: "6", offset: "0" })
  });

  const reviewsQuery = useQuery({
    queryKey: ["profile", id, "reviews"],
    queryFn: () => profilesApi.listReviews(id)
  });

  if (profileQuery.isLoading) {
    return <div className="surface h-[320px] animate-pulse bg-white/70" />;
  }

  if (profileQuery.isError) {
    return (
      <Card className="border-destructive/20 bg-white/95">
        <CardContent className="pt-6">
          <p className="text-destructive">{getErrorMessage(profileQuery.error)}</p>
        </CardContent>
      </Card>
    );
  }

  const profile = profileQuery.data;
  if (!profile) return null;

  const cards = cardsQuery.data?.items ?? [];
  const reviews = reviewsQuery.data ?? [];

  return (
    <div className="space-y-8">
      <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
        <CardHeader className="space-y-4">
          <Badge className="bg-slate-900/90 text-white" variant="secondary">
            Публичный профиль
          </Badge>
          <div className="space-y-2">
            <CardTitle className="text-3xl text-slate-950">{profile.display_name}</CardTitle>
            <p>{profile.bio || "Автор пока не добавил описание профиля."}</p>
          </div>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-6">
          <div>
            <p className="text-sm text-slate-500">Рейтинг</p>
            <p className="text-xl font-semibold text-slate-950">{profile.rating.toFixed(1)}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">Отзывы</p>
            <p className="text-xl font-semibold text-slate-950">{profile.reviews_count}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">ID автора</p>
            <p className="font-mono text-sm text-slate-700">{profile.user_id}</p>
          </div>
        </CardContent>
      </Card>

      <section className="space-y-4">
        <div className="flex items-center justify-between">
          <h2>Карточки автора</h2>
          <Link href="/catalog" className="text-sm font-medium text-primary">
            Вернуться в каталог
          </Link>
        </div>

        {cards.length === 0 ? (
          <Card className="border-slate-200/80 bg-white/95">
            <CardContent className="pt-6">
              <p>У автора пока нет публичных карточек.</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            {cards.map((card) => (
              <CardPreview key={card.id} card={card} />
            ))}
          </div>
        )}
      </section>

      <section className="space-y-4">
        <div className="flex items-center justify-between">
          <h2>Последние отзывы</h2>
          <Badge variant="secondary">{reviews.length} показано</Badge>
        </div>

        {reviews.length === 0 ? (
          <Card className="border-slate-200/80 bg-white/95">
            <CardContent className="pt-6">
              <p>Пока отзывов нет.</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2">
            {reviews.slice(0, 4).map((review) => (
              <Card key={review.id} className="border-slate-200/80 bg-white/95">
                <CardHeader className="space-y-3">
                  <div className="flex items-center justify-between gap-3">
                    <CardTitle className="text-lg">Оценка {review.rating}/5</CardTitle>
                    <Badge variant="secondary">{new Date(review.created_at).toLocaleDateString("ru-RU")}</Badge>
                  </div>
                  <CardDescription className="font-mono text-xs">Заказ {review.order_id}</CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-sm leading-6 text-foreground/80">{review.text || "Проверенный отзыв."}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
