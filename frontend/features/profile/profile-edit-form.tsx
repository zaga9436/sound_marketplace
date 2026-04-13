"use client";

import Link from "next/link";
import { useEffect } from "react";
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

const profileSchema = z.object({
  display_name: z.string().min(2, "Имя должно быть не короче 2 символов."),
  bio: z.string().max(1000, "Описание должно быть короче 1000 символов.").default("")
});

type ProfileFormValues = z.infer<typeof profileSchema>;

export function ProfileEditForm() {
  const router = useRouter();
  const queryClient = useQueryClient();
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

  const mutation = useMutation({
    mutationFn: profilesApi.updateMe,
    onSuccess: async (profile) => {
      updateProfile(profile);
      await queryClient.invalidateQueries({
        predicate: (query) => Array.isArray(query.queryKey) && query.queryKey.some((part) => part === "profile")
      });
      router.push("/profile");
    }
  });

  const onSubmit = form.handleSubmit((values) => mutation.mutate(values));

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

  return (
    <Card className="border-slate-200/80 bg-white/95 shadow-[0_24px_80px_-40px_rgba(15,23,42,0.32)]">
      <CardHeader className="space-y-4">
        <Badge className="bg-slate-900/90 text-white" variant="secondary">
          Редактирование профиля
        </Badge>
        <div className="space-y-2">
          <CardTitle className="text-3xl text-slate-950">Обновите публичный образ</CardTitle>
          <CardDescription className="text-base leading-7 text-slate-600">
            Профиль — это первая точка доверия после карточки. Сделайте его коротким, понятным и убедительным.
          </CardDescription>
        </div>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="space-y-6">
          <div className="space-y-2">
            <Label htmlFor="display_name">Отображаемое имя</Label>
            <Input id="display_name" className="rounded-2xl border-slate-300" placeholder="Алексей, mix engineer" {...form.register("display_name")} />
            {form.formState.errors.display_name ? <p className="text-sm text-destructive">{form.formState.errors.display_name.message}</p> : null}
          </div>

          <div className="space-y-2">
            <Label htmlFor="bio">Описание</Label>
            <Textarea id="bio" className="rounded-2xl border-slate-300" placeholder="Расскажите о своем опыте, жанрах, специализации и том, какие проекты вам особенно интересны." {...form.register("bio")} />
            {form.formState.errors.bio ? <p className="text-sm text-destructive">{form.formState.errors.bio.message}</p> : null}
          </div>

          {mutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(mutation.error)}</p> : null}

          <div className="flex flex-wrap gap-3">
            <Button type="submit" size="lg" className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800" disabled={mutation.isPending}>
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
