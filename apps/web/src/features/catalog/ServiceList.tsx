import { Boxes, History, Pencil, Play, Rocket, ScrollText, Square, Trash2 } from "lucide-react";
import {
  DataPanel,
  EmptyState,
  EnvironmentBadge,
  StatusBadge,
  truncateMiddle,
  userInitials,
} from "@/components/app";
import { Button } from "@/components/ui/button";
import type { Environment } from "../environments/types";
import { canDeployService, type Service } from "./types";

type ServiceListProps = {
  services: Service[];
  environments: Environment[];
  loading: boolean;
  canWrite: boolean;
  onEdit: (svc: Service) => void;
  onDelete: (svc: Service) => void;
  onCreate?: () => void;
  onCreateFromGit?: () => void;
  onDeploy: (svc: Service) => void;
  onOpenHistory: (svc: Service) => void;
  onStop: (svc: Service, environment: string) => void;
  onStart: (svc: Service, environment: string) => void;
  onLogs: (svc: Service, environment: string) => void;
};

const FALLBACK_ENVIRONMENTS: Pick<Environment, "slug" | "name" | "sort_order">[] = [
  { slug: "dev", name: "Development", sort_order: 0 },
  { slug: "staging", name: "Staging", sort_order: 1 },
  { slug: "prod", name: "Production", sort_order: 2 },
];

function deployBadgeClass(status: string) {
  if (status === "running") {
    return "bg-emerald-500/12 text-emerald-700 dark:text-emerald-400";
  }
  if (status === "failed") {
    return "bg-red-500/12 text-red-700 dark:text-red-400";
  }
  if (status === "building" || status === "claimed" || status === "pending") {
    return "bg-amber-500/12 text-amber-700 dark:text-amber-400";
  }
  if (status === "stopped") {
    return "bg-muted text-muted-foreground";
  }
  return undefined;
}

function shortRepoName(repoURL: string): string {
  const cleaned = repoURL.replace(/\.git$/i, "").replace(/\/$/, "");
  const parts = cleaned.split("/").filter(Boolean);
  if (parts.length >= 2) {
    return `${parts[parts.length - 2]}/${parts[parts.length - 1]}`;
  }
  return truncateMiddle(repoURL, 32);
}

function OriginCell({ svc }: { svc: Service }) {
  if (svc.source_type === "git") {
    return (
      <div className="flex min-w-0 flex-col gap-0.5">
        <span className="inline-flex max-w-full flex-wrap items-center gap-1.5">
          <span className="rounded-md bg-sky-500/12 px-1.5 py-0.5 text-[10px] font-medium tracking-wide text-sky-700 uppercase dark:text-sky-400">
            Git
          </span>
          <span className="truncate text-[11px] text-muted-foreground">
            {svc.branch || "main"}
          </span>
          {svc.auto_deploy_enabled ? (
            <span
              className="rounded-md bg-amber-500/12 px-1.5 py-0.5 text-[10px] font-medium tracking-wide text-amber-800 uppercase dark:text-amber-400"
              title={`Auto-deploy → ${svc.auto_deploy_environment || "staging"}`}
            >
              Auto
            </span>
          ) : null}
        </span>
        {svc.repo_url ? (
          <span
            className="block truncate font-mono text-[11px] text-muted-foreground"
            title={svc.repo_url}
          >
            {shortRepoName(svc.repo_url)}
          </span>
        ) : null}
      </div>
    );
  }

  if (svc.template_id) {
    return (
      <div className="flex min-w-0 flex-col gap-0.5">
        <span className="inline-flex max-w-full items-center gap-1.5">
          <span className="rounded-md bg-emerald-500/12 px-1.5 py-0.5 text-[10px] font-medium tracking-wide text-emerald-700 uppercase dark:text-emerald-400">
            Template
          </span>
          <span className="truncate text-[11px] text-muted-foreground">
            {svc.template_id}
          </span>
        </span>
        {svc.workspace_path ? (
          <span
            className="block truncate font-mono text-[11px] text-muted-foreground"
            title={svc.workspace_path}
          >
            {truncateMiddle(svc.workspace_path, 28)}
          </span>
        ) : null}
      </div>
    );
  }

  return (
    <span className="rounded-md bg-muted px-1.5 py-0.5 text-[10px] font-medium tracking-wide text-muted-foreground uppercase">
      Registered
    </span>
  );
}

function ServiceGlyph({ name }: { name: string }) {
  return (
    <span className="flex size-8 shrink-0 items-center justify-center rounded-md bg-muted text-[11px] font-semibold tracking-tight text-muted-foreground">
      {userInitials(name)}
    </span>
  );
}

