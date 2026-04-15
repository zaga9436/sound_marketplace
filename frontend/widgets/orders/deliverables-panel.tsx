"use client";

import { useMutation, useQuery } from "@tanstack/react-query";
import { Download, FileAudio, History } from "lucide-react";

import { cardsApi } from "@/entities/card/api/cards";
import { deliverablesApi } from "@/entities/deliverable/api/deliverables";
import { UploadDeliverableForm } from "@/features/deliverable/upload-deliverable-form";
import { getErrorMessage } from "@/lib/api/errors";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { Card as MarketplaceCard, Deliverable, Order, User } from "@/shared/types/api";

function formatBytes(size: number) {
  if (!size) return "0 Б";
  const units = ["Б", "КБ", "МБ", "ГБ"];
  let value = size;
  let unitIndex = 0;

  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024;
    unitIndex += 1;
  }

  return `${value.toFixed(value >= 10 || unitIndex === 0 ? 0 : 1)} ${units[unitIndex]}`;
}

function canUploadDeliverable(order: Order, user?: User | null) {
  if (!user) return false;
  if (user.role === "admin") return true;
  if (user.role !== "engineer" || user.id !== order.engineer_id) return false;
  return order.status === "in_progress" || order.status === "review" || order.status === "dispute";
}

export function DeliverablesPanel({ order, user, sourceCard }: { order: Order; user?: User | null; sourceCard?: MarketplaceCard }) {
  const deliverablesQuery = useQuery({
    queryKey: ["deliverables", order.id],
    queryFn: () => deliverablesApi.list(order.id),
    refetchInterval: 5000
  });

  const downloadMutation = useMutation({
    mutationFn: async (deliverable: Deliverable) => {
      const { url } = await deliverablesApi.getDownloadUrl(order.id, deliverable.id);
      window.open(url, "_blank", "noopener,noreferrer");
    }
  });

  const readyProduct = sourceCard?.card_type === "offer" && sourceCard.kind === "product";
  const canAccessReadyProduct = Boolean(readyProduct && ["in_progress", "review", "completed"].includes(order.status));
  const cardFullDownloadMutation = useMutation({
    mutationFn: async () => {
      const { url } = await cardsApi.getFullDownloadUrl(sourceCard!.id);
      window.open(url, "_blank", "noopener,noreferrer");
    }
  });

  const canUpload = canUploadDeliverable(order, user) && !readyProduct;
  const deliverables = deliverablesQuery.data ?? [];

  return (
    <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
      <CardHeader className="space-y-3">
        <div className="flex items-center justify-between gap-3">
          <div>
            <CardTitle>Результаты по заказу</CardTitle>
            <CardDescription>
              {readyProduct
                ? "Для готового продукта здесь доступен full-файл, который инженер заранее приложил к карточке."
                : "Здесь хранятся версии результата, которые исполнитель загружает по этой сделке."}
            </CardDescription>
          </div>
          <Badge variant="outline">{readyProduct ? "Готовый файл" : `${deliverables.length} версий`}</Badge>
        </div>
      </CardHeader>

      <CardContent className="space-y-5">
        {canUpload ? (
          <div className="rounded-2xl border border-slate-200 bg-slate-50 p-4">
            <div className="mb-3 flex items-center gap-2 text-sm font-medium text-slate-900">
              <History className="h-4 w-4" />
              Загрузка новой версии
            </div>
            <UploadDeliverableForm orderId={order.id} />
          </div>
        ) : null}

        {readyProduct ? (
          <div className="rounded-2xl border border-slate-200 bg-slate-50 p-4">
            <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
              <div className="min-w-0 space-y-1">
                <p className="font-medium text-slate-950">Полный файл готового продукта</p>
                <p className="text-sm leading-6 text-slate-600">
                  Это приватный full-файл из карточки. Он открывается участникам после того, как заказ взят в работу.
                </p>
              </div>
              <Button
                type="button"
                variant="outline"
                className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100"
                disabled={!canAccessReadyProduct || cardFullDownloadMutation.isPending}
                onClick={() => cardFullDownloadMutation.mutate()}
              >
                <Download className="mr-2 h-4 w-4" />
                {cardFullDownloadMutation.isPending ? "Готовим ссылку..." : canAccessReadyProduct ? "Открыть full-файл" : "Доступ после старта"}
              </Button>
            </div>
            {cardFullDownloadMutation.isError ? <p className="mt-3 text-sm text-red-600">{getErrorMessage(cardFullDownloadMutation.error)}</p> : null}
          </div>
        ) : null}

        {deliverablesQuery.isLoading ? (
          <div className="space-y-3">
            {Array.from({ length: 2 }).map((_, index) => (
              <div key={index} className="surface h-24 animate-pulse bg-slate-100/80" />
            ))}
          </div>
        ) : deliverablesQuery.isError ? (
          <p className="text-sm text-red-600">{getErrorMessage(deliverablesQuery.error)}</p>
        ) : deliverables.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-slate-300 bg-slate-50 p-5 text-sm leading-6 text-slate-600">
            {readyProduct
              ? "Для готового продукта отдельные deliverables обычно не нужны: итоговый файл уже приложен к карточке."
              : "Пока результаты не загружены. Как только исполнитель добавит файл, он появится здесь со своей версией и статусом."}
          </div>
        ) : (
          <div className="space-y-3">
            {deliverables.map((deliverable) => (
              <div key={deliverable.id} className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
                <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
                  <div className="min-w-0 space-y-2">
                    <div className="flex flex-wrap items-center gap-2">
                      <Badge className={deliverable.is_active ? "bg-slate-900 text-white" : ""} variant={deliverable.is_active ? "secondary" : "outline"}>
                        Версия {deliverable.version}
                      </Badge>
                      {deliverable.is_active ? <Badge variant="outline">Актуальная</Badge> : <Badge variant="outline">Архив</Badge>}
                    </div>
                    <div className="flex items-start gap-3">
                      <div className="mt-0.5 flex h-10 w-10 items-center justify-center rounded-xl bg-slate-100 text-slate-700">
                        <FileAudio className="h-4 w-4" />
                      </div>
                      <div className="min-w-0">
                        <p className="truncate font-medium text-slate-950">{deliverable.original_filename}</p>
                        <p className="text-sm text-slate-500">
                          {formatBytes(deliverable.size_bytes)} • {new Date(deliverable.created_at).toLocaleString("ru-RU")}
                        </p>
                      </div>
                    </div>
                  </div>

                  <Button
                    type="button"
                    variant="outline"
                    className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100"
                    disabled={downloadMutation.isPending}
                    onClick={() => downloadMutation.mutate(deliverable)}
                  >
                    <Download className="mr-2 h-4 w-4" />
                    {downloadMutation.isPending ? "Готовим ссылку..." : "Скачать"}
                  </Button>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
