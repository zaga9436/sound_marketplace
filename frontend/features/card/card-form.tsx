"use client";

import Link from "next/link";
import { useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { cardsApi } from "@/entities/card/api/cards";
import { CardMediaManager } from "@/features/card/card-media-manager";
import { useAuthStore } from "@/lib/auth/session-store";
import { getErrorMessage } from "@/lib/api/errors";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Textarea } from "@/shared/ui/textarea";
import { CardType } from "@/shared/types/api";

const cardSchema = z.object({
  card_type: z.enum(["offer", "request"]),
  kind: z.enum(["product", "service"]),
  title: z.string().min(3, "Название должно быть не короче 3 символов."),
  description: z.string().min(10, "Описание должно быть не короче 10 символов."),
  price: z.coerce.number().int().min(1, "Цена должна быть больше нуля."),
  tags: z.string().optional(),
  is_published: z.boolean()
});

type CardFormValues = z.infer<typeof cardSchema>;

type CardFormProps = {
  mode: "create" | "edit";
  cardId?: string;
  initialCardType?: CardType;
};

const selectClassName =
  "flex h-11 w-full rounded-2xl border border-slate-300 bg-white px-4 py-2 text-sm text-slate-900 shadow-sm transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-500";

function tagsToString(tags: string[] | undefined) {
  return tags?.join(", ") ?? "";
}

function parseTags(value?: string) {
  return (value ?? "")
    .split(",")
    .map((tag) => tag.trim())
    .filter(Boolean);
}

function allowedCardTypeForRole(role?: string | null) {
  if (role === "customer") return "request";
  if (role === "engineer") return "offer";
  return null;
}

