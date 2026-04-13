"use client";

import { useEffect, useRef, useState } from "react";
import { Pause, Play, Waves } from "lucide-react";

import { useCardCover } from "@/lib/cards/card-cover-store";
import { cn } from "@/shared/utils/cn";

type AudioCoverPreviewProps = {
  cardId: string;
  audioUrl?: string;
  title: string;
  compact?: boolean;
  className?: string;
};

export function AudioCoverPreview({ cardId, audioUrl, title, compact = false, className }: AudioCoverPreviewProps) {
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const cover = useCardCover(cardId);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const onPause = () => setIsPlaying(false);
    const onEnded = () => setIsPlaying(false);

    audio.addEventListener("pause", onPause);
    audio.addEventListener("ended", onEnded);

    return () => {
      audio.removeEventListener("pause", onPause);
      audio.removeEventListener("ended", onEnded);
    };
  }, []);

  const togglePlayback = async () => {
    const audio = audioRef.current;
    if (!audio || !audioUrl) return;

    if (audio.paused) {
      await audio.play();
      setIsPlaying(true);
      return;
    }

    audio.pause();
    setIsPlaying(false);
  };

  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-[1.75rem] border border-slate-200 bg-[linear-gradient(145deg,rgba(15,23,42,0.92),rgba(51,65,85,0.92))] text-white shadow-[0_20px_60px_-30px_rgba(15,23,42,0.55)]",
        compact ? "aspect-[16/11]" : "aspect-[16/10]",
        className
      )}
      style={
        cover
          ? {
              backgroundImage: `linear-gradient(180deg, rgba(15,23,42,0.18), rgba(15,23,42,0.82)), url(${cover.dataUrl})`,
              backgroundSize: "cover",
              backgroundPosition: "center"
            }
          : undefined
      }
    >
      {!cover ? (
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(148,163,184,0.24),transparent_28%),radial-gradient(circle_at_bottom_right,rgba(56,189,248,0.18),transparent_22%)]" />
      ) : null}

      <div className="relative flex h-full flex-col justify-between p-5">
        <div className="flex items-center justify-between gap-3">
          <div className="inline-flex items-center gap-2 rounded-full bg-white/12 px-3 py-1 text-xs font-medium backdrop-blur">
            <Waves className="h-3.5 w-3.5" />
            Аудио preview
          </div>
          {cover ? <div className="rounded-full bg-black/25 px-3 py-1 text-xs backdrop-blur">С обложкой</div> : null}
        </div>

        <div className="space-y-3">
          <button
            type="button"
            onClick={togglePlayback}
            disabled={!audioUrl}
            className="flex h-14 w-14 items-center justify-center rounded-full bg-white text-slate-950 shadow-lg transition hover:scale-[1.02] disabled:cursor-not-allowed disabled:bg-white/70"
          >
            {isPlaying ? <Pause className="h-5 w-5" /> : <Play className="ml-0.5 h-5 w-5" />}
          </button>
          <div className={compact ? "space-y-1" : "space-y-2"}>
            <p className={cn("max-w-[85%] font-semibold leading-tight", compact ? "line-clamp-2 text-base" : "line-clamp-2 text-2xl")}>{title}</p>
            <p className="text-sm text-white/80">
              {audioUrl ? "Нажмите, чтобы прослушать фрагмент." : "Загрузите аудио preview, чтобы оживить карточку."}
            </p>
          </div>
        </div>
      </div>

      {audioUrl ? <audio ref={audioRef} preload="none" src={audioUrl} /> : null}
    </div>
  );
}