function EnvDeployCell({
  svc,
  environments,
  canWrite,
  onOpenHistory,
  onStop,
  onStart,
  onLogs,
}: {
  svc: Service;
  environments: Environment[];
  canWrite: boolean;
  onOpenHistory: (svc: Service) => void;
  onStop: (svc: Service, environment: string) => void;
  onStart: (svc: Service, environment: string) => void;
  onLogs: (svc: Service, environment: string) => void;
}) {
  const envs = svc.env_deployments ?? {};
  const ordered =
    environments.length > 0
      ? [...environments].sort((a, b) => a.sort_order - b.sort_order)
      : FALLBACK_ENVIRONMENTS;

  return (
    <div className="flex min-w-0 flex-col gap-1.5">
      {ordered.map((env) => {
        const slug = env.slug;
        const d = envs[slug];
        return (
          <div key={slug} className="flex min-w-0 items-center gap-1.5">
            <EnvironmentBadge slug={slug} />
            {!d ? (
              <span className="text-[11px] text-muted-foreground">Not deployed</span>
            ) : (
              <>
                <button
                  type="button"
                  className="rounded-md outline-none hover:opacity-90 focus-visible:ring-2 focus-visible:ring-ring/40"
                  onClick={() => onOpenHistory(svc)}
                  title={`History ${svc.name} (${slug})`}
                >
                  <StatusBadge
                    status={d.status}
                    className={deployBadgeClass(d.status)}
                  />
                </button>
                {d.status === "running" && d.port != null ? (
                  <a
                    href={`http://localhost:${d.port}/healthz`}
                    target="_blank"
                    rel="noreferrer"
                    className="text-[11px] text-muted-foreground underline-offset-2 hover:text-foreground hover:underline"
                    onClick={(e) => e.stopPropagation()}
                  >
                    :{d.port}
                  </a>
                ) : null}
                {(d.status === "running" || d.status === "stopped") ? (
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon-sm"
                    className="size-7 text-muted-foreground hover:text-foreground"
                    aria-label={`Logs ${svc.name} ${slug}`}
                    onClick={() => onLogs(svc, slug)}
                  >
                    <ScrollText className="size-3" />
                  </Button>
                ) : null}
                {canWrite && d.status === "running" ? (
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon-sm"
                    className="size-7 text-muted-foreground hover:text-foreground"
                    aria-label={`Stop ${svc.name} ${slug}`}
                    onClick={() => onStop(svc, slug)}
                  >
                    <Square className="size-3" />
                  </Button>
                ) : null}
                {canWrite && d.status === "stopped" ? (
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon-sm"
                    className="size-7 text-muted-foreground hover:text-foreground"
                    aria-label={`Start ${svc.name} ${slug}`}
                    onClick={() => onStart(svc, slug)}
                  >
                    <Play className="size-3" />
                  </Button>
                ) : null}
              </>
            )}
          </div>
        );
      })}
    </div>
  );
}

function RowActions({
  svc,
  canWrite,
  onDeploy,
  onEdit,
  onDelete,
  onOpenHistory,
}: {
  svc: Service;
  canWrite: boolean;
  onDeploy: (svc: Service) => void;
  onEdit: (svc: Service) => void;
  onDelete: (svc: Service) => void;
  onOpenHistory: (svc: Service) => void;
}) {
  return (
    <div className="flex shrink-0 items-center justify-end gap-0.5 whitespace-nowrap">
      <Button
        type="button"
        variant="ghost"
        size="icon-sm"
        className="text-muted-foreground hover:text-foreground"
        aria-label={`History ${svc.name}`}
        onClick={() => onOpenHistory(svc)}
      >
        <History className="size-3.5" />
      </Button>
      {canWrite ? (
        <>
          {canDeployService(svc) ? (
            <Button
              type="button"
              variant="ghost"
              size="icon-sm"
              className="text-muted-foreground hover:text-foreground"
              aria-label={`Deploy ${svc.name}`}
              onClick={() => onDeploy(svc)}
            >
              <Rocket className="size-3.5" />
            </Button>
          ) : null}
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            className="text-muted-foreground hover:text-foreground"
            aria-label={`Edit ${svc.name}`}
            onClick={() => onEdit(svc)}
          >
            <Pencil className="size-3.5" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            className="text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
            aria-label={`Delete ${svc.name}`}
            onClick={() => onDelete(svc)}
          >
            <Trash2 className="size-3.5" />
          </Button>
        </>
      ) : null}
    </div>
  );
}

