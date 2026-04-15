"use client";

import { useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { chatApi } from "@/entities/chat/api/chat";
import { MessageInput } from "@/features/chat/message-input";
import { getErrorMessage } from "@/lib/api/errors";
import { useAuthStore } from "@/lib/auth/session-store";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";

export function OrderChatPanel({ orderId }: { orderId: string }) {
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);

  const messagesQuery = useQuery({
    queryKey: ["chat", orderId, "messages"],
    queryFn: () => chatApi.listMessages(orderId, 50),
    refetchInterval: 5000
  });

  const markReadMutation = useMutation({
    mutationFn: () => chatApi.markRead(orderId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["chats"] });
      await queryClient.invalidateQueries({ queryKey: ["notifications"] });
    }
  });

  useEffect(() => {
    if ((messagesQuery.data?.length ?? 0) > 0 && !markReadMutation.isPending) {
      markReadMutation.mutate();
    }
  }, [messagesQuery.data?.length]);

  return (
    <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
      <CardHeader>
        <CardTitle>Чат по заказу</CardTitle>
      </CardHeader>
      <CardContent className="space-y-5">
        {messagesQuery.isLoading ? (
          <div className="space-y-3">
            {Array.from({ length: 3 }).map((_, index) => (
              <div key={index} className="h-20 animate-pulse rounded-2xl border border-slate-200 bg-slate-50" />
            ))}
          </div>
        ) : messagesQuery.isError ? (
          <p className="text-sm text-red-600">{getErrorMessage(messagesQuery.error)}</p>
        ) : messagesQuery.data && messagesQuery.data.length > 0 ? (
          <div className="max-h-[420px] space-y-3 overflow-y-auto pr-1">
            {messagesQuery.data.map((message) => {
              const isOwn = message.sender_id === user?.id;
              return (
                <div key={message.id} className={`flex ${isOwn ? "justify-end" : "justify-start"}`}>
                  <div
                    className={`max-w-[82%] rounded-[1.25rem] border px-4 py-3 shadow-sm ${
                      isOwn ? "border-slate-300 bg-slate-100 text-slate-950" : "border-slate-200 bg-white text-slate-900"
                    }`}
                  >
                    <p className="whitespace-pre-wrap break-words text-sm leading-6">{message.body}</p>
                    <p className="mt-2 text-xs text-slate-500">
                      {isOwn ? "Вы • " : ""}
                      {new Date(message.created_at).toLocaleString("ru-RU")}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        ) : (
          <div className="rounded-2xl border border-dashed border-slate-300 bg-slate-50 px-4 py-5 text-sm text-slate-500">
            Пока сообщений нет. Здесь можно обсуждать детали работы, правки, сроки и ход выполнения заказа.
          </div>
        )}

        <MessageInput orderId={orderId} />
      </CardContent>
    </Card>
  );
}
