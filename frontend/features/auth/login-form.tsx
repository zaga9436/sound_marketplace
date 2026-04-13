"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { authApi } from "@/entities/auth/api/auth";
import { getErrorMessage } from "@/lib/api/errors";
import { useAuthStore } from "@/lib/auth/session-store";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";

const loginSchema = z.object({
  email: z.string().email("Введите корректный email."),
  password: z.string().min(6, "Минимум 6 символов.")
});

type LoginValues = z.infer<typeof loginSchema>;

export function LoginForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const setSession = useAuthStore((state) => state.setSession);

  const form = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: ""
    }
  });

  const mutation = useMutation({
    mutationFn: authApi.login,
    onSuccess: (payload) => {
      setSession(payload);
      const next = searchParams.get("next");
      router.push(next || (payload.user.role === "admin" ? "/admin" : "/dashboard"));
    }
  });

  const onSubmit = form.handleSubmit((values) => mutation.mutate(values));

  return (
    <Card className="w-full max-w-md border-slate-200/80 bg-white/95 shadow-[0_24px_80px_-40px_rgba(15,23,42,0.32)]">
      <CardHeader className="space-y-2">
        <CardTitle>Вход в SoundMarket</CardTitle>
        <CardDescription>Используйте email и пароль от аккаунта заказчика, инженера или администратора.</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="space-y-5">
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" placeholder="you@example.com" className="rounded-2xl border-slate-300" {...form.register("email")} />
            {form.formState.errors.email ? <p className="text-sm text-destructive">{form.formState.errors.email.message}</p> : null}
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">Пароль</Label>
            <Input id="password" type="password" placeholder="••••••••" className="rounded-2xl border-slate-300" {...form.register("password")} />
            {form.formState.errors.password ? <p className="text-sm text-destructive">{form.formState.errors.password.message}</p> : null}
          </div>

          {mutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(mutation.error)}</p> : null}

          <Button className="w-full rounded-2xl bg-slate-900 text-white hover:bg-slate-800" size="lg" type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? "Входим..." : "Войти"}
          </Button>

          <p className="text-center text-sm text-muted-foreground">
            Еще нет аккаунта?{" "}
            <Link href="/register" className="font-medium text-primary">
              Зарегистрироваться
            </Link>
          </p>
        </form>
      </CardContent>
    </Card>
  );
}
