"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";

import { notificationsApi } from "@/entities/notification/api/notifications";
import { MarkNotificationsReadButton } from "@/features/notifications/mark-read-button";
import { getErrorMessage } from "@/lib/api/errors";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";
import { NotificationItem } from "@/shared/types/api";

const statusLabels: Record<string, string> = {
  created: "создан",
  on_hold: "средства зарезервированы",
  in_progress: "в работе",
  review: "на проверке",
  completed: "завершен",
  dispute: "открыт спор",
  cancelled: "отменен",
  open: "открыт",
  closed: "закрыт",
  succeeded: "успешно оплачен",
  pending: "ожидает оплаты"
};

const phraseLabels: Array<[RegExp, string]> = [
  [/New message in order chat/gi, "Новое сообщение в чате заказа"],
  [/Card published/gi, "Карточка опубликована"],
  [/Card hidden/gi, "Карточка скрыта"],
  [/Card unhidden/gi, "Карточка возвращена в каталог"],
  [/New deliverable version uploaded/gi, "Загружена новая версия результата"],
  [/New deliverable uploaded/gi, "Загружен результат по заказу"],
  [/Payment succeeded/gi, "Платеж успешно зачислен"],
  [/Balance replenished/gi, "Баланс успешно пополнен"],
  [/Order status changed/gi, "Статус заказа изменен"],
  [/Order created/gi, "Создан новый заказ"],
  [/New bid/gi, "Новый отклик"],
  [/Dispute opened/gi, "Открыт спор"],
  [/Dispute closed/gi, "Спор закрыт"],
  [/Review received/gi, "Получен новый отзыв"]
];

function localizeText(value: string) {
  let result = value;
  phraseLabels.forEach(([pattern, replacement]) => {
    result = result.replace(pattern, replacement);
  });
  Object.entries(statusLabels).forEach(([raw, label]) => {
    result = result.replace(new RegExp(`\\b${raw}\\b`, "gi"), label);
  });
  return result;
}

function getNotificationLink(item: NotificationItem) {
  const match = item.message.match(/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/i);
  const id = match?.[0];
  if (!id) return null;

  if (item.type.includes("order") || item.type.includes("message") || item.type.includes("deliverable") || item.type.includes("dispute") || item.type.includes("review")) {
    return `/orders/${id}`;
  }
  if (item.type.includes("bid") || item.type.includes("card")) return `/cards/${id}`;
  return null;
}

function typeLabel(type: string) {
  const labels: Record<string, string> = {
    new_bid: "Новый отклик",
    order_created: "Новый заказ",
    order_status_changed: "Статус заказа",
    dispute_opened: "Открыт спор",
    dispute_closed: "Спор закрыт",
    new_message: "Новое сообщение",
    review_received: "Новый отзыв",
    payment_succeeded: "Платеж зачислен",
    deliverable_uploaded: "Загружен результат",
    deliverable_updated: "Обновлен результат",
    card_published: "Карточка опубликована"
  };
  return labels[type] ?? localizeText(type.replaceAll("_", " "));
}

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
            <h1>Центр событий аккаунта</h1>
            <p>Здесь собраны отклики, сообщения, изменения заказов, споры, отзывы и платежи.</p>
          </div>
          <MarkNotificationsReadButton />
        </div>
      </section>

      {items.length === 0 ? (
        <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
          <CardContent className="pt-6">
            <p>Пока уведомлений нет. Когда в заказах или карточках появится активность, она отобразится здесь.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {items.map((item) => {
            const href = getNotificationLink(item);
            return (
              <Card key={item.id} className={`border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.16)] ${item.is_read ? "" : "ring-1 ring-slate-300"}`}>
                <CardHeader className="flex flex-row items-start justify-between gap-4">
                  <div className="space-y-2">
                    <div className="flex flex-wrap items-center gap-2">
                      <Badge variant="outline">{typeLabel(item.type)}</Badge>
                      <Badge variant={item.is_read ? "outline" : "secondary"}>{item.is_read ? "Прочитано" : "Новое"}</Badge>
                    </div>
                    <CardTitle className="text-lg">{localizeText(item.message)}</CardTitle>
                    <p className="text-sm text-slate-500">{new Date(item.created_at).toLocaleString("ru-RU")}</p>
                  </div>
                  <div className="flex flex-col items-end gap-2">
                    {href ? (
                      <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                        <Link href={href}>Открыть</Link>
                      </Button>
                    ) : null}
                    {!item.is_read ? <MarkNotificationsReadButton ids={[item.id]} /> : null}
                  </div>
                </CardHeader>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}
