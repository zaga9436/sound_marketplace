"use client";

import Link from "next/link";
import { ChangeEvent, useEffect, useRef } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { profilesApi } from "@/entities/profile/api/profiles";
import { getErrorMessage } from "@/lib/api/errors";
import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Textarea } from "@/shared/ui/textarea";
import { UserAvatar } from "@/shared/ui/user-avatar";

const profileSchema = z.object({
  display_name: z.string().min(2, "Имя должно быть не короче 2 символов."),
  bio: z.string().max(1000, "Описание должно быть короче 1000 символов.").default("")
});

type ProfileFormValues = z.infer<typeof profileSchema>;

export function ProfileEditForm() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const user = useAuthStore((state) => state.user);
  const updateProfile = useAuthStore((state) => state.updateProfile);

  const profileQuery = useQuery({
    queryKey: ["profile", "me"],
    queryFn: () => profilesApi.me()
  });

  const form = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      display_name: "",
      bio: ""
    }
  });

  useEffect(() => {
    if (profileQuery.data) {
      form.reset({
        display_name: profileQuery.data.display_name ?? "",
        bio: profileQuery.data.bio ?? ""
      });
    }
  }, [form, profileQuery.data]);

  const invalidateProfileQueries = async () => {
    await queryClient.invalidateQueries({
      predicate: (query) => Array.isArray(query.queryKey) && query.queryKey.some((part) => part === "profile" || part === "auth")
    });
  };

  const mutation = useMutation({
    mutationFn: profilesApi.updateMe,
    onSuccess: async (profile) => {
      updateProfile(profile);
      await invalidateProfileQueries();
      router.push("/profile");
    }
  });

  const avatarMutation = useMutation({
    mutationFn: profilesApi.uploadAvatar,
    onSuccess: async () => {
      const profile = await queryClient.fetchQuery({
        queryKey: ["profile", "me"],
        queryFn: () => profilesApi.me()
      });
      updateProfile(profile);
      await invalidateProfileQueries();
    }
  });

  const onSubmit = form.handleSubmit((values) => mutation.mutate(values));

  const handleAvatarChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;
    const formData = new FormData();
    formData.append("file", file);
    avatarMutation.mutate(formData);
    event.target.value = "";
  };

  if (profileQuery.isLoading) {
    return <div className="surface h-[420px] animate-pulse bg-white/70" />;
  }

  if (profileQuery.isError) {
    return (
      <Card className="border-destructive/20 bg-white/95">
        <CardContent className="pt-6">
          <p className="text-destructive">{getErrorMessage(profileQuery.error)}</p>
        </CardContent>
      </Card>
    );
  }

  const profile = profileQuery.data;

  return (
    <Card className="border-slate-200/80 bg-white/95 shadow-[0_24px_80px_-40px_rgba(15,23,42,0.32)]">
      <CardHeader className="space-y-4">
        <Badge className="bg-slate-950 text-white" variant="secondary">
          Редактирование профиля
        </Badge>
        <div className="space-y-2">
          <CardTitle className="text-3xl text-slate-950">Обновите публичный образ</CardTitle>
          <CardDescription className="text-base leading-7 text-slate-600">
            Аккуратная аватарка, понятное имя и короткое описание помогают быстрее вызывать доверие и выглядеть как реальный участник музыкального сервиса.
          </CardDescription>
        </div>
      </CardHeader>
      <CardContent>
        <div className="mb-8 rounded-[1.75rem] border border-slate-200 bg-slate-50/80 p-5">
          <div className="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
            <div className="flex items-center gap-4">
              <UserAvatar
                avatarUrl={profile?.avatar_url}
                name={profile?.display_name}
                email={user?.email}
                className="h-20 w-20 rounded-[1.75rem] text-xl"
              />
              <div className="space-y-1">
                <p className="text-sm font-medium text-slate-950">{profile?.display_name || user?.email || "Профиль"}</p>
                <p className="text-sm text-slate-500">Аватар используется в шапке, кабинете и публичном профиле.</p>
              </div>
            </div>
            <div className="flex flex-wrap gap-3">
              <input ref={fileInputRef} type="file" accept="image/*" className="hidden" onChange={handleAvatarChange} />
              <Button
                type="button"
                variant="outline"
                className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100"
                disabled={avatarMutation.isPending}
                onClick={() => fileInputRef.current?.click()}
              >
                {avatarMutation.isPending ? "Загружаем аватарку..." : "Загрузить аватарку"}
              </Button>
            </div>
          </div>
          {avatarMutation.isError ? <p className="mt-3 text-sm text-destructive">{getErrorMessage(avatarMutation.error)}</p> : null}
          {avatarMutation.isSuccess ? <p className="mt-3 text-sm text-emerald-700">Аватарка обновлена и уже используется в интерфейсе.</p> : null}
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="space-y-2">
            <Label htmlFor="display_name">Отображаемое имя</Label>
            <Input
              id="display_name"
              className="rounded-2xl border-slate-300"
              placeholder="Например, Алексей — beatmaker и sound designer"
              {...form.register("display_name")}
            />
            {form.formState.errors.display_name ? <p className="text-sm text-destructive">{form.formState.errors.display_name.message}</p> : null}
          </div>

          <div className="space-y-2">
            <Label htmlFor="bio">Описание</Label>
            <Textarea
              id="bio"
              className="rounded-2xl border-slate-300"
              placeholder="Расскажите о жанрах, опыте, любимых форматах работы и о том, чем вы особенно полезны клиентам или артистам."
              {...form.register("bio")}
            />
            {form.formState.errors.bio ? <p className="text-sm text-destructive">{form.formState.errors.bio.message}</p> : null}
          </div>

          {mutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(mutation.error)}</p> : null}

          <div className="flex flex-wrap gap-3">
            <Button type="submit" size="lg" className="rounded-2xl bg-slate-950 text-white hover:bg-slate-800" disabled={mutation.isPending}>
              {mutation.isPending ? "Сохраняем..." : "Сохранить профиль"}
            </Button>
            <Button asChild variant="outline" size="lg" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100">
              <Link href="/profile">Отмена</Link>
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
