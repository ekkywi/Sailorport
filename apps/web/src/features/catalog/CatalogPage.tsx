import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { Plus, RefreshCw } from "lucide-react";
import { ErrorBanner, Toolbar, useToast } from "@/components/app";
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
import { cn } from "@/lib/utils";
import { CreateServiceForm } from "../scaffold/CreateServiceForm";
import {
  createService,
  deleteService,
  listServices,
  updateService,
} from "./api";
import { ServiceForm } from "./ServiceForm";
import { ServiceList } from "./ServiceList";
import type { Service, ServiceFormValues } from "./types";
import { DeployDialog } from "../deployments/DeployDialog";
import { DeploymentsDialog } from "../deployments/DeploymentsDialog";
import { LogsDialog } from "./LogsDialog";
import { startService, stopService } from "../runtime/api";
import type { AuthUser } from "@/features/auth/types";
import { canWriteCatalog } from "@/lib/rbac";
import { listEnvironments } from "../environments/api";
import type { Environment } from "../environments/types";

const emptyForm: ServiceFormValues = {
  name: "",
  description: "",
  owner: "",
};

type DialogMode = "none" | "create" | "register" | "edit";

export function CatalogPage({currentUser}: {currentUser: AuthUser}) {
  const canWrite = canWriteCatalog(currentUser.role);
  const { toast } = useToast();
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [listError, setListError] = useState("");
  const [formError, setFormError] = useState("");
  const [values, setValues] = useState<ServiceFormValues>(emptyForm);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [dialog, setDialog] = useState<DialogMode>("none");
  const [deleteTarget, setDeleteTarget] = useState<Service | null>(null);
  const [createdPath, setCreatedPath] = useState("");
  const [historyTarget, setHistoryTarget] = useState<Service | null>(null);
  const [deployDialogTarget, setDeployDialogTarget] = useState<Service | null>(null);
  const [runtimeTarget, setRuntimeTarget] = useState<{
    service: Service;
    environment: string;
    action: "stop" | "start";
  } | null>(null);
  const [runtimePending, setRuntimePending] = useState(false);
  const [logsTarget, setLogsTarget] = useState<{
    service: Service;
    environment: string;
  } | null>(null);
  const [environments, setEnvironments] = useState<Environment[]>([]);

  async function load(options?: { silent?: boolean }) {
    if (!options?.silent) {
      setLoading(true);
    }
    setListError("");
    try {
      setServices(await listServices());
    } catch (err) {
      setListError(err instanceof Error ? err.message : "Failed to load catalog");
    } finally {
      if (!options?.silent) {
        setLoading(false);
      }
    }
  }

  function openHistory(svc: Service) {
    setHistoryTarget(svc);
  }

  function handleDeployed(svc: Service, environment: string) {
    toast(`Deploy to ${environment} started for "${svc.name}"`);
    setHistoryTarget(svc);
    void load();
  }

  function refreshAfterRuntime() {
    void load({ silent: true });
    window.setTimeout(() => void load({ silent: true }), 3000);
    window.setTimeout(() => void load({ silent: true }), 6000);
  }

  async function confirmRuntime() {
    if (!runtimeTarget) return;
    const { service: svc, environment, action } = runtimeTarget;
    setRuntimePending(true);
    setListError("");
    try {
      if (action === "stop") {
        await stopService(svc.id, environment);
        toast(`Stopping ${environment} for "${svc.name}"…`);
      } else {
        await startService(svc.id, environment);
        toast(`Starting ${environment} for "${svc.name}"…`);
      }
      setRuntimeTarget(null);
      refreshAfterRuntime();
    } catch (err) {
      setListError(
        err instanceof Error
          ? err.message
          : action === "stop"
            ? "Failed to stop service"
            : "Failed to start service",
      );
      setRuntimeTarget(null);
    } finally {
      setRuntimePending(false);
    }
  }

  useEffect(() => {
    void load();
    void listEnvironments()
      .then(setEnvironments)
      .catch(() => {
        /* fallback handled in ServiceList */
      });
  }, []);

  useEffect(() => {
    const hasActive = services.some((svc) => {
      const envActive = Object.values(svc.env_deployments ?? {}).some(
        (d) =>
          d.status === "pending" ||
          d.status === "claimed" ||
          d.status === "building",
      );
      const latest = svc.latest_deployment?.status;
      const latestActive =
        latest === "pending" ||
        latest === "claimed" ||
        latest === "building";
      return envActive || latestActive;
    });
    if (!hasActive) return;

    const id = window.setInterval(() => {
      void load({ silent: true });
    }, 5000);
    return () => window.clearInterval(id);
  }, [services]);

  function closeDialog() {
    setValues(emptyForm);
    setEditingId(null);
    setFormError("");
    setCreatedPath("");
    setDialog("none");
  }

  function startCreate() {
    setEditingId(null);
    setValues(emptyForm);
    setFormError("");
    setCreatedPath("");
    setDialog("create");
  }

  function startRegister() {
    setEditingId(null);
    setValues(emptyForm);
    setFormError("");
    setCreatedPath("");
    setDialog("register");
  }

  function startEdit(svc: Service) {
    setEditingId(svc.id);
    setValues({
      name: svc.name,
      description: svc.description,
      owner: svc.owner,
    });
    setFormError("");
    setDialog("edit");
  }

  function onChange(field: keyof ServiceFormValues, value: string) {
    setValues((prev) => ({ ...prev, [field]: value }));
  }

  async function onSubmitMetadata(e: FormEvent) {
    e.preventDefault();
    setSaving(true);
    setFormError("");
    try {
      if (editingId) {
        await updateService(editingId, values);
        toast("Service updated");
      } else {
        await createService(values);
        toast("Service registered");
      }
      closeDialog();
      await load();
    } catch (err) {
      setFormError(
        err instanceof Error
          ? err.message
          : editingId
            ? "Failed to update service"
            : "Failed to register service",
      );
    } finally {
      setSaving(false);
    }
  }

  const dialogOpen = dialog !== "none";
  const countLabel =
    loading && services.length === 0
      ? "Loading…"
      : `${services.length} service${services.length === 1 ? "" : "s"}`;

  const deleteRunningEnvs =
    deleteTarget == null
      ? []
      : Object.entries(deleteTarget.env_deployments ?? {})
          .filter(([, d]) => d.status === "running")
          .map(([slug]) => slug);

  const deleteBlocked = deleteRunningEnvs.length > 0;
  const deleteHasProdRunning = deleteRunningEnvs.includes("prod");

  const deleteContainers =
    deleteTarget == null
      ? []
      : Object.entries(deleteTarget.env_deployments ?? {}).map(([slug]) => ({
          slug,
          name: `sailorport-${deleteTarget.name}-${slug}`,
        }));

  async function confirmDelete() {
    if (!deleteTarget || deleteBlocked) return;
    const name = deleteTarget.name;
    setDeleting(true);
    setListError("");
    try {
      await deleteService(deleteTarget.id);
      if (editingId === deleteTarget.id) {
        closeDialog();
      }
      setDeleteTarget(null);
      toast(`Deleted “${name}”`);
      await load();
    } catch (err) {
      setListError(err instanceof Error ? err.message : "Failed to delete service");
      setDeleteTarget(null);
    } finally {
      setDeleting(false);
    }
  }

  return (
    <div className="space-y-4">
      <Toolbar
        meta={countLabel}
        actions={
          <>
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
            {canWrite ? (
              <Button
                type="button"
                size="sm"
                className="h-8 gap-1.5 text-[13px]"
                onClick={startCreate}
              >
                <Plus className="size-3.5" />
                Create service
              </Button>
            ) : null}
          </>
        }
      />

      {listError ? (
        <ErrorBanner message={listError} onRetry={() => void load()} />
      ) : null}

      <ServiceList
        services={services}
        environments={environments}
        loading={loading}
        canWrite={canWrite}
        onEdit={startEdit}
        onDelete={setDeleteTarget}
        onDeploy={setDeployDialogTarget}
        onOpenHistory={openHistory}
        onStop={(svc, environment) =>
          setRuntimeTarget({ service: svc, environment, action: "stop" })
        }
        onStart={(svc, environment) =>
          setRuntimeTarget({ service: svc, environment, action: "start" })
        }
        onLogs={(svc, environment) => setLogsTarget({ service: svc, environment })}
        onCreate={canWrite ? startCreate : undefined}
      />

      <Dialog
        open={dialogOpen}
        onOpenChange={(open) => {
          if (!open) closeDialog();
        }}
      >
        <DialogContent className="sm:max-w-lg" showCloseButton={!saving}>
          {dialog === "create" ? (
            <>
              <DialogHeader>
                <DialogTitle>
                  {createdPath ? "Service created" : "Create service"}
                </DialogTitle>
                <DialogDescription>
                  {createdPath
                    ? "Workspace generated and registered in the catalog."
                    : "Generate a workspace from a template and register it in the catalog."}
                </DialogDescription>
              </DialogHeader>
              {createdPath ? (
                <div className="space-y-4">
                  <div className="rounded-lg border border-border bg-muted/40 px-3 py-2.5">
                    <p className="text-[12px] text-muted-foreground">Workspace</p>
                    <p
                      className="mt-1 break-all font-mono text-[12px] text-foreground"
                      title={createdPath}
                    >
                      {createdPath}
                    </p>
                  </div>
                  <Button
                    type="button"
                    size="sm"
                    className="h-8 text-[13px]"
                    onClick={closeDialog}
                  >
                    Done
                  </Button>
                </div>
              ) : (
                <CreateServiceForm
                  onSuccess={(path) => {
                    setCreatedPath(path);
                    toast("Service created");
                    void load();
                  }}
                  onRegisterExisting={startRegister}
                />
              )}
            </>
          ) : null}

          {dialog === "register" ? (
            <>
              <DialogHeader>
                <DialogTitle>Register existing service</DialogTitle>
                <DialogDescription>
                  Catalog only — no workspace is generated. Use this for services
                  that already exist elsewhere.
                </DialogDescription>
              </DialogHeader>
              <ServiceForm
                mode="register"
                values={values}
                saving={saving}
                error={formError}
                onChange={onChange}
                onSubmit={onSubmitMetadata}
                onCancel={closeDialog}
                onCreateInstead={startCreate}
              />
            </>
          ) : null}

          {dialog === "edit" ? (
            <>
              <DialogHeader>
                <DialogTitle>Edit service</DialogTitle>
                <DialogDescription>
                  Update catalog metadata for this service.
                </DialogDescription>
              </DialogHeader>
              <ServiceForm
                mode="edit"
                values={values}
                saving={saving}
                error={formError}
                onChange={onChange}
                onSubmit={onSubmitMetadata}
                onCancel={closeDialog}
              />
            </>
          ) : null}
        </DialogContent>
      </Dialog>

      <AlertDialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open && !deleting) setDeleteTarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete service?</AlertDialogTitle>
            <AlertDialogDescription>
              This removes{" "}
              <span className="font-medium text-foreground">
                {deleteTarget?.name}
              </span>{" "}
              from the catalog and deletes its workspace folder (when under the configured
              workspace directory).
            </AlertDialogDescription>
            {deleteContainers.length > 0 ? (
              <>
                <p className="text-[13px] text-muted-foreground">
                  Cleanup jobs will remove these Docker containers on the agent node:
                </p>
                <ul className="list-inside list-disc space-y-0.5 text-[13px] text-muted-foreground">
                  {deleteContainers.map(({ slug, name }) => (
                    <li key={slug}>
                      <span className="font-mono text-[12px]">{name}</span>
                      {deleteRunningEnvs.includes(slug) ? (
                        <span className="text-destructive"> (still running)</span>
                      ) : null}
                    </li>
                  ))}
                </ul>
              </>
            ) : null}
            {deleteBlocked ? (
              <p className="text-[13px] text-destructive">
                {deleteHasProdRunning
                  ? "Production is still running. Stop prod before deleting this service."
                  : `Stop running environments first: ${deleteRunningEnvs.join(", ")}.`}
              </p>
            ) : null}
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogClose
              render={
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="h-8 text-[13px]"
                  disabled={deleting}
                />
              }
            >
              Cancel
            </AlertDialogClose>
            <Button
              type="button"
              variant="destructive"
              size="sm"
              className="h-8 text-[13px]"
              disabled={deleting || deleteBlocked}
              onClick={() => void confirmDelete()}
            >
              {deleting ? "Deleting…" : "Delete"}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
      <AlertDialog
        open={runtimeTarget !== null}
        onOpenChange={(open) => {
          if (!open && !runtimePending) setRuntimeTarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {runtimeTarget?.action === "stop" ? "Stop service?" : "Start service?"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {runtimeTarget?.action === "stop" ? (
                <>
                  This stops the running container for{" "}
                  <span className="font-medium text-foreground">
                    {runtimeTarget.service.name}
                  </span>
                  {" "}
                  in{" "}
                  <span className="font-mono text-[12px]">
                    {runtimeTarget.environment}
                  </span>
                  . The deployment stays registered; use Start to bring it back.
                </>
              ) : (
                <>
                  This starts the container for{" "}
                  <span className="font-medium text-foreground">
                    {runtimeTarget?.service.name}
                  </span>
                  {" "}
                  in{" "}
                  <span className="font-mono text-[12px]">
                    {runtimeTarget?.environment}
                  </span>
                  .
                </>
              )}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogClose
              render={
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="h-8 text-[13px]"
                  disabled={runtimePending}
                />
              }
            >
              Cancel
            </AlertDialogClose>
            <Button
              type="button"
              size="sm"
              className="h-8 text-[13px]"
              disabled={runtimePending}
              onClick={() => void confirmRuntime()}
            >
              {runtimePending
                ? runtimeTarget?.action === "stop"
                  ? "Stopping…"
                  : "Starting…"
                : runtimeTarget?.action === "stop"
                  ? "Stop"
                  : "Start"}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
      <DeployDialog
        service={deployDialogTarget}
        open={deployDialogTarget !== null}
        onOpenChange={(open) => {
          if (!open) setDeployDialogTarget(null);
        }}
        onDeployed={handleDeployed}
      />
      <DeploymentsDialog
        open={historyTarget !== null}
        serviceId={historyTarget?.id ?? null}
        serviceName={historyTarget?.name ?? ""}
        onOpenChange={(open) => {
          if (!open) setHistoryTarget(null);
        }}
        onRefreshCatalog={() => void load({ silent: true })}
      />
      <LogsDialog
        open={logsTarget !== null}
        serviceId={logsTarget?.service.id ?? null}
        serviceName={logsTarget?.service.name ?? ""}
        environment={logsTarget?.environment ?? "dev"}
        onOpenChange={(open) => {
          if (!open) setLogsTarget(null);
        }}
      />
    </div>
  );
}