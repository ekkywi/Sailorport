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
import { createDeployment } from "../deployments/api";
import { DeploymentsDialog } from "../deployments/DeploymentsDialog";

const emptyForm: ServiceFormValues = {
  name: "",
  description: "",
  owner: "",
};

type DialogMode = "none" | "create" | "register" | "edit";

export function CatalogPage() {
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
  const [deployTarget, setDeployTarget] = useState<Service | null>(null);
  const [, setDeploying] = useState(false);

  async function load() {
    setLoading(true);
    setListError("");
    try {
      setServices(await listServices());
    } catch (err) {
      setListError(err instanceof Error ? err.message : "Failed to load catalog");
    } finally {
      setLoading(false);
    }
  }

  async function startDeploy(svc: Service) {
    setDeploying(true);
    setListError("");
    try {
      await createDeployment(svc.id);
      toast(`Deploy started for "${svc.name}"`);
      setDeployTarget(svc);
    } catch (err) {
      setListError(err instanceof Error ? err.message : "Failed to deploy");
    } finally {
      setDeploying(false);
    }
  }

  useEffect(() => {
    void load();
  }, []);

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

  async function confirmDelete() {
    if (!deleteTarget) return;
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

  const dialogOpen = dialog !== "none";
  const countLabel =
    loading && services.length === 0
      ? "Loading…"
      : `${services.length} service${services.length === 1 ? "" : "s"}`;

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
            <Button
              type="button"
              size="sm"
              className="h-8 gap-1.5 text-[13px]"
              onClick={startCreate}
            >
              <Plus className="size-3.5" />
              Create service
            </Button>
          </>
        }
      />

      {listError ? (
        <ErrorBanner message={listError} onRetry={() => void load()} />
      ) : null}

      <ServiceList
        services={services}
        loading={loading}
        onEdit={startEdit}
        onDelete={setDeleteTarget}
        onDeploy={(svc) => void startDeploy(svc)}
        onCreate={startCreate}
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
              from the catalog and deletes its workspace folder on disk (if it
              lives under the configured workspace directory).
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
              disabled={deleting}
              onClick={() => void confirmDelete()}
            >
              {deleting ? "Deleting…" : "Delete"}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
      <DeploymentsDialog
        open={deployTarget !== null}
        serviceId={deployTarget?.id ?? null}
        serviceName={deployTarget?.name ?? ""}
        onOpenChange={(open) => {
          if (!open) setDeployTarget(null);
        }}
      />
    </div>
  );
}