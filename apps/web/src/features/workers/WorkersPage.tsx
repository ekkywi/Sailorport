import { useCallback, useEffect, useState, type FormEvent } from "react";
import { RefreshCw, Server } from "lucide-react";
import {
  DataPanel,
  EmptyState,
  ErrorBanner,
  StatusBadge,
  Toolbar,
  formatAbsoluteTime,
  formatRelativeTime,
  skeletonClass,
  useToast,
} from "@/components/app";
import {
  AlertDialog,
  AlertDialogClose,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { me } from "@/features/auth/api";
import { isAdmin } from "@/lib/rbac";
import {
  decommissionWorker,
  listWorkers,
  restoreWorker,
  updateWorkerLabels,
} from "./api";
import {
  formatWorkerEnvironments,
  workerExtraLabels,
  workerLabelString,
  workerTier,
} from "./labels";
import type { Worker } from "./types";

function LabelChips({
  labels,
}: {
  labels: { key: string; value: string }[];
}) {
  if (labels.length === 0) {
    return <span className="text-muted-foreground">—</span>;
  }
  return (
    <div className="flex flex-wrap gap-1">
      {labels.map(({ key, value }) => (
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

function TierBadge({ tier }: { tier: string }) {
  if (!tier) {
    return <span className="text-muted-foreground">—</span>;
  }
  const tone =
    tier === "prod"
      ? "bg-rose-500/12 text-rose-700 dark:text-rose-400"
      : tier === "nonprod"
        ? "bg-sky-500/12 text-sky-700 dark:text-sky-400"
        : "bg-muted text-muted-foreground";
  return (
    <span
      className={cn(
        "inline-flex rounded-full px-2 py-0.5 text-[11px] font-medium capitalize",
        tone,
      )}
    >
      {tier}
    </span>
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
  const { toast } = useToast();
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [isAdminUser, setIsAdminUser] = useState(false);
  const [busyId, setBusyId] = useState<string | null>(null);
  const [editTarget, setEditTarget] = useState<Worker | null>(null);
  const [editTier, setEditTier] = useState("");
  const [editEnvs, setEditEnvs] = useState("");
  const [decommissionTarget, setDecommissionTarget] = useState<Worker | null>(
    null,
  );

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const [meUser, list] = await Promise.all([me(), listWorkers()]);
      setIsAdminUser(isAdmin(meUser.role));
      setWorkers(list);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load workers");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  function openEdit(w: Worker) {
    setEditTarget(w);
    setEditTier(workerTier(w));
    setEditEnvs(workerLabelString(w.labels, "environments"));
  }

  async function onSaveLabels(e: FormEvent) {
    e.preventDefault();
    if (!editTarget) return;
    setBusyId(editTarget.id);
    setError("");
    try {
      const updated = await updateWorkerLabels(editTarget.id, {
        tier: editTier.trim(),
        environments: editEnvs.trim(),
      });
      setWorkers((prev) => prev.map((w) => (w.id === updated.id ? updated : w)));
      setEditTarget(null);
      toast(`Labels updated for ${updated.name}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update labels");
    } finally {
      setBusyId(null);
    }
  }

  async function onDecommission() {
    if (!decommissionTarget) return;
    const id = decommissionTarget.id;
    setBusyId(id);
    setError("");
    try {
      const updated = await decommissionWorker(id);
      setWorkers((prev) => prev.map((w) => (w.id === updated.id ? updated : w)));
      setDecommissionTarget(null);
      toast(`${updated.name} decommissioned (draining)`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to decommission");
    } finally {
      setBusyId(null);
    }
  }

  async function onRestore(w: Worker) {
    setBusyId(w.id);
    setError("");
    try {
      const updated = await restoreWorker(w.id);
      setWorkers((prev) => prev.map((x) => (x.id === updated.id ? updated : x)));
      toast(`${updated.name} restored (offline until next heartbeat)`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to restore");
    } finally {
      setBusyId(null);
    }
  }

  const online = workers.filter((w) => w.status === "online").length;
  const draining = workers.filter((w) => w.status === "draining").length;
  const meta =
    loading && workers.length === 0
      ? "Loading…"
      : workers.length > 0
        ? `${workers.length} worker${workers.length === 1 ? "" : "s"} · ${online} online${
            draining > 0 ? ` · ${draining} draining` : ""
          }`
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
            <table className="w-full min-w-[720px] text-left text-[13px]">
              <thead>
                <tr className="border-b border-border bg-muted/40 text-[11px] font-medium tracking-[0.02em] text-muted-foreground uppercase">
                  <th className="px-4 py-2.5 font-medium">Name</th>
                  <th className="px-4 py-2.5 font-medium">Host</th>
                  <th className="px-4 py-2.5 font-medium">Status</th>
                  <th className="px-4 py-2.5 font-medium">Tier</th>
                  <th className="px-4 py-2.5 font-medium">Environments</th>
                  <th className="px-4 py-2.5 font-medium">Last seen</th>
                  <th className="px-4 py-2.5 font-medium">Other labels</th>
                  {isAdminUser ? (
                    <th className="px-4 py-2.5 font-medium">Actions</th>
                  ) : null}
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
                    <td className="px-4 py-2.5">
                      <TierBadge tier={workerTier(w)} />
                    </td>
                    <td className="px-4 py-2.5 font-mono text-[11px] text-muted-foreground uppercase">
                      {formatWorkerEnvironments(w)}
                    </td>
                    <td
                      className="px-4 py-2.5 tabular-nums text-muted-foreground"
                      title={formatAbsoluteTime(w.last_seen_at)}
                    >
                      {formatRelativeTime(w.last_seen_at)}
                    </td>
                    <td className="px-4 py-2.5">
                      <LabelChips labels={workerExtraLabels(w.labels)} />
                    </td>
                    {isAdminUser ? (
                      <td className="px-4 py-2.5">
                        <div className="flex flex-wrap gap-1.5">
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            className="h-7 text-[12px]"
                            disabled={busyId === w.id}
                            onClick={() => openEdit(w)}
                          >
                            Edit labels
                          </Button>
                          {w.status === "draining" ? (
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              className="h-7 text-[12px]"
                              disabled={busyId === w.id}
                              onClick={() => void onRestore(w)}
                            >
                              Restore
                            </Button>
                          ) : (
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              className="h-7 text-[12px] text-amber-700 dark:text-amber-400"
                              disabled={busyId === w.id}
                              onClick={() => setDecommissionTarget(w)}
                            >
                              Decommission
                            </Button>
                          )}
                        </div>
                      </td>
                    ) : null}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : null}
      </DataPanel>

      <Dialog
        open={editTarget != null}
        onOpenChange={(open) => {
          if (!open) setEditTarget(null);
        }}
      >
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Edit labels</DialogTitle>
            <DialogDescription>
              Soft override for {editTarget?.name}. Agent re-register may
              overwrite these values.
            </DialogDescription>
          </DialogHeader>
          <form className="space-y-3" onSubmit={(e) => void onSaveLabels(e)}>
            <div className="space-y-1.5">
              <Label htmlFor="worker-tier">Tier</Label>
              <Input
                id="worker-tier"
                value={editTier}
                onChange={(e) => setEditTier(e.target.value)}
                placeholder="nonprod"
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="worker-envs">Environments</Label>
              <Input
                id="worker-envs"
                value={editEnvs}
                onChange={(e) => setEditEnvs(e.target.value)}
                placeholder="dev,staging"
              />
              <p className="text-[11px] text-muted-foreground">
                Comma-separated slugs (same as agent ENVIRONMENTS). Empty allows
                all.
              </p>
            </div>
            <div className="flex justify-end gap-2 pt-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => setEditTarget(null)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={busyId === editTarget?.id}>
                Save
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <AlertDialog
        open={decommissionTarget != null}
        onOpenChange={(open) => {
          if (!open && busyId === null) setDecommissionTarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Decommission worker?</AlertDialogTitle>
            <AlertDialogDescription>
              <span className="font-medium text-foreground">
                {decommissionTarget?.name}
              </span>{" "}
              will be set to draining and cannot be targeted for new deploys
              until restored. Heartbeats will not bring it back online.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogClose
              render={
                <Button
                  type="button"
                  variant="outline"
                  disabled={busyId === decommissionTarget?.id}
                />
              }
            >
              Cancel
            </AlertDialogClose>
            <Button
              type="button"
              variant="destructive"
              disabled={busyId === decommissionTarget?.id}
              onClick={() => void onDecommission()}
            >
              Decommission
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