export function CardForm({ mode, cardId, initialCardType }: CardFormProps) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const allowedCardType = allowedCardTypeForRole(user?.role);

  const cardQuery = useQuery({
    queryKey: ["card", cardId, "edit"],
    queryFn: () => cardsApi.getById(cardId ?? ""),
    enabled: mode === "edit" && Boolean(cardId)
  });

  const form = useForm<CardFormValues>({
    resolver: zodResolver(cardSchema),
    defaultValues: {
      card_type: initialCardType ?? allowedCardType ?? "offer",
      kind: "service",
      title: "",
      description: "",
      price: 1000,
      tags: "",
      is_published: true
    }
  });

  useEffect(() => {
    if (mode === "edit" && cardQuery.data) {
      form.reset({
        card_type: cardQuery.data.card_type,
        kind: cardQuery.data.kind,
        title: cardQuery.data.title,
        description: cardQuery.data.description,
        price: cardQuery.data.price,
        tags: tagsToString(cardQuery.data.tags),
        is_published: cardQuery.data.is_published
      });
    }
  }, [cardQuery.data, form, mode]);

  useEffect(() => {
    if (mode === "create" && allowedCardType) {
      form.setValue("card_type", initialCardType ?? allowedCardType);
    }
  }, [allowedCardType, form, initialCardType, mode]);

  const mutation = useMutation({
    mutationFn: async (values: CardFormValues) => {
      const payload = {
        card_type: values.card_type,
        kind: values.kind,
        title: values.title.trim(),
        description: values.description.trim(),
        price: values.price,
        tags: parseTags(values.tags),
        is_published: values.is_published
      };

      if (mode === "edit" && cardId) {
        const { card_type: _cardType, ...updatePayload } = payload;
        return cardsApi.update(cardId, updatePayload);
      }

      return cardsApi.create(payload);
    },
    onSuccess: async (card) => {
      await queryClient.invalidateQueries({
        predicate: (query) => Array.isArray(query.queryKey) && query.queryKey.some((part) => part === "cards" || part === "card")
      });
      router.push(mode === "create" ? `/cards/${card.id}/edit?created=1` : `/cards/${card.id}`);
    }
  });

  const onSubmit = form.handleSubmit((values) => {
    if (mode === "create" && allowedCardType && values.card_type !== allowedCardType) {
      form.setError("card_type", {
        message: allowedCardType === "offer" ? "Инженер может создавать только offer-карточки." : "Заказчик может создавать только request-карточки."
      });
      return;
    }

    mutation.mutate(values);
  });

  if (mode === "edit" && cardQuery.isLoading) {
    return <div className="surface h-[520px] animate-pulse bg-white/70" />;
  }

  if (mode === "edit" && cardQuery.isError) {
    return (
      <Card className="border-destructive/20 bg-white/95">
        <CardContent className="pt-6">
          <p className="text-destructive">{getErrorMessage(cardQuery.error)}</p>
        </CardContent>
      </Card>
    );
  }

  if (mode === "create" && !allowedCardType) {
    return (
      <Card className="border-white/60 bg-white/95">
        <CardHeader>
          <Badge>Управление карточками</Badge>
          <CardTitle>Создание карточки</CardTitle>
          <CardDescription>Карточки могут создавать только заказчики и инженеры.</CardDescription>
        </CardHeader>
        <CardContent>
          <Button asChild>
            <Link href="/dashboard">Вернуться в кабинет</Link>
          </Button>
        </CardContent>
      </Card>
    );
  }

  const cardTypeValue = form.watch("card_type");
  const card = cardQuery.data;

  return (
    <div className="space-y-6">
      <Card className="overflow-hidden border-slate-200/80 bg-white/95 shadow-[0_24px_80px_-40px_rgba(15,23,42,0.38)]">
        <CardHeader className="space-y-4 border-b border-slate-200/80 bg-[linear-gradient(135deg,rgba(15,23,42,0.04),rgba(71,85,105,0.02))]">
          <div className="flex flex-wrap items-center gap-3">
            <Badge className="bg-slate-900/90 text-white hover:bg-slate-900" variant="secondary">
              {mode === "create" ? "Новая карточка" : "Редактирование"}
            </Badge>
            {mode === "edit" ? <Badge variant="outline">ID: {cardId}</Badge> : null}
          </div>
          <div className="space-y-2">
            <CardTitle className="text-3xl text-slate-950">
              {mode === "create" ? "Создайте карточку для каталога SoundMarket" : "Обновите карточку и медиа"}
            </CardTitle>
            <CardDescription className="max-w-3xl text-base leading-7 text-slate-600">
              Чистая, понятная карточка помогает быстрее показать ценность предложения или запроса. После создания карточки вы
              сразу перейдете к шагу с обложкой, preview и приватным полным файлом.
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent className="p-6">
          <form onSubmit={onSubmit} className="space-y-6">
            <div className="grid gap-6 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="card_type">Тип карточки</Label>
                <select id="card_type" className={selectClassName} disabled={mode === "edit" || Boolean(allowedCardType)} {...form.register("card_type")}>
                  <option value="offer">Offer</option>
                  <option value="request">Request</option>
                </select>
                <p className="text-sm text-slate-500">
                  {cardTypeValue === "offer" ? "Карточка будет показана как предложение исполнителя." : "Карточка будет показана как запрос заказчика."}
                </p>
                {form.formState.errors.card_type ? <p className="text-sm text-red-600">{form.formState.errors.card_type.message}</p> : null}
              </div>

              <div className="space-y-2">
                <Label htmlFor="kind">Формат</Label>
                <select id="kind" className={selectClassName} {...form.register("kind")}>
                  <option value="service">Услуга</option>
                  <option value="product">Продукт</option>
                </select>
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="title">Название</Label>
              <Input id="title" placeholder="Сведение и мастеринг трека под релиз" {...form.register("title")} className="rounded-2xl border-slate-300" />
              {form.formState.errors.title ? <p className="text-sm text-red-600">{form.formState.errors.title.message}</p> : null}
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Описание</Label>
              <Textarea
                id="description"
                placeholder="Опишите задачу, результат, сроки, стилистику и важные условия работы."
                className="rounded-2xl border-slate-300"
                {...form.register("description")}
              />
              {form.formState.errors.description ? <p className="text-sm text-red-600">{form.formState.errors.description.message}</p> : null}
            </div>

            <div className="grid gap-6 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="price">Цена</Label>
                <Input id="price" type="number" min="1" step="1" className="rounded-2xl border-slate-300" {...form.register("price", { valueAsNumber: true })} />
                <p className="text-sm text-slate-500">Укажите итоговую стоимость в рублях.</p>
                {form.formState.errors.price ? <p className="text-sm text-red-600">{form.formState.errors.price.message}</p> : null}
              </div>

              <div className="space-y-2">
                <Label htmlFor="tags">Теги</Label>
                <Input id="tags" placeholder="mixing, vocal, podcast" className="rounded-2xl border-slate-300" {...form.register("tags")} />
                <p className="text-sm text-slate-500">Указывайте теги через запятую.</p>
              </div>
            </div>

            <label className="flex items-start gap-3 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-4">
              <input type="checkbox" className="mt-1 h-4 w-4 rounded border-slate-300" {...form.register("is_published")} />
              <div className="space-y-1">
                <span className="text-sm font-medium text-slate-900">Опубликовать сразу</span>
                <p className="text-sm text-slate-500">Если снять флажок, карточка сохранится, но не будет показана в публичном каталоге.</p>
              </div>
            </label>

            {mode === "create" ? (
              <div className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-4 text-sm leading-6 text-slate-600">
                После создания карточки мы сразу переведем вас на экран редактирования, где можно загрузить обложку, preview и
                приватный полный файл.
              </div>
            ) : null}

            {mutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(mutation.error)}</p> : null}

            <div className="flex flex-wrap gap-3">
              <Button type="submit" size="lg" className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800" disabled={mutation.isPending}>
                {mutation.isPending ? (mode === "create" ? "Создаем..." : "Сохраняем...") : mode === "create" ? "Создать карточку" : "Сохранить изменения"}
              </Button>
              <Button asChild variant="outline" size="lg" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                <Link href={mode === "edit" && cardId ? `/cards/${cardId}` : "/dashboard"}>Отмена</Link>
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {mode === "edit" && card ? (
        <CardMediaManager
          cardId={card.id}
          coverUrl={card.cover_url}
          previewUrls={card.preview_urls ?? []}
          cardType={card.card_type}
          kind={card.kind}
          role={user?.role}
        />
      ) : null}
    </div>
  );
}
