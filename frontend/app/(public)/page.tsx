import Link from "next/link";
import { ArrowRight, CirclePlay, MessagesSquare, ShieldCheck, Sparkles, WalletCards } from "lucide-react";

import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";

const benefits = [
  {
    title: "Каталог музыки и услуг",
    description: "Публикуйте биты, услуги по сведению, мастерингу и саунд-дизайну или находите исполнителей под свой запрос.",
    icon: CirclePlay
  },
  {
    title: "Сделка в одном окне",
    description: "Отклики, заказ, чат, файлы, спор и отзыв собраны в одном понятном сценарии без лишней путаницы.",
    icon: MessagesSquare
  },
  {
    title: "Безопасные расчёты",
    description: "Пополнение баланса, удержание средств и завершение сделки работают как единый прозрачный поток.",
    icon: WalletCards
  },
  {
    title: "Доверие к исполнителям",
    description: "Публичные профили, проверенные отзывы и история работы помогают быстрее выбирать людей для сотрудничества.",
    icon: ShieldCheck
  }
];

const steps = [
  "Найдите готовое предложение или создайте собственный запрос на нужный результат.",
  "Обсудите детали в откликах и чате заказа, договоритесь о формате и сроках.",
  "Получите файлы, проверьте результат и завершите сделку с отзывом."
];

export default function HomePage() {
  return (
    <main className="page-frame">
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,_rgba(15,23,42,0.12),_transparent_32%),radial-gradient(circle_at_bottom_right,_rgba(71,85,105,0.14),_transparent_28%)]" />
        <div className="container relative py-16 md:py-24">
          <div className="grid gap-10 lg:grid-cols-[minmax(0,1.2fr)_minmax(320px,0.8fr)] lg:items-center">
            <div className="space-y-8">
              <Badge className="w-fit rounded-full bg-slate-950 px-4 py-1.5 text-white" variant="secondary">
                SoundMarket
              </Badge>

              <div className="space-y-5">
                <h1 className="max-w-4xl text-5xl font-semibold tracking-tight text-slate-950 md:text-6xl">
                  Музыкальный маркетплейс для сделок, файлов и рабочего общения.
                </h1>
                <p className="max-w-2xl text-lg leading-8 text-slate-600">
                  Находите готовые биты и услуги, собирайте заказы под конкретный запрос, безопасно ведите оплату и передавайте результат
                  в одном аккуратном пространстве.
                </p>
              </div>

              <div className="flex flex-wrap gap-3">
                <Button asChild size="lg" className="rounded-2xl bg-slate-950 px-6 text-white hover:bg-slate-800">
                  <Link href="/catalog">
                    Перейти в каталог
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Link>
                </Button>
                <Button asChild size="lg" variant="outline" className="rounded-2xl border-slate-300 bg-white px-6 text-slate-900 hover:bg-slate-100">
                  <Link href="/register">Создать профиль</Link>
                </Button>
              </div>

              <div className="grid max-w-3xl gap-3 sm:grid-cols-3">
                <div className="rounded-3xl border border-slate-200 bg-white/80 px-5 py-4 shadow-sm">
                  <p className="text-xs uppercase tracking-[0.18em] text-slate-500">Для заказчиков</p>
                  <p className="mt-2 text-sm leading-6 text-slate-700">Запросы, отклики, безопасная сделка и контроль результата.</p>
                </div>
                <div className="rounded-3xl border border-slate-200 bg-white/80 px-5 py-4 shadow-sm">
                  <p className="text-xs uppercase tracking-[0.18em] text-slate-500">Для исполнителей</p>
                  <p className="mt-2 text-sm leading-6 text-slate-700">Профиль, карточки, deliverables и понятный поток работы по заказу.</p>
                </div>
                <div className="rounded-3xl border border-slate-200 bg-white/80 px-5 py-4 shadow-sm">
                  <p className="text-xs uppercase tracking-[0.18em] text-slate-500">Для сделок</p>
                  <p className="mt-2 text-sm leading-6 text-slate-700">Чат, статусы, споры, отзывы и приватная выдача файлов.</p>
                </div>
              </div>
            </div>

            <div className="rounded-[2rem] border border-slate-200/80 bg-white/95 p-6 shadow-[0_30px_120px_-52px_rgba(15,23,42,0.45)]">
              <div className="space-y-5">
                <div className="flex items-center justify-between">
                  <p className="text-sm font-medium text-slate-500">Как работает SoundMarket</p>
                  <Sparkles className="h-5 w-5 text-slate-400" />
                </div>
                <div className="space-y-4">
                  {steps.map((step, index) => (
                    <div key={step} className="flex gap-4 rounded-3xl border border-slate-200 bg-slate-50/70 p-4">
                      <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl bg-slate-950 text-sm font-semibold text-white">
                        {index + 1}
                      </div>
                      <p className="text-sm leading-7 text-slate-700">{step}</p>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="container py-12 md:py-16">
        <div className="mb-8 flex items-end justify-between gap-4">
          <div className="space-y-2">
            <Badge variant="secondary" className="rounded-full bg-slate-100 text-slate-700">
              Возможности платформы
            </Badge>
            <h2 className="text-3xl font-semibold tracking-tight text-slate-950">Всё, что нужно для музыкальной сделки, уже рядом</h2>
          </div>
        </div>

        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          {benefits.map((item) => {
            const Icon = item.icon;

            return (
              <Card key={item.title} className="border-slate-200/80 bg-white/95 shadow-[0_18px_60px_-36px_rgba(15,23,42,0.25)]">
                <CardHeader className="space-y-4">
                  <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-slate-950 text-white">
                    <Icon className="h-5 w-5" />
                  </div>
                  <CardTitle className="text-xl text-slate-950">{item.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm leading-7 text-slate-600">{item.description}</p>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </section>
    </main>
  );
}
