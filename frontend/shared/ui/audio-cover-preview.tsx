"use client";

import { MouseEvent, useEffect, useRef, useState } from "react";
import { Pause, Play, Waves } from "lucide-react";

import { cn } from "@/shared/utils/cn";

type AudioCoverPreviewProps = {
  coverUrl?: string;
  audioUrl?: string;
  title: string;
  compact?: boolean;
  className?: string;
};

export function AudioCoverPreview({ coverUrl, audioUrl, title, compact = false, className }: AudioCoverPreviewProps) {
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);

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

  const togglePlayback = async (event: MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    event.stopPropagation();

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
        coverUrl
          ? {
              backgroundImage: `linear-gradient(180deg, rgba(15,23,42,0.18), rgba(15,23,42,0.82)), url(${coverUrl})`,
              backgroundSize: "cover",
              backgroundPosition: "center"
            }
          : undefined
      }
    >
      {!coverUrl ? (
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(148,163,184,0.24),transparent_28%),radial-gradient(circle_at_bottom_right,rgba(71,85,105,0.22),transparent_22%)]" />
      ) : null}

      <div className="relative flex h-full flex-col justify-between p-5">
        <div className="flex items-center justify-between gap-3">
          <div className="inline-flex items-center gap-2 rounded-full bg-white/12 px-3 py-1 text-xs font-medium backdrop-blur">
            <Waves className="h-3.5 w-3.5" />
            Аудио preview
          </div>
          {coverUrl ? <div className="rounded-full bg-black/25 px-3 py-1 text-xs backdrop-blur">С обложкой</div> : null}
        </div>

        <div className={compact ? "space-y-1" : "space-y-2"}>
          <button
            type="button"
            onClick={togglePlayback}
            onMouseDown={(event) => event.stopPropagation()}
            disabled={!audioUrl}
            aria-label={isPlaying ? "Поставить preview на паузу" : "Прослушать preview"}
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

      {audioUrl ? <audio ref={audioRef} preload="none" src={audioUrl} onClick={(event) => event.stopPropagation()} /> : null}
    </div>
  );
}
