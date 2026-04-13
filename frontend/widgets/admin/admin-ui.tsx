"use client";

import { PropsWithChildren, ReactNode } from "react";

import { Badge } from "@/shared/ui/badge";
import { cn } from "@/shared/utils/cn";

export function AdminPageHeader({
  title,
  description,
  actions
}: {
  title: string;
  description: string;
  actions?: ReactNode;
}) {
  return (
    <div className="flex flex-wrap items-end justify-between gap-4">
      <div className="space-y-2">
        <h1 className="text-3xl font-semibold tracking-tight text-slate-950">{title}</h1>
        <p className="max-w-3xl text-sm leading-7 text-slate-500">{description}</p>
      </div>
      {actions}
    </div>
  );
}

export function AdminSection({ children, className }: PropsWithChildren<{ className?: string }>) {
  return <section className={cn("rounded-[1.75rem] border border-slate-200 bg-white/95 p-5 shadow-[0_20px_60px_-36px_rgba(15,23,42,0.22)]", className)}>{children}</section>;
}

export function StatusBadge({ tone, children }: PropsWithChildren<{ tone?: "green" | "yellow" | "red" | "slate" | "blue" }>) {
  const toneClass =
    tone === "green"
      ? "bg-emerald-100 text-emerald-800"
      : tone === "yellow"
        ? "bg-amber-100 text-amber-800"
        : tone === "red"
          ? "bg-rose-100 text-rose-800"
          : tone === "blue"
            ? "bg-sky-100 text-sky-800"
            : "bg-slate-100 text-slate-700";

  return (
    <Badge variant="secondary" className={cn("rounded-full px-3 py-1", toneClass)}>
      {children}
    </Badge>
  );
}

export function formatRole(role?: string) {
  if (role === "customer") return "Заказчик";
  if (role === "engineer") return "Исполнитель";
  if (role === "admin") return "Администратор";
  return "Неизвестно";
}
