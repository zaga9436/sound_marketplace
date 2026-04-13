"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";

import { cardsApi } from "@/entities/card/api/cards";
import { disputesApi } from "@/entities/dispute/api/disputes";
import { ordersApi } from "@/entities/order/api/orders";
import { paymentsApi } from "@/entities/payment/api/payments";
import { profilesApi } from "@/entities/profile/api/profiles";
import { DisputeActions } from "@/features/dispute/dispute-actions";
import { OrderStatusActions } from "@/features/order/order-status-actions";
import { ReviewForm } from "@/features/review/review-form";
import { getErrorMessage } from "@/lib/api/errors";
import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { User } from "@/shared/types/api";
import { OrderChatPanel } from "@/widgets/chat/order-chat-panel";
import { DeliverablesPanel } from "@/widgets/orders/deliverables-panel";
import { OrderStatusBadge } from "@/widgets/orders/order-status-badge";

function formatPrice(value: number) {
  return new Intl.NumberFormat("ru-RU", {
    style: "currency",
    currency: "RUB",
    maximumFractionDigits: 0
  }).format(value);
}

function getNextStep(orderStatus: string, user?: User | null) {
  if (orderStatus === "on_hold") {
    return user?.role === "engineer" ? "Возьмите заказ в работу и начните производство результата." : "Ожидайте, пока исполнитель возьмет заказ в работу.";
  }
  if (orderStatus === "in_progress") {
    return user?.role === "engineer" ? "Загрузите промежуточный или финальный deliverable и отправьте заказ на проверку." : "Следите за чатом и ожидайте новую версию результата.";
  }
  if (orderStatus === "review") {
    return user?.role === "customer" ? "Проверьте deliverables, обсудите детали в чате и подтвердите завершение, если все готово." : "Ожидайте проверки заказчиком и будьте готовы ответить в чате.";
  }
  if (orderStatus === "dispute") {
    return "Спор открыт. Все действия по заказу сейчас проходят через dispute flow и историю сообщений.";
  }
  if (orderStatus === "completed") {
    return "Сделка завершена. Можно скачать актуальный результат и, если доступно, оставить отзыв.";
  }
  if (orderStatus === "cancelled") {
    return "Заказ отменен. Средства возвращены или спор закрыт соответствующим решением.";
  }
  return "Следите за статусом заказа и связанными действиями в блоках ниже.";
}

