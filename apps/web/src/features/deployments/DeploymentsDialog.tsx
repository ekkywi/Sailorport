import { useCallback, useEffect, useState } from "react";
import { ExternalLink, RefreshCw } from "lucide-react";
import {
  EmptyState,
  ErrorBanner,
  StatusBadge,
  formatRelativeTime,
} from "@/components/app";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import { listDeploymentsByService } from "./api";
import type { Deployment } from "./types";

type DeploymentsDialogProps = {
  serviceId: string | null;
  serviceName: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

function isActiveStatus(status: string) {
  return status === "pending" || status === "claimed" || status === "building";
}

function deployBadgeClass(status: string) {
  if (status === "running") {
    return "bg-emerald-500/12 text-emerald-700 dark:text-emerald-400";
  }
  if (status === "failed") {
    return "bg-red-500/12 text-red-700 dark:text-red-400";
  }
  if (status === "building" || status === "claimed") {
    return "bg-amber-500/12 text-amber-700 dark:text-amber-400";
  }
  return undefined;
}

export function DeploymentsDialog({
  serviceId,
  serviceName,
  open,
  onOpenChange,
}: DeploymentsDialogProps) {
  const [items, setItems] = useState<Deployment[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const load = useCallback(async () => {
    if (!serviceId) return;
    setLoading(true);
    setError("");
    try {
      setItems(await listDeploymentsByService(serviceId));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load deployments");
    } finally {
      setLoading(false);
    }
  }, [serviceId]);

  useEffect(() => {
    if (open && serviceId) {
      void load();
    }
  }, [open, serviceId, load]);

  useEffect(() => {
    if (!open || !serviceId) return;
    const hasActive = items.some((d) => isActiveStatus(d.status));
    if (!hasActive) return;

    const id = window.setInterval(() => {
      void load();
    }, 3000);
    return () => window.clearInterval(id);
  }, [open, serviceId, items, load]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Deployments</DialogTitle>
          <DialogDescription>
            Recent deploys for{" "}
            <span className="font-medium text-foreground">{serviceName}</span>.
          </DialogDescription>
        </DialogHeader>

        <div className="flex justify-end">
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 gap-1.5 text-[13px]"
            onClick={() => void load()}
            disabled={loading}
          >
            <RefreshCw className={cn("size-3.5", loading && "animate-spin")} />
            Refresh
          </Button>
        </div>

        {error ? (
          <ErrorBanner message={error} onRetry={() => void load()} />
        ) : null}

        {!error && items.length === 0 && !loading ? (
          <EmptyState
            title="No deployments yet"
            description="Click Deploy on this service to create the first one."
            className="py-10"
          />
        ) : null}

        {items.length > 0 ? (
          <ul className="max-h-[360px] space-y-0 divide-y divide-border overflow-y-auto rounded-lg border border-border">
            {items.map((d) => (
              <li key={d.id} className="px-3 py-3">
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                      <StatusBadge
                        status={d.status}
                        className={deployBadgeClass(d.status)}
                      />
                      <span className="font-mono text-[11px] text-muted-foreground">
                        {d.id.slice(0, 8)}
                      </span>
                    </div>
                    <p className="mt-1 text-[12px] text-muted-foreground">
                      {formatRelativeTime(d.created_at)}
                      {d.image_tag ? ` · ${d.image_tag}` : ""}
                    </p>
                    {d.error_message ? (
                      <p className="mt-1 line-clamp-2 text-[12px] text-red-600 dark:text-red-400">
                        {d.error_message}
                      </p>
                    ) : null}
                  </div>
                  {d.status === "running" && d.port != null ? (
                    <a
                      href={`http://localhost:${d.port}/healthz`}
                      target="_blank"
                      rel="noreferrer"
                      className="inline-flex shrink-0 items-center gap-1 text-[12px] text-foreground underline-offset-2 hover:underline"
                    >
                      healthz
                      <ExternalLink className="size-3" />
                    </a>
                  ) : null}
                </div>
              </li>
            ))}
          </ul>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}