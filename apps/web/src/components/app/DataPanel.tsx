import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

export function DataPanel({
  children,
  className,
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "overflow-hidden rounded-lg border border-border bg-card/80",
        className,
      )}
    >
      {children}
    </div>
  );
}
