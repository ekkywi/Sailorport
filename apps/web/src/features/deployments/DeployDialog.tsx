import { useEffect, useMemo, useState } from "react";
import { Check, Rocket, Search, Server } from "lucide-react";
import { ErrorBanner, StatusBadge } from "@/components/app";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import type { Service } from "../catalog/types";
import { listEnvironments } from "../environments/api";
import type { Environment } from "../environments/types";
import { listWorkers } from "../workers/api";
import { workerAllowsEnvironment, formatWorkerEnvironments } from "../workers/labels";
import type { Worker } from "../workers/types";
import { createDeployment } from "./api";

type DeployDialogProps = {
  service: Service | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onDeployed?: (service: Service, environment: string) => void;
};

const ANY_WORKER = "";
const WORKER_SEARCH_THRESHOLD = 4;
const WORKER_LIST_MAX_HEIGHT = "max-h-52";

function defaultWorkerForEnv(
  service: Service | null,
  envSlug: string,
  eligibleWorkers: Worker[],
): string {
  const prior = service?.env_deployments?.[envSlug]?.worker_id?.trim();
  if (prior && eligibleWorkers.some((w) => w.id === prior)) {
    return prior;
  }
  return ANY_WORKER;
}

function workerStatusClass(status: string) {
  if (status === "online") {
    return "bg-emerald-500/12 text-emerald-700 dark:text-emerald-400";
  }
  if (status === "draining") {
    return "bg-amber-500/12 text-amber-700 dark:text-amber-400";
  }
  return "bg-muted text-muted-foreground";
}

function matchesWorkerSearch(worker: Worker, query: string): boolean {
  const q = query.trim().toLowerCase();
  if (!q) return true;
  return (
    worker.name.toLowerCase().includes(q) ||
    worker.hostname.toLowerCase().includes(q) ||
    worker.id.toLowerCase().includes(q)
  );
}

function WorkerOption({
  selected,
  onSelect,
  children,
  className,
}: {
  selected: boolean;
  onSelect: () => void;
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <button
      type="button"
      aria-pressed={selected}
      onClick={onSelect}
      className={cn(
        "flex w-full items-center justify-between gap-3 rounded-md border px-3 py-2.5 text-left text-[13px] transition-colors",
        selected
          ? "border-ring bg-accent/40"
          : "border-border hover:bg-muted/50",
        className,
      )}
    >
      {children}
      {selected ? (
        <Check className="size-4 shrink-0 text-foreground" />
      ) : (
        <span className="size-4 shrink-0" aria-hidden />
      )}
    </button>
  );
}

