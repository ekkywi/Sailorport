import { useCallback, useEffect, useState } from "react";
import { RefreshCw, Server } from "lucide-react";
import {
  DataPanel,
  EmptyState,
  ErrorBanner,
  StatusBadge,
  Toolbar,
  formatAbsoluteTime,
  formatRelativeTime,
  labelEntries,
  skeletonClass,
} from "@/components/app";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { listWorkers } from "./api";
import type { Worker } from "./types";

function LabelChips({ labels }: { labels: Worker["labels"] }) {
  const entries = labelEntries(labels);
  if (entries.length === 0) {
    return <span className="text-muted-foreground">—</span>;
  }
  return (
    <div className="flex flex-wrap gap-1">
      {entries.map(({ key, value }) => (
        <span
          key={key}
          className="inline-flex max-w-[140px] truncate rounded bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground"
          title={`${key}=${value}`}
        >
          {value ? `${key}=${value}` : key}
        </span>
      ))}
    </div>
  );
}

function WorkersTableSkeleton() {
  return (
    <div className="space-y-0 divide-y divide-border">
      {Array.from({ length: 4 }).map((_, i) => (
        <div key={i} className="flex items-center gap-4 px-4 py-3">
          <div className={skeletonClass("h-3.5 w-28")} />
          <div className={skeletonClass("h-3.5 w-24")} />
          <div className={skeletonClass("h-5 w-16 rounded-full")} />
          <div className={skeletonClass("h-3.5 w-20")} />
        </div>
      ))}
    </div>
  );
}

export function WorkersPage() {
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setWorkers(await listWorkers());
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load workers");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  const online = workers.filter((w) => w.status === "online").length;
  const meta =
    loading && workers.length === 0
      ? "Loading…"
      : workers.length > 0
        ? `${workers.length} worker${workers.length === 1 ? "" : "s"} · ${online} online`
        : "No workers registered";

  return (
    <div className="space-y-4">
      <Toolbar
        meta={meta}
        actions={
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
        }
      />

      {error ? (
        <ErrorBanner message={error} onRetry={() => void load()} />
      ) : null}

      <DataPanel>
        {loading && workers.length === 0 ? <WorkersTableSkeleton /> : null}

        {!loading && workers.length === 0 && !error ? (
          <EmptyState
            icon={Server}
            title="No workers yet"
            description="Start a Sailorport agent on a node to register it and begin heartbeats."
            className="py-14"
          />
        ) : null}

        {workers.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[640px] text-left text-[13px]">
              <thead>
                <tr className="border-b border-border bg-muted/40 text-[11px] font-medium tracking-[0.02em] text-muted-foreground uppercase">
                  <th className="px-4 py-2.5 font-medium">Name</th>
                  <th className="px-4 py-2.5 font-medium">Host</th>
                  <th className="px-4 py-2.5 font-medium">Status</th>
                  <th className="px-4 py-2.5 font-medium">Last seen</th>
                  <th className="px-4 py-2.5 font-medium">Labels</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {workers.map((w) => (
                  <tr key={w.id} className="hover:bg-muted/30">
                    <td className="px-4 py-2.5 font-medium tracking-[-0.01em]">
                      {w.name}
                    </td>
                    <td className="px-4 py-2.5 text-muted-foreground">
                      {w.hostname || "—"}
                    </td>
                    <td className="px-4 py-2.5">
                      <StatusBadge status={w.status} />
                    </td>
                    <td
                      className="px-4 py-2.5 tabular-nums text-muted-foreground"
                      title={formatAbsoluteTime(w.last_seen_at)}
                    >
                      {formatRelativeTime(w.last_seen_at)}
                    </td>
                    <td className="px-4 py-2.5">
                      <LabelChips labels={w.labels} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : null}
      </DataPanel>
    </div>
  );
}
