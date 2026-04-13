import Link from "next/link";
import { ArrowRight, ShieldCheck, Sparkles, Waves } from "lucide-react";

import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";

const highlights = [
  {
    title: "Каталог и сделки",
    description: "Предложения, запросы, отклики, заказы и deliverables собраны в одном понятном продуктовым потоке.",
    icon: Waves
  },
  {
    title: "Прозрачный workflow",
    description: "Отзывы, verified сделки, чат по заказу и понятные статусы помогают быстро доверять платформе.",
    icon: ShieldCheck
  },
  {
    title: "Дипломный продукт",
    description: "Интерфейс уже собран вокруг реального backend flow, поэтому можно сразу показывать полноценный сценарий SoundMarket.",
    icon: Sparkles
  }
];

export default function HomePage() {
  return (
    <main className="page-frame">
      <div className="container flex min-h-[calc(100vh-72px)] flex-col justify-center py-16">
        <div className="mx-auto flex max-w-5xl flex-col gap-10">
          <div className="flex flex-col gap-5">
            <Badge className="w-fit bg-slate-900/90 text-white" variant="secondary">
              SoundMarket
            </Badge>
            <div className="max-w-3xl space-y-4">
              <h1>Современный маркетплейс для музыкальных услуг, запросов и сделок.</h1>
              <p className="text-lg text-slate-600">
                Каталог, карточки, отклики, заказы, чат, споры и отзывы уже собраны в единый интерфейс. Можно сразу переходить к реальным пользовательским сценариям.
              </p>
            </div>
            <div className="flex flex-wrap gap-3">
              <Button asChild size="lg" className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800">
                <Link href="/catalog">
                  Открыть каталог
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
              <Button asChild variant="secondary" size="lg" className="rounded-2xl bg-slate-100 text-slate-900 hover:bg-slate-200">
                <Link href="/register">Создать аккаунт</Link>
              </Button>
            </div>
          </div>

          <div className="grid gap-4 md:grid-cols-3">
            {highlights.map((item) => {
              const Icon = item.icon;

              return (
                <Card key={item.title} className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
                  <CardHeader className="space-y-4">
                    <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-900/10 text-slate-900">
                      <Icon className="h-5 w-5" />
                    </div>
                    <CardTitle className="text-slate-950">{item.title}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <p className="text-slate-600">{item.description}</p>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        </div>
      </div>
    </main>
  );
}