export function ServiceList({
  services,
  environments,
  loading,
  canWrite,
  onEdit,
  onDelete,
  onDeploy,
  onOpenHistory,
  onStop,
  onStart,
  onLogs,
  onCreate,
  onCreateFromGit,
}: ServiceListProps) {
  if (loading && services.length === 0) {
    return (
      <DataPanel>
        <div className="space-y-0 divide-y divide-border">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="flex items-center gap-3 px-4 py-3.5">
              <div className="size-8 animate-pulse rounded-md bg-muted" />
              <div className="flex-1 space-y-2">
                <div className="h-3.5 w-1/3 animate-pulse rounded bg-muted" />
                <div className="h-3 w-1/2 animate-pulse rounded bg-muted" />
              </div>
            </div>
          ))}
        </div>
      </DataPanel>
    );
  }

  if (!loading && services.length === 0) {
    return (
      <DataPanel>
        <EmptyState
          icon={Boxes}
          title="No services yet"
          description={
            canWrite
              ? "Add a service from Git to deploy your repo, or start from a template."
              : "No services yet. Ask an admin or developer to add one."
          }
          action={
            canWrite && (onCreateFromGit || onCreate) ? (
              <div className="flex flex-wrap justify-center gap-2">
                {onCreateFromGit ? (
                  <Button
                    type="button"
                    size="sm"
                    className="h-8 text-[13px]"
                    onClick={onCreateFromGit}
                  >
                    Add from Git
                  </Button>
                ) : null}
                {onCreate ? (
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="h-8 text-[13px]"
                    onClick={onCreate}
                  >
                    More options
                  </Button>
                ) : null}
              </div>
            ) : undefined
          }
          className="py-16"
        />
      </DataPanel>
    );
  }

  return (
    <DataPanel>
      <div className="hidden overflow-x-auto md:block">
        <table className="w-full text-left text-[13px]">
          <thead>
            <tr className="border-b border-border text-[11px] font-medium tracking-[0.04em] text-muted-foreground uppercase">
              <th className="min-w-[220px] px-4 py-2.5 font-medium">Service</th>
              <th className="w-[12%] px-4 py-2.5 font-medium">Owner</th>
              <th className="min-w-[200px] px-4 py-2.5 font-medium">Deploy</th>
              <th className="min-w-[240px] px-4 py-2.5 font-medium">Origin</th>
              <th className="w-[120px] px-4 py-2.5 font-medium">
                <span className="sr-only">Actions</span>
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {services.map((svc) => (
              <tr
                key={svc.id}
                className="group transition-colors hover:bg-muted/35"
              >
                <td className="px-4 py-3">
                  <div className="flex min-w-0 items-start gap-3">
                    <ServiceGlyph name={svc.name} />
                    <div className="min-w-0 pt-0.5">
                      <p className="truncate font-medium tracking-[-0.01em]">
                        {svc.name}
                      </p>
                      <p className="mt-0.5 truncate text-[12px] text-muted-foreground">
                        {svc.description || "No description"}
                      </p>
                    </div>
                  </div>
                </td>
                <td className="px-4 py-3">
                  {svc.owner ? (
                    <span className="inline-block max-w-full truncate rounded-md bg-muted px-2 py-0.5 text-[12px] text-foreground/80">
                      {svc.owner}
                    </span>
                  ) : (
                    <span className="text-muted-foreground">—</span>
                  )}
                </td>
                <td className="px-4 py-3">
                  <EnvDeployCell
                    svc={svc}
                    environments={environments}
                    canWrite={canWrite}
                    onOpenHistory={onOpenHistory}
                    onStop={onStop}
                    onStart={onStart}
                    onLogs={onLogs}
                  />
                </td>
                <td className="px-4 py-3">
                  <OriginCell svc={svc} />
                </td>
                <td className="px-4 py-3">
                  <RowActions
                    svc={svc}
                    canWrite={canWrite}
                    onDeploy={onDeploy}
                    onEdit={onEdit}
                    onDelete={onDelete}
                    onOpenHistory={onOpenHistory}
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <ul className="divide-y divide-border md:hidden">
        {services.map((svc) => (
          <li key={svc.id} className="px-4 py-3.5">
            <div className="flex items-start gap-3">
              <ServiceGlyph name={svc.name} />
              <div className="min-w-0 flex-1">
                <div className="flex items-start justify-between gap-2">
                  <div className="min-w-0">
                    <p className="truncate text-[13px] font-medium tracking-[-0.01em]">
                      {svc.name}
                    </p>
                    <p className="mt-0.5 truncate text-[12px] text-muted-foreground">
                      {svc.owner || "Unassigned"}
                      {svc.description ? ` · ${svc.description}` : ""}
                    </p>
                  </div>
                  <RowActions
                    svc={svc}
                    canWrite={canWrite}
                    onDeploy={onDeploy}
                    onEdit={onEdit}
                    onDelete={onDelete}
                    onOpenHistory={onOpenHistory}
                  />
                </div>
                <div className="mt-2">
                  <EnvDeployCell
                    svc={svc}
                    environments={environments}
                    canWrite={canWrite}
                    onOpenHistory={onOpenHistory}
                    onStop={onStop}
                    onStart={onStart}
                    onLogs={onLogs}
                  />
                </div>
                <div className="mt-2">
                  <OriginCell svc={svc} />
                </div>
              </div>
            </div>
          </li>
        ))}
      </ul>
    </DataPanel>
  );
}
