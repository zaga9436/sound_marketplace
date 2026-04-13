"use client";

import { Badge } from "@/shared/ui/badge";
import { OrderStatus } from "@/shared/types/api";

const labels: Record<OrderStatus, string> = {
  created: "Создан",
  on_hold: "На холде",
  in_progress: "В работе",
  review: "На проверке",
  completed: "Завершен",
  dispute: "Спор",
  cancelled: "Отменен"
};

const classes: Record<OrderStatus, string> = {
  created: "bg-slate-100 text-slate-700",
  on_hold: "bg-amber-100 text-amber-800",
  in_progress: "bg-sky-100 text-sky-800",
  review: "bg-violet-100 text-violet-800",
  completed: "bg-emerald-100 text-emerald-800",
  dispute: "bg-rose-100 text-rose-800",
  cancelled: "bg-slate-200 text-slate-700"
};

export function OrderStatusBadge({ status }: { status: OrderStatus }) {
  return (
    <Badge variant="secondary" className={classes[status]}>
      {labels[status]}
    </Badge>
  );
}
