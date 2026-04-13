"use client";

import { useQuery } from "@tanstack/react-query";

import { bidsApi } from "@/entities/bid/api/bids";
import { profilesApi } from "@/entities/profile/api/profiles";
import { CreateBidForm } from "@/features/bid/create-bid-form";
import { CreateOrderFromBidButton, CreateOrderFromOfferButton } from "@/features/order/create-order-actions";
import { getErrorMessage } from "@/lib/api/errors";
import { useAuthStore } from "@/lib/auth/session-store";
import { Badge } from "@/shared/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";
import { Bid, Card as MarketplaceCard } from "@/shared/types/api";

function formatPrice(value: number) {
  return new Intl.NumberFormat("ru-RU", {
    style: "currency",
    currency: "RUB",
    maximumFractionDigits: 0
  }).format(value);
}

export function CardDealPanel({ card }: { card: MarketplaceCard }) {
  const user = useAuthStore((state) => state.user);
  const isOwner = user?.id === card.author_id;
  const canViewBids = card.card_type === "request" && Boolean(user && (isOwner || user.role === "admin"));
  const canCreateBid = card.card_type === "request" && user?.role === "engineer" && !isOwner;
  const canCreateOrderFromOffer = card.card_type === "offer" && user?.role === "customer" && !isOwner;

  const bidsQuery = useQuery({
    queryKey: ["bids", card.id],
    queryFn: () => bidsApi.listByRequest(card.id),
    enabled: canViewBids
  });

  return (
    <div className="space-y-4">
      {card.card_type === "request" ? (
        <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
          <CardHeader className="space-y-3">
            <Badge className="bg-slate-900/90 text-white" variant="secondary">
              Отклики
            </Badge>
            <CardTitle>{canViewBids ? "Заявки по этому запросу" : "Отклик на запрос"}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-5">
            {canCreateBid ? (
              <CreateBidForm requestId={card.id} />
            ) : canViewBids ? (
              bidsQuery.isLoading ? (
                <div className="space-y-3">
                  {Array.from({ length: 2 }).map((_, index) => (
                    <div key={index} className="h-28 animate-pulse rounded-2xl border border-slate-200 bg-slate-50" />
                  ))}
                </div>
              ) : bidsQuery.isError ? (
                <p className="text-sm text-red-600">{getErrorMessage(bidsQuery.error)}</p>
              ) : bidsQuery.data && bidsQuery.data.length > 0 ? (
                <div className="space-y-3">
                  {bidsQuery.data.map((bid) => (
                    <BidRow key={bid.id} bid={bid} canCreateOrder={user?.role === "customer" && isOwner} />
                  ))}
                </div>
              ) : (
                <div className="rounded-2xl border border-dashed border-slate-300 bg-slate-50 px-4 py-5 text-sm text-slate-500">
                  Пока откликов нет. Они появятся здесь, когда инженеры начнут отвечать на запрос.
                </div>
              )
            ) : (
              <div className="rounded-2xl border border-dashed border-slate-300 bg-slate-50 px-4 py-5 text-sm text-slate-500">
                {user
                  ? "Эта зона доступна инженеру для отклика или владельцу карточки для просмотра заявок."
                  : "Войдите в систему, чтобы откликнуться на запрос или управлять своими заявками."}
              </div>
            )}
          </CardContent>
        </Card>
      ) : (
        <Card className="border-slate-200/80 bg-white/95 shadow-[0_20px_60px_-32px_rgba(15,23,42,0.24)]">
          <CardHeader className="space-y-3">
            <Badge className="bg-slate-900/90 text-white" variant="secondary">
              Заказ
            </Badge>
            <CardTitle>Покупка предложения</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-slate-600">
              Стоимость заказа: <span className="font-semibold text-slate-950">{formatPrice(card.price)}</span>
            </p>

            {canCreateOrderFromOffer ? (
              <CreateOrderFromOfferButton cardId={card.id} />
            ) : (
              <div className="rounded-2xl border border-dashed border-slate-300 bg-slate-50 px-4 py-5 text-sm text-slate-500">
                {isOwner
                  ? "Это ваша карточка. Заказы по offer создают заказчики."
                  : user
                    ? "Заказ по offer может создать только заказчик."
                    : "Войдите как заказчик, чтобы создать заказ по этому offer."}
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}

function BidRow({ bid, canCreateOrder }: { bid: Bid; canCreateOrder: boolean }) {
  const engineerProfileQuery = useQuery({
    queryKey: ["profile", bid.engineer_id],
    queryFn: () => profilesApi.getById(bid.engineer_id)
  });

  return (
    <div className="rounded-2xl border border-slate-200 bg-slate-50/80 p-4">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div className="space-y-2">
          <div className="flex flex-wrap items-center gap-2">
            <Badge variant="outline">{formatPrice(bid.price)}</Badge>
            <Badge variant="secondary">{new Date(bid.created_at).toLocaleDateString("ru-RU")}</Badge>
          </div>
          <p className="font-medium text-slate-950">{engineerProfileQuery.data?.display_name ?? "Инженер"}</p>
          <p className="text-sm leading-6 text-slate-600">{bid.message}</p>
        </div>
        {canCreateOrder ? <CreateOrderFromBidButton bidId={bid.id} /> : null}
      </div>
    </div>
  );
}
