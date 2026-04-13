"use client";

import { useEffect, useState } from "react";

const PREFIX = "sm-card-cover:";
const EVENT_NAME = "soundmarket:card-cover-updated";

export type CardCoverRecord = {
  dataUrl: string;
  fileName: string;
  updatedAt: string;
};

function key(cardId: string) {
  return `${PREFIX}${cardId}`;
}

export function getCardCover(cardId: string): CardCoverRecord | null {
  if (typeof window === "undefined") return null;
  const raw = window.localStorage.getItem(key(cardId));
  if (!raw) return null;

  try {
    return JSON.parse(raw) as CardCoverRecord;
  } catch {
    return null;
  }
}

export function saveCardCover(cardId: string, value: CardCoverRecord) {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(key(cardId), JSON.stringify(value));
  window.dispatchEvent(new CustomEvent(EVENT_NAME, { detail: { cardId } }));
}

export function removeCardCover(cardId: string) {
  if (typeof window === "undefined") return;
  window.localStorage.removeItem(key(cardId));
  window.dispatchEvent(new CustomEvent(EVENT_NAME, { detail: { cardId } }));
}

export function useCardCover(cardId?: string) {
  const [cover, setCover] = useState<CardCoverRecord | null>(null);

  useEffect(() => {
    if (!cardId) {
      setCover(null);
      return;
    }

    const sync = () => setCover(getCardCover(cardId));
    sync();

    const handler = (event: Event) => {
      const customEvent = event as CustomEvent<{ cardId?: string }>;
      if (!customEvent.detail?.cardId || customEvent.detail.cardId === cardId) {
        sync();
      }
    };

    window.addEventListener(EVENT_NAME, handler);
    window.addEventListener("storage", sync);

    return () => {
      window.removeEventListener(EVENT_NAME, handler);
      window.removeEventListener("storage", sync);
    };
  }, [cardId]);

  return cover;
}
