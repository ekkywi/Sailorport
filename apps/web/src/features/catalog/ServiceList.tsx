import { Boxes, FolderGit2, Pencil, Rocket, Trash2 } from "lucide-react";
import {
  DataPanel,
  EmptyState,
  truncateMiddle,
  userInitials,
} from "@/components/app";
import { Button } from "@/components/ui/button";
import type { Service } from "./types";

type ServiceListProps = {
  services: Service[];
  loading: boolean;
  onEdit: (svc: Service) => void;
  onDelete: (svc: Service) => void;
  onCreate?: () => void;
  onDeploy: (svc: Service) => void;
};

function ServiceGlyph({ name }: { name: string }) {
  return (
    <span className="flex size-8 shrink-0 items-center justify-center rounded-md bg-muted text-[11px] font-semibold tracking-tight text-muted-foreground">
      {userInitials(name)}
    </span>
  );
}

export function ServiceList({
  services,
  loading,
  onEdit,
  onDelete,
  onDeploy,
  onCreate,
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
          title="No services in catalog"
          description="Create a service from a template to generate a workspace and register it here."
          action={
            onCreate ? (
              <Button
                type="button"
                size="sm"
                className="h-8 text-[13px]"
                onClick={onCreate}
              >
                Create service
              </Button>
            ) : undefined
          }
          className="py-16"
        />
      </DataPanel>
    );
  }

  return (
    <DataPanel>
      {/* Desktop table */}
      <div className="hidden overflow-x-auto md:block">
        <table className="w-full table-fixed text-left text-[13px]">
          <thead>
            <tr className="border-b border-border text-[11px] font-medium tracking-[0.04em] text-muted-foreground uppercase">
              <th className="w-[36%] px-4 py-2.5 font-medium">Service</th>
              <th className="w-[18%] px-4 py-2.5 font-medium">Owner</th>
              <th className="w-[28%] px-4 py-2.5 font-medium">Origin</th>
              <th className="w-[18%] px-4 py-2.5 font-medium">
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
                  <div className="flex min-w-0 flex-col gap-0.5">
                    {svc.template_id ? (
                      <span className="inline-flex max-w-full items-center gap-1 truncate rounded-md border border-border px-1.5 py-0.5 text-[11px] text-muted-foreground">
                        <FolderGit2 className="size-3 shrink-0" />
                        <span className="truncate">{svc.template_id}</span>
                      </span>
                    ) : (
                      <span className="text-[12px] text-muted-foreground">
                        Registered only
                      </span>
                    )}
                    {svc.workspace_path ? (
                      <span
                        className="truncate font-mono text-[11px] text-muted-foreground"
                        title={svc.workspace_path}
                      >
                        {truncateMiddle(svc.workspace_path, 36)}
                      </span>
                    ) : null}
                  </div>
                </td>
                <td className="px-4 py-3">
                  <div className="flex shrink-0 items-center justify-end gap-0.5 whitespace-nowrap">
                    {svc.workspace_path ? (
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
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Mobile cards */}
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
                  <div className="flex shrink-0 gap-0.5">
                    {svc.workspace_path ? (
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon-sm"
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
                      aria-label={`Edit ${svc.name}`}
                      onClick={() => onEdit(svc)}
                    >
                      <Pencil className="size-3.5" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon-sm"
                      className="text-destructive"
                      aria-label={`Delete ${svc.name}`}
                      onClick={() => onDelete(svc)}
                    >
                      <Trash2 className="size-3.5" />
                    </Button>
                  </div>
                </div>
                {(svc.template_id || svc.workspace_path) && (
                  <p className="mt-2 truncate font-mono text-[11px] text-muted-foreground">
                    {svc.template_id || "registered"}
                    {svc.workspace_path
                      ? ` · ${truncateMiddle(svc.workspace_path, 24)}`
                      : ""}
                  </p>
                )}
              </div>
            </div>
          </li>
        ))}
      </ul>
    </DataPanel>
  );
}
