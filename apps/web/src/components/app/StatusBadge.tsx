import { cn } from "@/lib/utils";

type WorkerStatus = "online" | "offline" | "draining" | string;

export function StatusDot({
  status,
  className,
}: {
  status: WorkerStatus;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-block size-1.5 shrink-0 rounded-full",
        status === "online" && "bg-emerald-500",
        status === "draining" && "bg-amber-500",
        status !== "online" && status !== "draining" && "bg-muted-foreground/40",
        className,
      )}
      aria-hidden
    />
  );
}

export function StatusBadge({
  status,
  className,
}: {
  status: WorkerStatus;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-[11px] font-medium capitalize",
        status === "online" &&
          "bg-emerald-500/12 text-emerald-700 dark:text-emerald-400",
        status === "draining" &&
          "bg-amber-500/12 text-amber-700 dark:text-amber-400",
        status !== "online" &&
          status !== "draining" &&
          "bg-muted text-muted-foreground",
        className,
      )}
    >
      <StatusDot status={status} />
      {status}
    </span>
  );
}
