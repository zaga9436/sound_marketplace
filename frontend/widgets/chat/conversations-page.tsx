"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";

import { chatApi } from "@/entities/chat/api/chat";
import { getErrorMessage } from "@/lib/api/errors";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";

export function ConversationsPage() {
  const query = useQuery({
    queryKey: ["chats"],
    queryFn: () => chatApi.listConversations(50),
    refetchInterval: 5000
  });

  if (query.isLoading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 4 }).map((_, index) => (
          <div key={index} className="surface h-28 animate-pulse bg-white/70" />
        ))}
      </div>
    );
  }

  if (query.isError) {
    return (
      <Card className="border-destructive/20 bg-white/95">
        <CardContent className="pt-6">
          <p className="text-destructive">{getErrorMessage(query.error)}</p>
        </CardContent>
      </Card>
    );
  }

  const conversations = query.data ?? [];

  return (
    <div className="space-y-8">
      <section className="space-y-3">
        <Badge className="bg-slate-900/90 text-white" variant="secondary">
          Чаты
        </Badge>
        <div className="space-y-2">
          <h1>Рабочие диалоги по сделкам</h1>
          <p>Каждый чат привязан к конкретному заказу. Отсюда удобно перейти в нужную сделку и продолжить общение.</p>
        </div>
      </section>

      {conversations.length === 0 ? (
        <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
          <CardContent className="pt-6">
            <p>Пока активных чатов нет. Они появятся, когда по заказам начнут приходить сообщения.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {conversations.map((conversation) => (
            <Card key={conversation.order_id} className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.16)]">
              <CardHeader className="flex flex-row items-start justify-between gap-4">
                <div className="space-y-2">
                  <CardTitle className="text-lg">Заказ {conversation.order_id}</CardTitle>
                  <p className="text-sm text-slate-500">{conversation.last_message || "Сообщений пока нет"}</p>
                  {conversation.last_message_at ? <p className="text-xs text-slate-500">{new Date(conversation.last_message_at).toLocaleString("ru-RU")}</p> : null}
                </div>
                <div className="flex flex-col items-end gap-2">
                  {conversation.unread_count > 0 ? <Badge variant="secondary">{conversation.unread_count} новых</Badge> : <Badge variant="outline">Все прочитано</Badge>}
                  <Button asChild className="rounded-2xl bg-slate-900 text-white hover:bg-slate-800">
                    <Link href={`/orders/${conversation.order_id}`}>Открыть заказ</Link>
                  </Button>
                </div>
              </CardHeader>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
