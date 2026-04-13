"use client";

import { ChangeEvent, useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AudioLines, FileAudio, ImagePlus, LockKeyhole, UploadCloud } from "lucide-react";

import { cardsApi } from "@/entities/card/api/cards";
import { getErrorMessage } from "@/lib/api/errors";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { MediaFile } from "@/shared/types/api";

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

function createFileFormData(file: File) {
  const formData = new FormData();
  formData.append("file", file);
  return formData;
}

type CardMediaManagerProps = {
  cardId: string;
  coverUrl?: string;
  previewUrls: string[];
};

export function CardMediaManager({ cardId, coverUrl, previewUrls }: CardMediaManagerProps) {
  const queryClient = useQueryClient();
  const [coverUpload, setCoverUpload] = useState<MediaFile | null>(null);
  const [previewUpload, setPreviewUpload] = useState<MediaFile | null>(null);
  const [fullUpload, setFullUpload] = useState<MediaFile | null>(null);

  const fullDownloadQuery = useQuery({
    queryKey: ["card", cardId, "full-download"],
    queryFn: () => cardsApi.getFullDownloadUrl(cardId),
    retry: false
  });

  const invalidateCardQueries = async () => {
    await queryClient.invalidateQueries({
      predicate: (query) => Array.isArray(query.queryKey) && query.queryKey.some((part) => part === "card" || part === "cards")
    });
    await queryClient.invalidateQueries({ queryKey: ["card", cardId, "full-download"] });
  };

  const coverMutation = useMutation({
    mutationFn: (file: File) => cardsApi.uploadCover(cardId, createFileFormData(file)),
    onSuccess: async (media) => {
      setCoverUpload(media);
      await invalidateCardQueries();
    }
  });

  const previewMutation = useMutation({
    mutationFn: (file: File) => cardsApi.uploadPreview(cardId, createFileFormData(file)),
    onSuccess: async (media) => {
      setPreviewUpload(media);
      await invalidateCardQueries();
    }
  });

  const fullMutation = useMutation({
    mutationFn: (file: File) => cardsApi.uploadFull(cardId, createFileFormData(file)),
    onSuccess: async (media) => {
      setFullUpload(media);
      await invalidateCardQueries();
    }
  });

  const previewUrl = previewUrls[0];
  const hasFullDownload = fullDownloadQuery.isSuccess && Boolean(fullDownloadQuery.data?.url);

  const coverStatusText = useMemo(() => {
    if (coverMutation.isPending) return "Загружаем новую обложку...";
    if (coverMutation.isSuccess) return "Обложка успешно обновлена и сразу доступна в карточке.";
    return coverUrl ? "Обложка уже загружена и используется в каталоге и на странице карточки." : "Обложка пока не загружена.";
  }, [coverMutation.isPending, coverMutation.isSuccess, coverUrl]);

  const previewStatusText = useMemo(() => {
    if (previewMutation.isPending) return "Загружаем preview-файл...";
    if (previewMutation.isSuccess) return "Preview-файл успешно обновлен.";
    return previewUrl ? "Preview уже загружен и доступен в карточке." : "Preview пока не загружен.";
  }, [previewMutation.isPending, previewMutation.isSuccess, previewUrl]);

  const fullStatusText = useMemo(() => {
    if (fullMutation.isPending) return "Загружаем полный файл...";
    if (fullMutation.isSuccess) return "Полный файл успешно загружен.";
    if (hasFullDownload) return "Полный файл загружен и доступен по приватной ссылке.";
    if (fullDownloadQuery.isFetching) return "Проверяем наличие полного файла...";
    return "Полный файл пока не загружен.";
  }, [fullDownloadQuery.isFetching, fullMutation.isPending, fullMutation.isSuccess, hasFullDownload]);

  const handleFile =
    (callback: (file: File) => void) =>
    (event: ChangeEvent<HTMLInputElement>) => {
      const file = event.target.files?.[0];
      if (!file) return;
      callback(file);
      event.target.value = "";
    };

  return (
    <Card className="overflow-hidden border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.32)]">
      <CardHeader className="space-y-4 border-b border-slate-200/80 bg-[linear-gradient(135deg,rgba(15,23,42,0.04),rgba(71,85,105,0.02))]">
        <div className="flex flex-wrap items-center gap-3">
          <Badge className="bg-slate-900/90 text-white hover:bg-slate-900" variant="secondary">
            Медиа карточки
          </Badge>
          <Badge variant="outline">SoundMarket</Badge>
        </div>
        <div className="space-y-2">
          <CardTitle className="text-2xl text-slate-950">Обложка, preview и приватный исходник</CardTitle>
          <CardDescription className="max-w-3xl text-base leading-7 text-slate-600">
            Обложка и preview работают как публичная витрина карточки. Полный файл остается приватным и открывается только по
            защищенной ссылке.
          </CardDescription>
        </div>
      </CardHeader>

      <CardContent className="grid gap-6 p-6 xl:grid-cols-3">
        <section className="rounded-2xl border border-slate-200 bg-slate-50/80 p-5">
          <div className="mb-4 flex items-center gap-3">
            <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-900 text-white">
              <ImagePlus className="h-5 w-5" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-slate-950">Обложка карточки</h3>
              <p className="text-sm leading-6 text-slate-600">Публичное изображение для каталога, detail page и блока preview.</p>
            </div>
          </div>

          <div className="space-y-4">
            <div className="overflow-hidden rounded-[1.5rem] border border-slate-200 bg-white">
              {coverUrl ? (
                <img src={coverUrl} alt="Обложка карточки" className="aspect-[16/10] w-full object-cover" />
              ) : (
                <div className="flex aspect-[16/10] items-center justify-center bg-[linear-gradient(145deg,rgba(15,23,42,0.92),rgba(51,65,85,0.88))] px-6 text-center text-sm text-white/80">
                  Пока без обложки. Добавьте изображение, чтобы карточка выглядела сильнее в каталоге и на detail page.
                </div>
              )}
            </div>

            <div className="rounded-2xl border border-slate-200 bg-white p-4 text-sm text-slate-700">
              <p>{coverStatusText}</p>
              {coverUpload ? (
                <p className="mt-2 text-slate-500">
                  Последний файл: {coverUpload.original_filename} • {formatBytes(coverUpload.size_bytes)}
                </p>
              ) : null}
            </div>

            {coverMutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(coverMutation.error)}</p> : null}

            <label className="flex cursor-pointer items-center justify-center gap-3 rounded-2xl bg-slate-900 px-4 py-3 text-sm font-medium text-white transition hover:bg-slate-800">
              <UploadCloud className="h-4 w-4" />
              {coverMutation.isPending ? "Загрузка..." : coverUrl ? "Обновить обложку" : "Загрузить обложку"}
              <input type="file" accept="image/*" className="hidden" onChange={handleFile((file) => coverMutation.mutate(file))} disabled={coverMutation.isPending} />
            </label>
          </div>
        </section>

        <section className="rounded-2xl border border-slate-200 bg-slate-50/80 p-5">
          <div className="mb-4 flex items-center gap-3">
            <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-900 text-white">
              <AudioLines className="h-5 w-5" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-slate-950">Preview-файл</h3>
              <p className="text-sm leading-6 text-slate-600">Публичный аудио-фрагмент, который увидят посетители карточки.</p>
            </div>
          </div>

          <div className="space-y-4">
            {previewUrl ? (
              <div className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
                <audio controls className="w-full">
                  <source src={previewUrl} />
                </audio>
              </div>
            ) : (
              <div className="rounded-2xl border border-dashed border-slate-300 bg-white/80 p-5 text-sm leading-6 text-slate-600">
                Preview еще не добавлен. После загрузки он сразу появится здесь и на публичной странице карточки.
              </div>
            )}

            <div className="rounded-2xl border border-slate-200 bg-white p-4 text-sm text-slate-700">
              <p>{previewStatusText}</p>
              {previewUpload ? (
                <p className="mt-2 text-slate-500">
                  Последний файл: {previewUpload.original_filename} • {formatBytes(previewUpload.size_bytes)}
                </p>
              ) : null}
            </div>

            {previewMutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(previewMutation.error)}</p> : null}

            <label className="flex cursor-pointer items-center justify-center gap-3 rounded-2xl bg-slate-900 px-4 py-3 text-sm font-medium text-white transition hover:bg-slate-800">
              <UploadCloud className="h-4 w-4" />
              {previewMutation.isPending ? "Загрузка..." : "Загрузить preview"}
              <input type="file" accept="audio/*" className="hidden" onChange={handleFile((file) => previewMutation.mutate(file))} disabled={previewMutation.isPending} />
            </label>
          </div>
        </section>

        <section className="rounded-2xl border border-slate-200 bg-slate-50/80 p-5">
          <div className="mb-4 flex items-center gap-3">
            <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-200 text-slate-900">
              <LockKeyhole className="h-5 w-5" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-slate-950">Полный файл</h3>
              <p className="text-sm leading-6 text-slate-600">Приватный исходник. Он не отображается публично и доступен только по защищенной ссылке.</p>
            </div>
          </div>

          <div className="space-y-4">
            <div className="rounded-2xl border border-slate-200 bg-white p-4">
              <div className="flex items-start gap-3">
                <div className="mt-1 flex h-9 w-9 items-center justify-center rounded-xl bg-slate-100 text-slate-700">
                  <FileAudio className="h-4 w-4" />
                </div>
                <div className="min-w-0 space-y-2">
                  <p className="text-sm leading-6 text-slate-700">{fullStatusText}</p>
                  {fullUpload ? (
                    <p className="text-sm text-slate-500">
                      Последний файл: {fullUpload.original_filename} • {formatBytes(fullUpload.size_bytes)}
                    </p>
                  ) : hasFullDownload ? (
                    <p className="text-sm text-slate-500">Приватная ссылка уже доступна для владельца карточки.</p>
                  ) : null}
                </div>
              </div>
            </div>

            {fullMutation.isError ? <p className="text-sm text-red-600">{getErrorMessage(fullMutation.error)}</p> : null}
            {fullDownloadQuery.isError ? <p className="text-sm text-slate-500">Пока приватного файла нет, либо он еще не загружен.</p> : null}

            <div className="flex flex-wrap gap-3">
              <label className="flex cursor-pointer items-center justify-center gap-3 rounded-2xl bg-slate-900 px-4 py-3 text-sm font-medium text-white transition hover:bg-slate-800">
                <UploadCloud className="h-4 w-4" />
                {fullMutation.isPending ? "Загрузка..." : "Загрузить полный файл"}
                <input type="file" accept="audio/*" className="hidden" onChange={handleFile((file) => fullMutation.mutate(file))} disabled={fullMutation.isPending} />
              </label>

              {hasFullDownload ? (
                <Button asChild variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
                  <a href={fullDownloadQuery.data.url} target="_blank" rel="noreferrer">
                    Открыть приватную ссылку
                  </a>
                </Button>
              ) : null}
            </div>
          </div>
        </section>
      </CardContent>
    </Card>
  );
}
