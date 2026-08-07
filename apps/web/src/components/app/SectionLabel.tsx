import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

export function SectionLabel({
  children,
  className,
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <p
      className={cn(
        "px-2 text-[11px] font-medium tracking-[0.04em] text-muted-foreground uppercase",
        className,
      )}
    >
      {children}
    </p>
  );
}
