import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

export function Toolbar({
  meta,
  actions,
  className,
}: {
  /** Compact count / status line — page title lives in the topbar */
  meta?: ReactNode;
  actions?: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between",
        className,
      )}
    >
      {meta ? (
        <p className="min-w-0 text-[12px] text-muted-foreground">{meta}</p>
      ) : (
        <span className="hidden sm:block" />
      )}
      {actions ? (
        <div className="flex shrink-0 flex-wrap items-center gap-1.5 sm:ml-auto">
          {actions}
        </div>
      ) : null}
    </div>
  );
}
