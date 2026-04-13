"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { chatApi } from "@/entities/chat/api/chat";
import { getErrorMessage } from "@/lib/api/errors";
import { Button } from "@/shared/ui/button";
import { Textarea } from "@/shared/ui/textarea";

const schema = z.object({
  body: z.string().min(1, "Введите сообщение.")
});

type Values = z.infer<typeof schema>;

export function MessageInput({ orderId }: { orderId: string }) {
  const queryClient = useQueryClient();
  const form = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: { body: "" }
  });

  const mutation = useMutation({
    mutationFn: (values: Values) => chatApi.sendMessage(orderId, values.body),
    onSuccess: async () => {
      form.reset();
      await queryClient.invalidateQueries({ queryKey: ["chat", orderId, "messages"] });
      await queryClient.invalidateQueries({ queryKey: ["chats"] });
      await queryClient.invalidateQueries({ queryKey: ["notifications"] });
    }
  });

  const onSubmit = form.handleSubmit((values) => mutation.mutate(values));

  return (
    <form onSubmit={onSubmit} className="space-y-3">
      <Textarea
        className="min-h-[110px] rounded-2xl border-slate-300"
        placeholder="Напишите сообщение по заказу: детали, правки, сроки, уточнения."
        {...form.register("body")}
      />
      {form.formState.errors.body ? <p className="text-sm text-red-600">{form.formState.errors.body.message}</p> : null}
      {mutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(mutation.error)}</p> : null}
      <Button type="submit" className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800" disabled={mutation.isPending}>
        {mutation.isPending ? "Отправляем..." : "Отправить сообщение"}
      </Button>
    </form>
  );
}
