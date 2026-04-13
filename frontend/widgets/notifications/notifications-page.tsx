"use client";

import { useQuery } from "@tanstack/react-query";

import { notificationsApi } from "@/entities/notification/api/notifications";
import { MarkNotificationsReadButton } from "@/features/notifications/mark-read-button";
import { getErrorMessage } from "@/lib/api/errors";
import { Badge } from "@/shared/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";

export function NotificationsPage() {
  const query = useQuery({
    queryKey: ["notifications"],
    queryFn: () => notificationsApi.list(50),
    refetchInterval: 5000
  });

  if (query.isLoading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 4 }).map((_, index) => (
          <div key={index} className="surface h-24 animate-pulse bg-white/70" />
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

  const data = query.data;
  const items = data?.items ?? [];

  return (
    <div className="space-y-8">
      <section className="space-y-3">
        <div className="flex flex-wrap items-center gap-3">
          <Badge className="bg-slate-900/90 text-white" variant="secondary">
            Уведомления
          </Badge>
          <Badge variant="outline">Непрочитано: {data?.unread_count ?? 0}</Badge>
        </div>
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div className="space-y-2">
            <h1>Центр событий по вашему аккаунту</h1>
            <p>Здесь собраны новые отклики, сообщения, изменения статусов заказов и другие важные события в SoundMarket.</p>
          </div>
          <MarkNotificationsReadButton />
        </div>
      </section>

      {items.length === 0 ? (
        <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
          <CardContent className="pt-6">
            <p>Пока уведомлений нет. Как только в заказах или карточках появится активность, она отобразится здесь.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {items.map((item) => (
            <Card key={item.id} className={`border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.16)] ${item.is_read ? "" : "ring-1 ring-slate-300"}`}>
              <CardHeader className="flex flex-row items-start justify-between gap-4">
                <div className="space-y-2">
                  <CardTitle className="text-lg">{item.message}</CardTitle>
                  <p className="text-sm text-slate-500">{new Date(item.created_at).toLocaleString("ru-RU")}</p>
                </div>
                <div className="flex flex-col items-end gap-2">
                  <Badge variant={item.is_read ? "outline" : "secondary"}>{item.is_read ? "Прочитано" : "Новое"}</Badge>
                  {!item.is_read ? <MarkNotificationsReadButton ids={[item.id]} /> : null}
                </div>
              </CardHeader>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