export function OrderDetailPage({ id }: { id: string }) {
  const user = useAuthStore((state) => state.user);

  const orderQuery = useQuery({
    queryKey: ["order", id],
    queryFn: () => ordersApi.getById(id),
    refetchInterval: 5000
  });

  const balanceQuery = useQuery({
    queryKey: ["balance"],
    queryFn: () => paymentsApi.getBalance(),
    enabled: Boolean(user),
    refetchInterval: 5000
  });

  const order = orderQuery.data;
  const sourceCardId = order?.card_id || order?.request_id;

  const customerProfileQuery = useQuery({
    queryKey: ["profile", order?.customer_id],
    queryFn: () => profilesApi.getById(order!.customer_id),
    enabled: Boolean(order?.customer_id)
  });

  const engineerProfileQuery = useQuery({
    queryKey: ["profile", order?.engineer_id],
    queryFn: () => profilesApi.getById(order!.engineer_id),
    enabled: Boolean(order?.engineer_id)
  });

  const sourceCardQuery = useQuery({
    queryKey: ["card", sourceCardId, "order-detail"],
    queryFn: () => cardsApi.getById(sourceCardId!),
    enabled: Boolean(sourceCardId)
  });

  const disputeQuery = useQuery({
    queryKey: ["dispute", id],
    queryFn: () => disputesApi.getByOrderId(id),
    retry: false,
    enabled: Boolean(order)
  });

  if (orderQuery.isLoading) {
    return <div className="surface h-[420px] animate-pulse bg-white/70" />;
  }

  if (orderQuery.isError) {
    return (
      <Card className="border-destructive/20 bg-white/95">
        <CardContent className="pt-6">
          <p className="text-destructive">{getErrorMessage(orderQuery.error)}</p>
        </CardContent>
      </Card>
    );
  }

  if (!order) return null;

  const canReview = user?.role === "customer" && user.id === order.customer_id && order.status === "completed";
  const nextStep = getNextStep(order.status, user);

  return (
    <div className="space-y-8">
      <section className="space-y-4">
        <div className="flex flex-wrap items-center gap-3">
          <OrderStatusBadge status={order.status} />
          <Badge variant="outline">ID: {order.id}</Badge>
        </div>

        <div className="grid gap-4 xl:grid-cols-[minmax(0,1.3fr)_360px]">
          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader className="space-y-3">
              <CardTitle className="text-3xl text-slate-950">Заказ на {formatPrice(order.amount)}</CardTitle>
              <CardDescription className="max-w-3xl text-base leading-7 text-slate-600">
                Главный экран сделки: здесь видно текущий этап, участников, чат, deliverables, спор, отзыв и следующие действия.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-5">
              <div className="grid gap-5 md:grid-cols-2">
                <Info label="Сумма" value={formatPrice(order.amount)} />
                <Info label="Создан" value={new Date(order.created_at).toLocaleString("ru-RU")} />
                <Info label="Последнее изменение" value={new Date(order.last_status_time).toLocaleString("ru-RU")} />
                <Info label="Источник" value={order.card_id ? `Карточка ${order.card_id}` : order.request_id ? `Запрос ${order.request_id}` : "Сделка"} />
                {order.bid_id ? <Info label="Связанный отклик" value={order.bid_id} /> : null}
                {order.dispute_reason ? <Info label="Причина спора" value={order.dispute_reason} /> : null}
              </div>

              <div className="rounded-2xl border border-slate-200 bg-slate-50 p-4">
                <p className="text-sm font-medium text-slate-900">Что делать дальше</p>
                <p className="mt-2 text-sm leading-6 text-slate-600">{nextStep}</p>
              </div>
            </CardContent>
          </Card>

          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader>
              <CardTitle>Деньги по сделке</CardTitle>
              <CardDescription>Баланс помогает быстро понять, хватает ли средств на новые действия и пополнения.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {balanceQuery.isLoading ? (
                <div className="surface h-20 animate-pulse bg-slate-100/80" />
              ) : balanceQuery.isError ? (
                <p className="text-sm text-red-600">{getErrorMessage(balanceQuery.error)}</p>
              ) : (
                <div className="rounded-2xl border border-slate-200 bg-[linear-gradient(145deg,rgba(15,23,42,0.96),rgba(51,65,85,0.92))] p-5 text-white">
                  <p className="text-sm text-white/70">Текущий баланс</p>
                  <p className="mt-2 text-3xl font-semibold">{formatPrice(balanceQuery.data?.balance ?? 0)}</p>
                </div>
              )}

              <div className="flex flex-wrap gap-3">
                <Button asChild className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800">
                  <Link href="/balance">Пополнить баланс</Link>
                </Button>
                <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                  <Link href="/orders">Все заказы</Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </section>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.35fr)_360px]">
        <div className="space-y-6">
          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader>
              <CardTitle>Участники</CardTitle>
            </CardHeader>
            <CardContent className="grid gap-4 md:grid-cols-2">
              <ParticipantCard
                label="Заказчик"
                userId={order.customer_id}
                displayName={customerProfileQuery.data?.display_name}
                rating={customerProfileQuery.data?.rating}
                reviewsCount={customerProfileQuery.data?.reviews_count}
              />
              <ParticipantCard
                label="Исполнитель"
                userId={order.engineer_id}
                displayName={engineerProfileQuery.data?.display_name}
                rating={engineerProfileQuery.data?.rating}
                reviewsCount={engineerProfileQuery.data?.reviews_count}
              />
            </CardContent>
          </Card>

          <DeliverablesPanel order={order} user={user} />

          <OrderChatPanel orderId={order.id} />

          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader>
              <CardTitle>Спор по заказу</CardTitle>
              <CardDescription>Если по сделке возникла проблема, спор открывается и закрывается прямо здесь.</CardDescription>
            </CardHeader>
            <CardContent>
              {disputeQuery.isLoading ? (
                <p className="text-sm text-slate-500">Проверяем состояние спора...</p>
              ) : disputeQuery.isError ? (
                <DisputeActions order={order} user={user} />
              ) : (
                <DisputeActions order={order} dispute={disputeQuery.data} user={user} />
              )}
            </CardContent>
          </Card>

          {canReview ? (
            <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
              <CardHeader>
                <CardTitle>Отзыв после сделки</CardTitle>
                <CardDescription>После завершения заказа заказчик может оставить проверенный отзыв исполнителю.</CardDescription>
              </CardHeader>
              <CardContent>
                <ReviewForm orderId={order.id} targetUserId={order.engineer_id} />
              </CardContent>
            </Card>
          ) : null}
        </div>

        <div className="space-y-6">
          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader>
              <CardTitle>Действия по статусу</CardTitle>
              <CardDescription>Кнопки появляются только тогда, когда backend разрешает этот переход для вашей роли.</CardDescription>
            </CardHeader>
            <CardContent>
              <OrderStatusActions order={order} />
            </CardContent>
          </Card>

          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader>
              <CardTitle>Связанная карточка</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {sourceCardQuery.isLoading ? (
                <p className="text-sm text-slate-500">Загружаем карточку...</p>
              ) : sourceCardQuery.data ? (
                <>
                  <div className="space-y-2">
                    <p className="font-medium text-slate-950">{sourceCardQuery.data.title}</p>
                    <p className="text-sm text-slate-600">{sourceCardQuery.data.description}</p>
                  </div>
                  <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                    <Link href={`/cards/${sourceCardQuery.data.id}`}>Открыть карточку</Link>
                  </Button>
                </>
              ) : (
                <p className="text-sm text-slate-500">Связанная карточка недоступна.</p>
              )}
            </CardContent>
          </Card>

          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader>
              <CardTitle>Быстрые переходы</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <p className="text-sm leading-6 text-slate-600">Отсюда удобно перейти к остальным заказам, чатам, уведомлениям и балансу, не теряя контекст сделки.</p>
              <div className="flex flex-wrap gap-3">
                <Button asChild className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800">
                  <Link href="/orders">Все заказы</Link>
                </Button>
                <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                  <Link href="/chats">Чаты</Link>
                </Button>
                <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                  <Link href="/notifications">Уведомления</Link>
                </Button>
                <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                  <Link href="/balance">Баланс</Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}

function Info({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="text-sm text-slate-500">{label}</p>
      <p className="text-sm font-medium text-slate-900">{value}</p>
    </div>
  );
}

function ParticipantCard({
  label,
  userId,
  displayName,
  rating,
  reviewsCount
}: {
  label: string;
  userId: string;
  displayName?: string;
  rating?: number;
  reviewsCount?: number;
}) {
  return (
    <div className="rounded-2xl border border-slate-200 bg-slate-50/80 p-4">
      <p className="text-sm text-slate-500">{label}</p>
      <p className="mt-1 font-medium text-slate-950">{displayName ?? userId}</p>
      {rating != null ? (
        <p className="mt-1 text-sm text-slate-500">
          Рейтинг {rating.toFixed(1)} • {reviewsCount ?? 0} отзывов
        </p>
      ) : null}
    </div>
  );
}
