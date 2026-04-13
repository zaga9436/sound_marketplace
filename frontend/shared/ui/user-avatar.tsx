"use client";

import { cn } from "@/shared/utils/cn";

type UserAvatarProps = {
  name?: string | null;
  email?: string | null;
  avatarUrl?: string | null;
  className?: string;
  imageClassName?: string;
};

function getInitials(name?: string | null, email?: string | null) {
  const source = (name || email || "SoundMarket").trim();
  if (!source) return "SM";
  const parts = source.split(/\s+/).filter(Boolean);
  if (parts.length === 1) {
    return parts[0].slice(0, 2).toUpperCase();
  }
  return `${parts[0][0] ?? ""}${parts[1][0] ?? ""}`.toUpperCase();
}

export function UserAvatar({ name, email, avatarUrl, className, imageClassName }: UserAvatarProps) {
  const initials = getInitials(name, email);

  return (
    <div
      className={cn(
        "relative flex h-12 w-12 shrink-0 items-center justify-center overflow-hidden rounded-2xl border border-slate-200 bg-gradient-to-br from-slate-900 via-slate-800 to-slate-700 text-sm font-semibold text-white shadow-sm",
        className
      )}
    >
      {avatarUrl ? (
        <img
          src={avatarUrl}
          alt={name ? `Аватар ${name}` : "Аватар пользователя"}
          className={cn("h-full w-full object-cover", imageClassName)}
        />
      ) : (
        <span>{initials}</span>
      )}
    </div>
  );
}
