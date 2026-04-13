"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";

import { ordersApi } from "@/entities/order/api/orders";
import { getErrorMessage } from "@/lib/api/errors";
import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { OrderStatusBadge } from "@/widgets/orders/order-status-badge";

function formatPrice(value: number) {
  return new Intl.NumberFormat("ru-RU", {
    style: "currency",
    currency: "RUB",
    maximumFractionDigits: 0
  }).format(value);
}

export function OrdersPage() {
  const user = useAuthStore((state) => state.user);
  const query = useQuery({
    queryKey: ["orders"],
    queryFn: () => ordersApi.list()
  });

  const title = user?.role === "engineer" ? "Мои сделки и входящие заказы" : user?.role === "customer" ? "Мои покупки и заказы" : "Заказы";

  if (query.isLoading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 3 }).map((_, index) => (
          <div key={index} className="surface h-40 animate-pulse bg-white/70" />
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

  const orders = query.data ?? [];

  return (
    <div className="space-y-8">
      <section className="space-y-3">
        <Badge className="bg-slate-900/90 text-white" variant="secondary">
          Сделки
        </Badge>
        <div className="space-y-2">
          <h1>{title}</h1>
          <p>Здесь собраны все заказы, связанные с вашей ролью в SoundMarket. Отсюда удобно переходить в детали сделки и выполнять доступные действия.</p>
        </div>
      </section>

      {orders.length === 0 ? (
        <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
          <CardContent className="flex flex-col items-start gap-4 pt-6">
            <p>У вас пока нет заказов. Начните с каталога или со своих карточек.</p>
            <div className="flex flex-wrap gap-3">
              <Button asChild className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800">
                <Link href="/catalog">Перейти в каталог</Link>
              </Button>
              <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                <Link href="/dashboard">Открыть кабинет</Link>
              </Button>
            </div>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {orders.map((order) => (
            <Card key={order.id} className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
              <CardHeader className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
                <div className="space-y-2">
                  <div className="flex flex-wrap items-center gap-2">
                    <OrderStatusBadge status={order.status} />
                    <Badge variant="outline">ID: {order.id}</Badge>
                  </div>
                  <CardTitle className="text-2xl">Заказ на {formatPrice(order.amount)}</CardTitle>
                  <CardDescription className="text-sm text-slate-500">
                    Создан {new Date(order.created_at).toLocaleDateString("ru-RU")} • customer {order.customer_id} • engineer {order.engineer_id}
                  </CardDescription>
                </div>
                <Button asChild className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800">
                  <Link href={`/orders/${order.id}`}>Открыть заказ</Link>
                </Button>
              </CardHeader>
              <CardContent className="grid gap-4 md:grid-cols-3">
                <div>
                  <p className="text-sm text-slate-500">Источник</p>
                  <p className="text-sm font-medium text-slate-900">{order.card_id ? `Offer ${order.card_id}` : order.request_id ? `Request ${order.request_id}` : "Сделка"}</p>
                </div>
                <div>
                  <p className="text-sm text-slate-500">Последнее обновление</p>
                  <p className="text-sm font-medium text-slate-900">{new Date(order.last_status_time).toLocaleString("ru-RU")}</p>
                </div>
                <div>
                  <p className="text-sm text-slate-500">Сумма</p>
                  <p className="text-sm font-medium text-slate-900">{formatPrice(order.amount)}</p>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
