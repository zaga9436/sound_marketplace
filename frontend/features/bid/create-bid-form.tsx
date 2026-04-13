"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { bidsApi } from "@/entities/bid/api/bids";
import { getErrorMessage } from "@/lib/api/errors";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Textarea } from "@/shared/ui/textarea";

const bidSchema = z.object({
  price: z.coerce.number().int().min(1, "Цена должна быть больше нуля."),
  message: z.string().min(5, "Сообщение должно быть не короче 5 символов.")
});

type BidFormValues = z.infer<typeof bidSchema>;

export function CreateBidForm({ requestId }: { requestId: string }) {
  const queryClient = useQueryClient();
  const form = useForm<BidFormValues>({
    resolver: zodResolver(bidSchema),
    defaultValues: {
      price: 1000,
      message: ""
    }
  });

  const mutation = useMutation({
    mutationFn: (values: BidFormValues) => bidsApi.create(requestId, values),
    onSuccess: async () => {
      form.reset({ price: form.getValues("price"), message: "" });
      await queryClient.invalidateQueries({ queryKey: ["bids", requestId] });
    }
  });

  const onSubmit = form.handleSubmit((values) => mutation.mutate(values));

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="bid-price">Ваша цена</Label>
        <Input id="bid-price" type="number" min="1" step="1" className="rounded-2xl border-slate-300" {...form.register("price", { valueAsNumber: true })} />
        {form.formState.errors.price ? <p className="text-sm text-red-600">{form.formState.errors.price.message}</p> : null}
      </div>

      <div className="space-y-2">
        <Label htmlFor="bid-message">Комментарий</Label>
        <Textarea
          id="bid-message"
          className="rounded-2xl border-slate-300"
          placeholder="Кратко опишите, как вы планируете выполнить задачу и за какой срок."
          {...form.register("message")}
        />
        {form.formState.errors.message ? <p className="text-sm text-red-600">{form.formState.errors.message.message}</p> : null}
      </div>

      {mutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(mutation.error)}</p> : null}
      {mutation.isSuccess ? <p className="text-sm text-emerald-700">Отклик отправлен. Заказчик увидит его в списке заявок.</p> : null}

      <Button type="submit" className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800" disabled={mutation.isPending}>
        {mutation.isPending ? "Отправляем..." : "Отправить отклик"}
      </Button>
    </form>
  );
}
