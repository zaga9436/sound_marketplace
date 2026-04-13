"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
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

const registerSchema = z.object({
  email: z.string().email("Введите корректный email."),
  password: z.string().min(6, "Минимум 6 символов."),
  role: z.enum(["customer", "engineer"])
});

type RegisterValues = z.infer<typeof registerSchema>;

export function RegisterForm() {
  const router = useRouter();
  const setSession = useAuthStore((state) => state.setSession);

  const form = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: "",
      password: "",
      role: "customer"
    }
  });

  const mutation = useMutation({
    mutationFn: authApi.register,
    onSuccess: (payload) => {
      setSession(payload);
      router.push("/dashboard");
    }
  });

  const onSubmit = form.handleSubmit((values) => mutation.mutate(values));

  return (
    <Card className="w-full max-w-md border-slate-200/80 bg-white/95 shadow-[0_24px_80px_-40px_rgba(15,23,42,0.32)]">
      <CardHeader className="space-y-2">
        <CardTitle>Создать аккаунт</CardTitle>
        <CardDescription>Выберите роль и начните работать в SoundMarket как заказчик или инженер.</CardDescription>
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

          <div className="space-y-2">
            <Label htmlFor="role">Роль</Label>
            <select
              id="role"
              className="flex h-11 w-full rounded-2xl border border-slate-300 bg-white px-4 py-2 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-500"
              {...form.register("role")}
            >
              <option value="customer">Заказчик</option>
              <option value="engineer">Инженер</option>
            </select>
          </div>

          {mutation.isError ? <p className="text-sm text-destructive">{getErrorMessage(mutation.error)}</p> : null}

          <Button className="w-full rounded-2xl bg-slate-900 text-white hover:bg-slate-800" size="lg" type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? "Создаем аккаунт..." : "Зарегистрироваться"}
          </Button>

          <p className="text-center text-sm text-muted-foreground">
            Уже есть аккаунт?{" "}
            <Link href="/login" className="font-medium text-primary">
              Войти
            </Link>
          </p>
        </form>
      </CardContent>
    </Card>
  );
}
