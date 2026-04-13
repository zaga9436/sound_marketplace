"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { reviewsApi } from "@/entities/review/api/reviews";
import { getErrorMessage } from "@/lib/api/errors";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Textarea } from "@/shared/ui/textarea";

const schema = z.object({
  rating: z.coerce.number().int().min(1, "Оценка от 1 до 5.").max(5, "Оценка от 1 до 5."),
  text: z.string().min(3, "Добавьте короткий отзыв.")
});

type Values = z.infer<typeof schema>;

export function ReviewForm({ orderId, targetUserId }: { orderId: string; targetUserId: string }) {
  const queryClient = useQueryClient();
  const form = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: {
      rating: 5,
      text: ""
    }
  });

  const mutation = useMutation({
    mutationFn: (values: Values) => reviewsApi.create(orderId, values.rating, values.text),
    onSuccess: async () => {
      form.reset({ rating: 5, text: "" });
      await queryClient.invalidateQueries({ queryKey: ["profile", targetUserId] });
      await queryClient.invalidateQueries({ queryKey: ["profile", targetUserId, "reviews"] });
      await queryClient.invalidateQueries({ queryKey: ["order", orderId] });
    }
  });

  const onSubmit = form.handleSubmit((values) => mutation.mutate(values));

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="review-rating">Оценка</Label>
        <Input id="review-rating" type="number" min="1" max="5" step="1" className="rounded-2xl border-slate-300" {...form.register("rating", { valueAsNumber: true })} />
        {form.formState.errors.rating ? <p className="text-sm text-red-600">{form.formState.errors.rating.message}</p> : null}
      </div>
      <div className="space-y-2">
        <Label htmlFor="review-text">Отзыв</Label>
        <Textarea id="review-text" className="rounded-2xl border-slate-300" placeholder="Опишите, как прошла работа и что понравилось." {...form.register("text")} />
        {form.formState.errors.text ? <p className="text-sm text-red-600">{form.formState.errors.text.message}</p> : null}
      </div>
      {mutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(mutation.error)}</p> : null}
      {mutation.isSuccess ? <p className="text-sm text-emerald-700">Отзыв сохранен и уже участвует в рейтинге исполнителя.</p> : null}
      <Button type="submit" className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800" disabled={mutation.isPending}>
        {mutation.isPending ? "Публикуем..." : "Оставить отзыв"}
      </Button>
    </form>
  );
}
