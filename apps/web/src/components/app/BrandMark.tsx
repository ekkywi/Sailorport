import { cn } from "@/lib/utils";

export function BrandMark({ className }: { className?: string }) {
  return (
    <span
      className={cn(
        "flex size-[18px] shrink-0 items-center justify-center rounded-[4px] bg-primary text-[10px] font-semibold leading-none tracking-tight text-primary-foreground",
        className,
      )}
    >
      S
    </span>
  );
}
