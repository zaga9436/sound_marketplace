"use client";

import { ChangeEvent } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { UploadCloud } from "lucide-react";

import { deliverablesApi } from "@/entities/deliverable/api/deliverables";
import { getErrorMessage } from "@/lib/api/errors";
import { Button } from "@/shared/ui/button";

export function UploadDeliverableForm({ orderId }: { orderId: string }) {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (file: File) => deliverablesApi.upload(orderId, file),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["deliverables", orderId] });
      await queryClient.invalidateQueries({ queryKey: ["notifications"] });
    }
  });

  const handleFile = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;
    mutation.mutate(file);
    event.target.value = "";
  };

  return (
    <div className="space-y-3">
      {mutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(mutation.error)}</p> : null}
      {mutation.isSuccess ? <p className="text-sm text-emerald-700">Новая версия результата успешно загружена.</p> : null}

      <label className="inline-flex cursor-pointer items-center gap-3 rounded-2xl bg-slate-900 px-4 py-3 text-sm font-medium text-white transition hover:bg-slate-800">
        <UploadCloud className="h-4 w-4" />
        {mutation.isPending ? "Загрузка..." : "Загрузить результат"}
        <input type="file" className="hidden" onChange={handleFile} disabled={mutation.isPending} />
      </label>

      <p className="text-sm leading-6 text-slate-500">Файл будет загружен как новая версия deliverable и останется приватным.</p>
    </div>
  );
}
