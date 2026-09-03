import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { CatalogAppEnvField } from "./types";

type CatalogEnvFieldsProps = {
  fields: CatalogAppEnvField[];
  values: Record<string, string>;
  /** Edit mode: secret empty = keep existing value on server. */
  mode?: "create" | "edit";
  onChange: (name: string, value: string) => void;
};

export function CatalogEnvFields({
  fields,
  values,
  mode = "create",
  onChange,
}: CatalogEnvFieldsProps) {
  if (fields.length === 0) return null;

  const isEdit = mode === "edit";

  return (
    <div className="space-y-3 border-t border-border/60 pt-3 sm:col-span-2">
      <div>
        <p className="text-[13px] font-medium text-foreground">
          Environment variables
        </p>
        <p className="mt-0.5 text-[12px] text-muted-foreground">
          {isEdit
            ? "Secrets left blank keep the current value. Redeploy after changing env so the container picks them up."
            : "Values for the container image. Secrets are stored securely and not shown again after save."}
        </p>
      </div>
      {fields.map((field) => (
        <div key={field.name} className="space-y-1.5">
          <Label
            htmlFor={`catalog-env-${field.name}`}
            className="text-[12px] text-muted-foreground"
          >
            {field.name}
            {field.required && !isEdit ? " *" : ""}
            {field.secret && isEdit ? " (optional)" : ""}
          </Label>
          <Input
            id={`catalog-env-${field.name}`}
            type={field.secret ? "password" : "text"}
            value={values[field.name] ?? ""}
            onChange={(e) => onChange(field.name, e.target.value)}
            required={field.required && !isEdit}
            placeholder={
              field.secret && isEdit
                ? "Leave blank to keep current"
                : field.default
                  ? `Default: ${field.default}`
                  : field.secret
                    ? "Required"
                    : "Optional"
            }
            className="h-9 text-[13px]"
            autoComplete="off"
          />
          {field.description ? (
            <p className="text-[11px] text-muted-foreground">
              {field.description}
            </p>
          ) : null}
        </div>
      ))}
    </div>
  );
}
