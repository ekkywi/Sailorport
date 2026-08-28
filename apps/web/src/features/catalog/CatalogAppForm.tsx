import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { listCatalogApps } from "./api";
import type { CatalogApp, CatalogAppFormValues } from "./types";

type CatalogAppFormProps = {
  values: CatalogAppFormValues;
  saving: boolean;
  error?: string;
  onChange: (field: keyof CatalogAppFormValues, value: string) => void;
  onSubmit: (e: FormEvent) => void;
  onCancel: () => void;
  onBack?: () => void;
};

export function CatalogAppForm({
  values,
  saving,
  error,
  onChange,
  onSubmit,
  onCancel,
  onBack,
}: CatalogAppFormProps) {
  const [apps, setApps] = useState<CatalogApp[]>([]);
  const [loadingApps, setLoadingApps] = useState(true);
  const [loadError, setLoadError] = useState("");

  useEffect(() => {
    async function load() {
      setLoadingApps(true);
      setLoadError("");
      try {
        const data = await listCatalogApps();
        setApps(data);

        if (data.length > 0 && !values.catalog_app_id) {
          onChange("catalog_app_id", data[0].id);
        }
      } catch (err) {
        setLoadError(
          err instanceof Error ? err.message : "Failed to load catalog apps",
        );
      } finally {
        setLoadingApps(false);
      }
    }
    void load();
    // eslint-disable-next-line react-hooks/exhaustive-deps -- load once on mount
  }, []);

  if (loadingApps) {
    return (
      <p className="py-6 text-center text-[13px] text-muted-foreground">
        Loading catalog apps…
      </p>
    );
  }

  if (apps.length === 0) {
    return (
      <div className="space-y-4 py-2">
        <p className="text-[13px] text-muted-foreground">
          No catalog apps configured. Check the API{" "}
          <code className="text-[11px]">catalog-apps/</code> folder.
        </p>
        {loadError ? (
          <p className="text-[13px] text-destructive">{loadError}</p>
        ) : null}
        {onBack ? (
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 text-[13px]"
            onClick={onBack}
          >
            Back
          </Button>
        ) : null}
      </div>
    );
  }

  const selected = apps.find((a) => a.id === values.catalog_app_id);

  return (
    <form onSubmit={onSubmit} className="grid gap-3 sm:grid-cols-2">
      {error || loadError ? (
        <p className="text-[13px] text-destructive sm:col-span-2">
          {error || loadError}
        </p>
      ) : null}

      <div className="space-y-1.5 sm:col-span-2">
        <Label htmlFor="catalog-app" className="text-[12px] text-muted-foreground">
          Catalog app
        </Label>
        <select
          id="catalog-app"
          value={values.catalog_app_id}
          onChange={(e) => {
            const id = e.target.value;
            onChange("catalog_app_id", id);
            if (!values.name.trim()) {
              onChange("name", id);
            }
          }}
          required
          className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 text-[13px] shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
        >
          {apps.map((a) => (
            <option key={a.id} value={a.id}>
              {a.name} ({a.id})
            </option>
          ))}
        </select>
        {selected ? (
          <p className="text-[11px] text-muted-foreground">
            {selected.description || selected.image}
            {selected.image ? (
              <>
                {" "}
                · <span className="font-mono">{selected.image}</span>
                {selected.container_port
                  ? ` · port ${selected.container_port}`
                  : null}
              </>
            ) : null}
          </p>
        ) : null}
      </div>

      <div className="space-y-1.5 sm:col-span-2">
        <Label htmlFor="catalog-name" className="text-[12px] text-muted-foreground">
          Service name
        </Label>
        <Input
          id="catalog-name"
          value={values.name}
          onChange={(e) => onChange("name", e.target.value)}
          required
          placeholder="demo-pg"
          className="h-9 text-[13px]"
        />
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="catalog-owner" className="text-[12px] text-muted-foreground">
          Owner
        </Label>
        <Input
          id="catalog-owner"
          value={values.owner}
          onChange={(e) => onChange("owner", e.target.value)}
          placeholder="platform-team"
          className="h-9 text-[13px]"
        />
      </div>

      <div className="space-y-1.5">
        <Label
          htmlFor="catalog-description"
          className="text-[12px] text-muted-foreground"
        >
          Description
        </Label>
        <Input
          id="catalog-description"
          value={values.description}
          onChange={(e) => onChange("description", e.target.value)}
          placeholder="Optional"
          className="h-9 text-[13px]"
        />
      </div>

      <div className="flex flex-wrap items-center gap-3 pt-1 sm:col-span-2">
        <Button type="submit" size="sm" className="h-8 text-[13px]" disabled={saving}>
          {saving ? "Adding…" : "Add service"}
        </Button>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="h-8 text-[13px]"
          onClick={onCancel}
          disabled={saving}
        >
          Cancel
        </Button>
        {onBack ? (
          <button
            type="button"
            className="text-[12px] text-muted-foreground transition-colors hover:text-foreground"
            onClick={onBack}
            disabled={saving}
          >
            Back
          </button>
        ) : null}
      </div>
    </form>
  );
}
