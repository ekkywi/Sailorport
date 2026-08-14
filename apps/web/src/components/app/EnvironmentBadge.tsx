import { cn } from "@/lib/utils";

export function EnvironmentBadge({
  slug,
  className,
}: {
  slug: string;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex shrink-0 rounded-md border border-border bg-muted/40 px-1.5 py-0.5 font-mono text-[10px] font-medium tracking-wide text-muted-foreground uppercase",
        className,
      )}
    >
      {slug}
    </span>
  );
}