export function DeployDialog({
  service,
  open,
  onOpenChange,
  onDeployed,
}: DeployDialogProps) {
  const [envs, setEnvs] = useState<Environment[]>([]);
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [environment, setEnvironment] = useState("dev");
  const [workerId, setWorkerId] = useState(ANY_WORKER);
  const [workerSearch, setWorkerSearch] = useState("");
  const [loading, setLoading] = useState(false);
  const [deploying, setDeploying] = useState(false);
  const [error, setError] = useState("");

  const onlineWorkers = useMemo(
    () => workers.filter((w) => w.status === "online"),
    [workers],
  );

  const eligibleWorkers = useMemo(
    () => onlineWorkers.filter((w) => workerAllowsEnvironment(w, environment)),
    [onlineWorkers, environment],
  );

  const filteredWorkers = useMemo(
    () => eligibleWorkers.filter((w) => matchesWorkerSearch(w, workerSearch)),
    [eligibleWorkers, workerSearch],
  );

  const priorWorkerId =
    service?.env_deployments?.[environment]?.worker_id?.trim() ?? "";
  const priorWorkerOffline =
    priorWorkerId !== "" &&
    !onlineWorkers.some((w) => w.id === priorWorkerId);
  const priorWorkerIneligible =
    priorWorkerId !== "" &&
    onlineWorkers.some((w) => w.id === priorWorkerId) &&
    !eligibleWorkers.some((w) => w.id === priorWorkerId);

  useEffect(() => {
    if (!open) return;

    setError("");
    setWorkerSearch("");
    setLoading(true);
    void Promise.all([listEnvironments(), listWorkers()])
      .then(([envData, workerData]) => {
        setEnvs(envData);
        setWorkers(workerData);
        const online = workerData.filter((w) => w.status === "online");
        const dev = envData.find((e) => e.slug === "dev");
        const slug = dev?.slug ?? envData[0]?.slug ?? "dev";
        setEnvironment(slug);
        const eligible = online.filter((w) => workerAllowsEnvironment(w, slug));
        setWorkerId(defaultWorkerForEnv(service, slug, eligible));
      })
      .catch((err) => {
        setError(
          err instanceof Error ? err.message : "Failed to load deploy options",
        );
      })
      .finally(() => {
        setLoading(false);
      });
  }, [open, service]);

  useEffect(() => {
    if (!open || !service) return;
    setWorkerSearch("");
    setWorkerId(defaultWorkerForEnv(service, environment, eligibleWorkers));
  }, [environment, open, service, eligibleWorkers]);

  async function handleDeploy() {
    if (!service) return;

    setDeploying(true);
    setError("");
    try {
      await createDeployment(
        service.id,
        environment,
        workerId || undefined,
      );
      onOpenChange(false);
      onDeployed?.(service, environment);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to deploy");
    } finally {
      setDeploying(false);
    }
  }

  const pickerDisabled = deploying || loading || envs.length === 0;
  const selectedWorker = workers.find((w) => w.id === workerId);
  const showWorkerSearch =
    !loading && eligibleWorkers.length >= WORKER_SEARCH_THRESHOLD;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="flex max-h-[min(90vh,640px)] flex-col gap-0 overflow-hidden sm:max-w-md">
        <DialogHeader className="shrink-0">
          <DialogTitle>Deploy service</DialogTitle>
          <DialogDescription>
            Choose an environment and worker for{" "}
            <span className="font-medium text-foreground">
              {service?.name ?? "this service"}
            </span>
            . Each environment runs in its own container.
          </DialogDescription>
        </DialogHeader>

        <div className="min-h-0 flex-1 space-y-4 overflow-y-auto py-1">
          {error ? (
            <ErrorBanner message={error} onRetry={() => setError("")} />
          ) : null}

          <fieldset className="space-y-2" disabled={pickerDisabled}>
            <legend className="text-[12px] text-muted-foreground">
              Environment
            </legend>
            {loading ? (
              <div className="space-y-2">
                {Array.from({ length: 3 }).map((_, i) => (
                  <div
                    key={i}
                    className="h-11 animate-pulse rounded-md bg-muted"
                  />
                ))}
              </div>
            ) : (
              <div className="space-y-2">
                {envs.map((env) => {
                  const selected = environment === env.slug;
                  return (
                    <WorkerOption
                      key={env.id}
                      selected={selected}
                      onSelect={() => setEnvironment(env.slug)}
                    >
                      <span className="min-w-0">
                        <span className="block font-medium">{env.name}</span>
                        <span className="font-mono text-[11px] text-muted-foreground uppercase">
                          {env.slug}
                        </span>
                      </span>
                    </WorkerOption>
                  );
                })}
              </div>
            )}
          </fieldset>

          <fieldset className="space-y-2" disabled={pickerDisabled}>
            <div className="flex items-baseline justify-between gap-2">
              <legend className="text-[12px] text-muted-foreground">
                Worker
              </legend>
              {!loading ? (
                <span className="text-[11px] text-muted-foreground">
                  {eligibleWorkers.length} for{" "}
                  <span className="font-mono uppercase">{environment}</span>
                  {onlineWorkers.length > eligibleWorkers.length
                    ? ` · ${onlineWorkers.length} online`
                    : ""}
                </span>
              ) : null}
            </div>

            {loading ? (
              <div className="space-y-2">
                {Array.from({ length: 2 }).map((_, i) => (
                  <div
                    key={i}
                    className="h-11 animate-pulse rounded-md bg-muted"
                  />
                ))}
              </div>
            ) : (
              <div className="space-y-2">
                <WorkerOption
                  selected={workerId === ANY_WORKER}
                  onSelect={() => setWorkerId(ANY_WORKER)}
                >
                  <span className="min-w-0">
                    <span className="block font-medium">Any available</span>
                    <span className="text-[11px] text-muted-foreground">
                      Redeploy uses the previous worker for this environment
                    </span>
                  </span>
                </WorkerOption>

                {priorWorkerOffline ? (
                  <p className="rounded-md border border-amber-500/30 bg-amber-500/8 px-3 py-2 text-[11px] text-amber-800 dark:text-amber-300">
                    Previous worker for{" "}
                    <span className="font-mono uppercase">{environment}</span> is
                    offline. Pick an online worker or use Any available.
                  </p>
                ) : null}

                {priorWorkerIneligible ? (
                  <p className="rounded-md border border-amber-500/30 bg-amber-500/8 px-3 py-2 text-[11px] text-amber-800 dark:text-amber-300">
                    Previous worker for{" "}
                    <span className="font-mono uppercase">{environment}</span> does
                    not allow this environment. Pick another worker or use Any
                    available.
                  </p>
                ) : null}

                {showWorkerSearch ? (
                  <div className="relative">
                    <Search className="pointer-events-none absolute top-1/2 left-2.5 size-3.5 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      type="search"
                      value={workerSearch}
                      onChange={(e) => setWorkerSearch(e.target.value)}
                      placeholder="Search workers by name or host…"
                      aria-label="Search workers"
                      className="h-9 pl-8 text-[13px]"
                    />
                  </div>
                ) : null}

                {onlineWorkers.length === 0 ? (
                  <p className="rounded-md border border-dashed border-border px-3 py-2 text-[12px] text-muted-foreground">
                    No online workers. Start an agent or use Any available for
                    auto-affinity on redeploy.
                  </p>
                ) : eligibleWorkers.length === 0 ? (
                  <p className="rounded-md border border-dashed border-border px-3 py-2 text-[12px] text-muted-foreground">
                    No online workers allow{" "}
                    <span className="font-mono uppercase">{environment}</span>.
                    Use Any available for redeploy affinity, or pick another
                    environment.
                  </p>
                ) : (
                  <div
                    className={cn(
                      "space-y-2 overflow-y-auto pr-0.5",
                      eligibleWorkers.length > 2 && WORKER_LIST_MAX_HEIGHT,
                    )}
                  >
                    {filteredWorkers.length === 0 ? (
                      <p className="px-1 py-2 text-[12px] text-muted-foreground">
                        No workers match &ldquo;{workerSearch.trim()}&rdquo;.
                      </p>
                    ) : (
                      filteredWorkers.map((worker) => (
                        <WorkerOption
                          key={worker.id}
                          selected={workerId === worker.id}
                          onSelect={() => setWorkerId(worker.id)}
                        >
                          <span className="flex min-w-0 items-center gap-2">
                            <Server className="size-3.5 shrink-0 text-muted-foreground" />
                            <span className="min-w-0">
                              <span className="flex items-center gap-2">
                                <span className="block truncate font-medium">
                                  {worker.name}
                                </span>
                                <StatusBadge
                                  status={worker.status}
                                  className={workerStatusClass(worker.status)}
                                />
                              </span>
                              <span className="block truncate font-mono text-[11px] text-muted-foreground">
                                {worker.hostname || worker.id.slice(0, 8)}
                                {" · "}
                                {formatWorkerEnvironments(worker)}
                              </span>
                            </span>
                          </span>
                        </WorkerOption>
                      ))
                    )}
                  </div>
                )}
              </div>
            )}
          </fieldset>

          {selectedWorker && workerId !== ANY_WORKER ? (
            <p className="text-[11px] text-muted-foreground">
              Deploy will target worker{" "}
              <span className="font-medium text-foreground">
                {selectedWorker.name}
              </span>
              .
            </p>
          ) : null}
        </div>

        <DialogFooter className="shrink-0 gap-2 border-t border-border pt-4 sm:gap-0">
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 text-[13px]"
            disabled={deploying}
            onClick={() => onOpenChange(false)}
          >
            Cancel
          </Button>
          <Button
            type="button"
            size="sm"
            className="h-8 gap-1.5 text-[13px]"
            disabled={deploying || loading || envs.length === 0 || !service}
            onClick={() => void handleDeploy()}
          >
            <Rocket className="size-3.5" />
            {deploying ? "Deploying…" : "Deploy"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
