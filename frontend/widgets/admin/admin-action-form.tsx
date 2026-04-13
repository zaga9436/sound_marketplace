"use client";

import { useState } from "react";

import { Button } from "@/shared/ui/button";
import { Textarea } from "@/shared/ui/textarea";

export function AdminActionForm({
  actionLabel,
  confirmLabel,
  placeholder,
  optional = false,
  disabled = false,
  pending = false,
  onSubmit
}: {
  actionLabel: string;
  confirmLabel: string;
  placeholder: string;
  optional?: boolean;
  disabled?: boolean;
  pending?: boolean;
  onSubmit: (reason: string) => void;
}) {
  const [open, setOpen] = useState(false);
  const [reason, setReason] = useState("");

  return (
    <div className="space-y-3">
      {!open ? (
        <Button type="button" variant="outline" className="rounded-2xl border-slate-300 bg-white text-slate-900 hover:bg-slate-100" disabled={disabled} onClick={() => setOpen(true)}>
          {actionLabel}
        </Button>
      ) : (
        <div className="rounded-2xl border border-slate-200 bg-slate-50 p-3">
          <div className="space-y-3">
            <Textarea
              value={reason}
              onChange={(event) => setReason(event.target.value)}
              placeholder={placeholder}
              className="min-h-[100px] rounded-2xl border-slate-300 bg-white"
            />
            <div className="flex flex-wrap gap-2">
              <Button
                type="button"
                className="rounded-2xl bg-slate-950 text-white hover:bg-slate-800"
                disabled={pending || (!optional && !reason.trim())}
                onClick={() => {
                  onSubmit(reason.trim());
                  setOpen(false);
                  setReason("");
                }}
              >
                {pending ? "Сохраняем..." : confirmLabel}
              </Button>
              <Button
                type="button"
                variant="ghost"
                className="rounded-2xl text-slate-700 hover:bg-slate-100"
                onClick={() => {
                  setOpen(false);
                  setReason("");
                }}
              >
                Отмена
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
