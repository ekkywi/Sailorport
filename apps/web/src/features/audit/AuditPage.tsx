import { useCallback, useEffect, useState } from "react";
import { ClipboardList, RefreshCw } from "lucide-react";
import {
  DataPanel,
  EmptyState,
  ErrorBanner,
  Toolbar,
  formatAbsoluteTime,
  formatRelativeTime,
  skeletonClass,
} from "@/components/app";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { listAuditEvents } from "./api";
import type { AuditEvent } from "./types";

function actionBadgeClass(action: string) {
  if (action.endsWith(".delete")) {
    return "bg-red-500/12 text-red-700 dark:text-red-400";
  }
  if (action.endsWith(".create")) {
    return "bg-emerald-500/12 text-emerald-700 dark:text-emerald-400";
  }
  if (action.includes("disable") || action.includes("password")) {
    return "bg-amber-500/12 text-amber-700 dark:text-amber-400";
  }
  return "bg-muted text-muted-foreground";
}

function formatPayload(payload: Record<string, unknown>) {
  const keys = Object.keys(payload);
  if (keys.length === 0) return "—";
  try {
    const text = JSON.stringify(payload);
    return text.length > 120 ? `${text.slice(0, 117)}…` : text;
  } catch {
    return "—";
  }
}

function AuditTableSkeleton() {
  return (
    <div className="space-y-0 divide-y divide-border">
      {Array.from({ length: 5 }).map((_, i) => (
        <div key={i} className="flex items-center gap-4 px-4 py-3">
          <div className={skeletonClass("h-3.5 w-24")} />
          <div className={skeletonClass("h-3.5 w-32")} />
          <div className={skeletonClass("h-5 w-28 rounded-full")} />
          <div className={skeletonClass("h-3.5 w-40")} />
        </div>
      ))}
    </div>
  );
}

export function AuditPage() {
  const [events, setEvents] = useState<AuditEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setEvents(await listAuditEvents(50));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load audit log");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  const meta =
    loading && events.length === 0
      ? "Loading…"
      : events.length > 0
        ? `${events.length} event${events.length === 1 ? "" : "s"} (latest 50)`
        : "No audit events yet";

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
        {loading && events.length === 0 ? <AuditTableSkeleton /> : null}

        {!loading && events.length === 0 && !error ? (
          <EmptyState
            icon={ClipboardList}
            title="No audit events yet"
            description="Actions such as catalog changes and user administration will appear here."
            className="py-14"
          />
        ) : null}

        {events.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[720px] text-left text-[13px]">
              <thead>
                <tr className="border-b border-border bg-muted/40 text-[11px] font-medium tracking-[0.02em] text-muted-foreground uppercase">
                  <th className="px-4 py-2.5 font-medium">When</th>
                  <th className="px-4 py-2.5 font-medium">Actor</th>
                  <th className="px-4 py-2.5 font-medium">Action</th>
                  <th className="px-4 py-2.5 font-medium">Resource</th>
                  <th className="px-4 py-2.5 font-medium">Details</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {events.map((ev) => (
                  <tr key={ev.id} className="hover:bg-muted/30">
                    <td
                      className="whitespace-nowrap px-4 py-2.5 tabular-nums text-muted-foreground"
                      title={formatAbsoluteTime(ev.at)}
                    >
                      {formatRelativeTime(ev.at)}
                    </td>
                    <td className="max-w-[180px] truncate px-4 py-2.5 text-muted-foreground">
                      {ev.actor_email || "—"}
                    </td>
                    <td className="px-4 py-2.5">
                      <span
                        className={cn(
                          "inline-flex rounded-md px-1.5 py-0.5 font-mono text-[11px]",
                          actionBadgeClass(ev.action),
                        )}
                      >
                        {ev.action}
                      </span>
                    </td>
                    <td className="px-4 py-2.5">
                      <div className="min-w-0">
                        <p className="truncate font-medium tracking-[-0.01em]">
                          {ev.resource_name || "—"}
                        </p>
                        <p className="truncate font-mono text-[11px] text-muted-foreground">
                          {ev.resource_type}
                          {ev.resource_id ? ` · ${ev.resource_id.slice(0, 8)}` : ""}
                        </p>
                      </div>
                    </td>
                    <td
                      className="max-w-[280px] truncate px-4 py-2.5 font-mono text-[11px] text-muted-foreground"
                      title={formatPayload(ev.payload)}
                    >
                      {formatPayload(ev.payload)}
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
