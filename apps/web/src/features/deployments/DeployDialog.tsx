import { useEffect, useState } from "react";
import { Rocket } from "lucide-react";
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
import { Label } from "@/components/ui/label";
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

const selectClassName =
  "flex h-9 w-full rounded-md border border-input bg-transparent px-3 text-[13px] shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50";

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

        <div className="space-y-1.5">
          <Label htmlFor="deploy-environment" className="text-[12px] text-muted-foreground">
            Environment
          </Label>
          {loadingEnvs ? (
            <div className="h-9 animate-pulse rounded-md bg-muted" />
          ) : (
            <select
              id="deploy-environment"
              value={environment}
              onChange={(e) => setEnvironment(e.target.value)}
              disabled={deploying || envs.length === 0}
              className={selectClassName}
            >
              {envs.map((env) => (
                <option key={env.id} value={env.slug}>
                  {env.name} ({env.slug})
                </option>
              ))}
            </select>
          )}
        </div>

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
