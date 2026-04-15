"use client";

import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ArrowUpRight, RefreshCcw, Wallet } from "lucide-react";

import { paymentsApi } from "@/entities/payment/api/payments";
import { getErrorMessage } from "@/lib/api/errors";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";

function formatPrice(value: number) {
  return new Intl.NumberFormat("ru-RU", {
    style: "currency",
    currency: "RUB",
    maximumFractionDigits: 0
  }).format(value);
}

const quickAmounts = [500, 1000, 3000, 5000];

export function BalancePage() {
  const queryClient = useQueryClient();
  const [amount, setAmount] = useState("1000");
  const [externalId, setExternalId] = useState("");

  const balanceQuery = useQuery({
    queryKey: ["balance"],
    queryFn: () => paymentsApi.getBalance(),
    refetchInterval: 5000,
    staleTime: 0,
    refetchOnMount: "always",
    refetchOnWindowFocus: true
  });

  const createPaymentMutation = useMutation({
    mutationFn: (depositAmount: number) => paymentsApi.createDeposit(depositAmount)
  });

  const syncMutation = useMutation({
    mutationFn: (id: string) => paymentsApi.sync(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["balance"] });
      await queryClient.invalidateQueries({ queryKey: ["notifications"] });
    }
  });

  const confirmationUrl = createPaymentMutation.data?.confirmation_url || createPaymentMutation.data?.redirect_url;
  const syncHint = useMemo(() => {
    if (!syncMutation.data) return null;
    return syncMutation.data.deposit_created
      ? "Платеж подтвержден, сумма зачислена на баланс."
      : "Статус платежа обновлен. Повторного зачисления не было.";
  }, [syncMutation.data]);

  const handleCreatePayment = () => {
    const parsed = Number(amount);
    if (!Number.isFinite(parsed) || parsed <= 0) return;
    createPaymentMutation.mutate(parsed);
  };

  const handleSync = () => {
    if (!externalId.trim()) return;
    syncMutation.mutate(externalId.trim());
  };

  return (
    <div className="space-y-8">
      <section className="space-y-3">
        <Badge className="bg-slate-900/90 text-white" variant="secondary">
          Баланс и платежи
        </Badge>
        <div className="space-y-2">
          <h1 className="text-3xl font-semibold tracking-tight text-slate-950">Средства внутри SoundMarket</h1>
          <p className="max-w-3xl text-base leading-7 text-slate-600">
            Пополняйте баланс, оплачивайте заказы и проверяйте поступление платежа после возврата со страницы оплаты.
          </p>
        </div>
      </section>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.2fr)_420px]">
        <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
          <CardHeader className="space-y-3">
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-slate-900 text-white">
                <Wallet className="h-5 w-5" />
              </div>
              <div>
                <CardTitle>Текущий баланс</CardTitle>
                <CardDescription>Доступные средства для покупки готовых продуктов и создания заказов.</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-6">
            {balanceQuery.isLoading ? (
              <div className="surface h-28 animate-pulse bg-slate-100/80" />
            ) : balanceQuery.isError ? (
              <p className="text-sm text-red-600">{getErrorMessage(balanceQuery.error)}</p>
            ) : (
              <div className="rounded-[1.75rem] border border-slate-200 bg-[linear-gradient(145deg,rgba(15,23,42,0.96),rgba(51,65,85,0.92))] p-6 text-white">
                <p className="text-sm text-white/70">Доступно сейчас</p>
                <p className="mt-3 text-4xl font-semibold tracking-tight text-white">{formatPrice(balanceQuery.data?.balance ?? 0)}</p>
                <p className="mt-4 text-sm text-white/75">Баланс используется для резервирования средств по заказам и возвратов при отмене сделки.</p>
              </div>
            )}

            <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
              {quickAmounts.map((value) => (
                <button
                  key={value}
                  type="button"
                  onClick={() => setAmount(String(value))}
                  className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-left text-sm font-medium text-slate-900 transition hover:bg-slate-100"
                >
                  {formatPrice(value)}
                </button>
              ))}
            </div>
          </CardContent>
        </Card>

        <div className="space-y-6">
          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader>
              <CardTitle>Пополнение баланса</CardTitle>
              <CardDescription>Создайте платеж и перейдите на страницу оплаты.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="deposit-amount">Сумма, ₽</Label>
                <Input id="deposit-amount" type="number" min="1" step="1" value={amount} onChange={(event) => setAmount(event.target.value)} className="rounded-2xl border-slate-300" />
              </div>

              {createPaymentMutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(createPaymentMutation.error)}</p> : null}
              {createPaymentMutation.isSuccess ? (
                <div className="rounded-2xl border border-emerald-200 bg-emerald-50 p-4 text-sm leading-6 text-emerald-800">
                  Платеж создан. После оплаты вернитесь на эту страницу и нажмите “Проверить платеж”, если баланс не обновился автоматически.
                </div>
              ) : null}

              <div className="flex flex-wrap gap-3">
                <Button type="button" className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800" disabled={createPaymentMutation.isPending} onClick={handleCreatePayment}>
                  {createPaymentMutation.isPending ? "Создаем платеж..." : "Создать платеж"}
                </Button>

                {confirmationUrl ? (
                  <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                    <a href={confirmationUrl} target="_blank" rel="noreferrer">
                      <ArrowUpRight className="mr-2 h-4 w-4" />
                      Перейти к оплате
                    </a>
                  </Button>
                ) : null}
              </div>

              {createPaymentMutation.data?.external_id ? (
                <div className="rounded-2xl border border-slate-200 bg-slate-50 p-4 text-sm text-slate-600">
                  Номер платежа: <span className="font-mono text-slate-900">{createPaymentMutation.data.external_id}</span>
                </div>
              ) : null}
            </CardContent>
          </Card>

          <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
            <CardHeader>
              <CardTitle>Проверка платежа</CardTitle>
              <CardDescription>Если после оплаты баланс еще не обновился, введите номер платежа и проверьте его статус.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="external-id">Номер платежа</Label>
                <Input
                  id="external-id"
                  value={externalId}
                  onChange={(event) => setExternalId(event.target.value)}
                  placeholder="Например, 2f1c2d7e-..."
                  className="rounded-2xl border-slate-300"
                />
              </div>

              {syncMutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(syncMutation.error)}</p> : null}
              {syncHint ? <p className="text-sm text-emerald-700">{syncHint}</p> : null}

              <Button type="button" variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100" disabled={syncMutation.isPending} onClick={handleSync}>
                <RefreshCcw className="mr-2 h-4 w-4" />
                {syncMutation.isPending ? "Проверяем..." : "Проверить платеж"}
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
