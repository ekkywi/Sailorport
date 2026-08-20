import { useEffect, useState } from "react";
import { Check, Rocket, Server } from "lucide-react";
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
import { cn } from "@/lib/utils";
import type { Service } from "../catalog/types";
import { listEnvironments } from "../environments/api";
import type { Environment } from "../environments/types";
import { listWorkers } from "../workers/api";
import type { Worker } from "../workers/types";
import { createDeployment } from "./api";

type DeployDialogProps = {
  service: Service | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onDeployed?: (service: Service, environment: string) => void;
};

const ANY_WORKER = "";

function defaultWorkerForEnv(service: Service | null, envSlug: string): string {
  const prior = service?.env_deployments?.[envSlug]?.worker_id;
  if (prior && prior.trim() !== "") {
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
  const [loading, setLoading] = useState(false);
  const [deploying, setDeploying] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!open) return;

    setError("");
    setLoading(true);
    void Promise.all([listEnvironments(), listWorkers()])
      .then(([envData, workerData]) => {
        setEnvs(envData);
        setWorkers(workerData);
        const dev = envData.find((e) => e.slug === "dev");
        const slug = dev?.slug ?? envData[0]?.slug ?? "dev";
        setEnvironment(slug);
        setWorkerId(defaultWorkerForEnv(service, slug));
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
    setWorkerId(defaultWorkerForEnv(service, environment));
  }, [environment, open, service]);

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

  const onlineWorkers = workers.filter((w) => w.status === "online");
  const pickerDisabled = deploying || loading || envs.length === 0;
  const selectedWorker = workers.find((w) => w.id === workerId);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Deploy service</DialogTitle>
          <DialogDescription>
            Choose an environment and worker for{" "}
            <span className="font-medium text-foreground">
              {service?.name ?? "this service"}
            </span>
            . Each environment runs in its own container.
          </DialogDescription>
        </DialogHeader>

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
                <div key={i} className="h-11 animate-pulse rounded-md bg-muted" />
              ))}
            </div>
          ) : (
            <div className="space-y-2">
              {envs.map((env) => {
                const selected = environment === env.slug;
                return (
                  <button
                    key={env.id}
                    type="button"
                    aria-pressed={selected}
                    onClick={() => setEnvironment(env.slug)}
                    className={cn(
                      "flex w-full items-center justify-between gap-3 rounded-md border px-3 py-2.5 text-left text-[13px] transition-colors",
                      selected
                        ? "border-ring bg-accent/40"
                        : "border-border hover:bg-muted/50",
                    )}
                  >
                    <span className="min-w-0">
                      <span className="block font-medium">{env.name}</span>
                      <span className="font-mono text-[11px] text-muted-foreground uppercase">
                        {env.slug}
                      </span>
                    </span>
                    {selected ? (
                      <Check className="size-4 shrink-0 text-foreground" />
                    ) : null}
                  </button>
                );
              })}
            </div>
          )}
        </fieldset>

        <fieldset className="space-y-2" disabled={pickerDisabled}>
          <legend className="text-[12px] text-muted-foreground">Worker</legend>
          {loading ? (
            <div className="space-y-2">
              {Array.from({ length: 2 }).map((_, i) => (
                <div key={i} className="h-11 animate-pulse rounded-md bg-muted" />
              ))}
            </div>
          ) : (
            <div className="space-y-2">
              <button
                type="button"
                aria-pressed={workerId === ANY_WORKER}
                onClick={() => setWorkerId(ANY_WORKER)}
                className={cn(
                  "flex w-full items-center justify-between gap-3 rounded-md border px-3 py-2.5 text-left text-[13px] transition-colors",
                  workerId === ANY_WORKER
                    ? "border-ring bg-accent/40"
                    : "border-border hover:bg-muted/50",
                )}
              >
                <span className="min-w-0">
                  <span className="block font-medium">Any available</span>
                  <span className="text-[11px] text-muted-foreground">
                    Redeploy uses the previous worker for this environment
                  </span>
                </span>
                {workerId === ANY_WORKER ? (
                  <Check className="size-4 shrink-0 text-foreground" />
                ) : null}
              </button>
              {onlineWorkers.length === 0 ? (
                <p className="rounded-md border border-dashed border-border px-3 py-2 text-[12px] text-muted-foreground">
                  No online workers. Start an agent or pick Any available for
                  auto-affinity on redeploy.
                </p>
              ) : (
                onlineWorkers.map((worker) => {
                  const selected = workerId === worker.id;
                  return (
                    <button
                      key={worker.id}
                      type="button"
                      aria-pressed={selected}
                      onClick={() => setWorkerId(worker.id)}
                      className={cn(
                        "flex w-full items-center justify-between gap-3 rounded-md border px-3 py-2.5 text-left text-[13px] transition-colors",
                        selected
                          ? "border-ring bg-accent/40"
                          : "border-border hover:bg-muted/50",
                      )}
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
                          </span>
                        </span>
                      </span>
                      {selected ? (
                        <Check className="size-4 shrink-0 text-foreground" />
                      ) : null}
                    </button>
                  );
                })
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

        <DialogFooter className="gap-2 sm:gap-0">
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
