import { useCallback, useEffect, useRef, useState } from "react";
import { ExternalLink, RefreshCw } from "lucide-react";
import {
  EmptyState,
  ErrorBanner,
  EnvironmentBadge,
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
  onRefreshCatalog?: () => void;
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
  onRefreshCatalog,
}: DeploymentsDialogProps) {
  const [items, setItems] = useState<Deployment[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [loadedOnce, setLoadedOnce] = useState(false);

  // Avoid putting parent inline callbacks in load deps (causes refresh loops / flicker).
  const onRefreshCatalogRef = useRef(onRefreshCatalog);
  onRefreshCatalogRef.current = onRefreshCatalog;

  const load = useCallback(async (options?: { silent?: boolean }) => {
    if (!serviceId) return;
    const silent = options?.silent ?? false;
    if (!silent) {
      setLoading(true);
    }
    setError("");
    try {
      const next = await listDeploymentsByService(serviceId);
      setItems(next);
      setLoadedOnce(true);
      if (next.some((d) => isActiveStatus(d.status))) {
        onRefreshCatalogRef.current?.();
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load deployments");
    } finally {
      if (!silent) {
        setLoading(false);
      }
    }
  }, [serviceId]);

  useEffect(() => {
    if (!open || !serviceId) {
      return;
    }
    setItems([]);
    setLoadedOnce(false);
    setError("");
    void load({ silent: false });
  }, [open, serviceId, load]);

  useEffect(() => {
    if (!open || !serviceId) return;
    const hasActive = items.some((d) => isActiveStatus(d.status));
    if (!hasActive) return;

    const id = window.setInterval(() => {
      void load({ silent: true });
    }, 3000);
    return () => window.clearInterval(id);
  }, [open, serviceId, items, load]);

  const showEmpty = !error && loadedOnce && !loading && items.length === 0;
  const showList = items.length > 0;

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
            onClick={() => void load({ silent: false })}
            disabled={loading}
          >
            <RefreshCw className={cn("size-3.5", loading && "animate-spin")} />
            Refresh
          </Button>
        </div>

        {error ? (
          <ErrorBanner
            message={error}
            onRetry={() => void load({ silent: false })}
          />
        ) : null}

        {loading && !loadedOnce ? (
          <div className="space-y-0 divide-y divide-border rounded-lg border border-border py-1">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="space-y-2 px-3 py-3">
                <div className="h-4 w-24 animate-pulse rounded bg-muted" />
                <div className="h-3 w-40 animate-pulse rounded bg-muted" />
              </div>
            ))}
          </div>
        ) : null}

        {showEmpty ? (
          <EmptyState
            title="No deployments yet"
            description="Click Deploy on this service to create the first one."
            className="py-10"
          />
        ) : null}

        {showList ? (
          <ul className="max-h-[360px] space-y-0 divide-y divide-border overflow-y-auto rounded-lg border border-border">
            {items.map((d) => (
              <li key={d.id} className="px-3 py-3">
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                      <EnvironmentBadge slug={d.environment_slug} />
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
