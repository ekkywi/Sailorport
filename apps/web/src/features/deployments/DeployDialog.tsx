import { useEffect, useState } from "react";
import { Check, Rocket } from "lucide-react";
import { ErrorBanner } from "@/components/app";
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
import { createDeployment } from "./api";

type DeployDialogProps = {
  service: Service | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onDeployed?: (service: Service, environment: string) => void;
};

export function DeployDialog({
  service,
  open,
  onOpenChange,
  onDeployed,
}: DeployDialogProps) {
  const [envs, setEnvs] = useState<Environment[]>([]);
  const [environment, setEnvironment] = useState("dev");
  const [loadingEnvs, setLoadingEnvs] = useState(false);
  const [deploying, setDeploying] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!open) return;

    setError("");
    setLoadingEnvs(true);
    void listEnvironments()
      .then((data) => {
        setEnvs(data);
        const dev = data.find((e) => e.slug === "dev");
        setEnvironment(dev?.slug ?? data[0]?.slug ?? "dev");
      })
      .catch((err) => {
        setError(
          err instanceof Error ? err.message : "Failed to load environments",
        );
      })
      .finally(() => {
        setLoadingEnvs(false);
      });
  }, [open]);

  async function handleDeploy() {
    if (!service) return;

    setDeploying(true);
    setError("");
    try {
      await createDeployment(service.id, environment);
      onOpenChange(false);
      onDeployed?.(service, environment);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to deploy");
    } finally {
      setDeploying(false);
    }
  }

  const pickerDisabled = deploying || loadingEnvs || envs.length === 0;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Deploy service</DialogTitle>
          <DialogDescription>
            Choose an environment for{" "}
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
          {loadingEnvs ? (
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
            disabled={deploying || loadingEnvs || envs.length === 0 || !service}
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
