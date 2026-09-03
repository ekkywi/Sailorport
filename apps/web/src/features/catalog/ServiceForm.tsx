import type { FormEvent } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { Environment } from "../environments/types";
import { CatalogEnvFields } from "./CatalogEnvFields";
import type { CatalogAppEnvField, ServiceFormValues } from "./types";
import { WebhookSettingsFields } from "./WebhookSettingsFields";

type ServiceFormProps = {
  mode: "register" | "edit";
  values: ServiceFormValues;
  saving: boolean;
  error?: string;
  /** Show GitHub webhook controls (edit + git services). */
  showWebhookSettings?: boolean;
  environments?: Environment[];
  /** Catalog app env schema when editing source_type=catalog_app. */
  catalogEnvFields?: CatalogAppEnvField[];
  catalogEnv?: Record<string, string>;
  onCatalogEnvChange?: (name: string, value: string) => void;
  onChange: (field: keyof ServiceFormValues, value: string | boolean) => void;
  onSubmit: (e: FormEvent) => void;
  onCancel: () => void;
  onBack?: () => void;
};

export function ServiceForm({
  mode,
  values,
  saving,
  error,
  showWebhookSettings = false,
  environments = [],
  catalogEnvFields,
  catalogEnv = {},
  onCatalogEnvChange,
  onChange,
  onSubmit,
  onCancel,
  onBack,
}: ServiceFormProps) {
  const isEdit = mode === "edit";
  const showCatalogEnv =
    isEdit &&
    Boolean(catalogEnvFields?.length) &&
    typeof onCatalogEnvChange === "function";

  return (
    <form onSubmit={onSubmit} className="grid gap-3 sm:grid-cols-2">
      {error ? (
        <p className="text-[13px] text-destructive sm:col-span-2">{error}</p>
      ) : null}

      <div className="space-y-1.5 sm:col-span-2">
        <Label htmlFor="svc-name" className="text-[12px] text-muted-foreground">
          Name
        </Label>
        <Input
          id="svc-name"
          value={values.name}
          onChange={(e) => onChange("name", e.target.value)}
          required
          placeholder="payments-api"
          className="h-9 text-[13px]"
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="svc-owner" className="text-[12px] text-muted-foreground">
          Owner
        </Label>
        <Input
          id="svc-owner"
          value={values.owner}
          onChange={(e) => onChange("owner", e.target.value)}
          placeholder="platform-team"
          className="h-9 text-[13px]"
        />
      </div>
      <div className="space-y-1.5">
        <Label
          htmlFor="svc-description"
          className="text-[12px] text-muted-foreground"
        >
          Description
        </Label>
        <Input
          id="svc-description"
          value={values.description}
          onChange={(e) => onChange("description", e.target.value)}
          placeholder="Handles payments"
          className="h-9 text-[13px]"
        />
      </div>

      {showWebhookSettings ? (
        <WebhookSettingsFields
          values={values}
          environments={environments}
          disabled={saving}
          onChange={onChange}
        />
      ) : null}

      {showCatalogEnv ? (
        <CatalogEnvFields
          fields={catalogEnvFields!}
          values={catalogEnv}
          mode="edit"
          onChange={onCatalogEnvChange!}
        />
      ) : null}

      <div className="flex flex-wrap items-center gap-3 pt-1 sm:col-span-2">
        <Button type="submit" size="sm" className="h-8 text-[13px]" disabled={saving}>
          {saving
            ? "Saving…"
            : isEdit
              ? "Save changes"
              : "Register in catalog"}
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
