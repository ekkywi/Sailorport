import type { FormEvent } from "react";
import type { ServiceFormValues } from "./types";

type ServiceFormProps = {
  title: string;
  values: ServiceFormValues;
  saving: boolean;
  isEditing: boolean;
  onChange: (field: keyof ServiceFormValues, value: string) => void;
  onSubmit: (e: FormEvent) => void;
  onCancel: () => void;
};

export function ServiceForm({
  title,
  values,
  saving,
  isEditing,
  onChange,
  onSubmit,
  onCancel,
}: ServiceFormProps) {
  return (
    <section className="card">
      <h2>{title}</h2>
      <form onSubmit={onSubmit} className="form">
        <label>
          Name
          <input
            value={values.name}
            onChange={(e) => onChange("name", e.target.value)}
            required
            placeholder="payments-api"
          />
        </label>
        <label>
          Owner
          <input
            value={values.owner}
            onChange={(e) => onChange("owner", e.target.value)}
            placeholder="platform-team"
          />
        </label>
        <label>
          Description
          <input
            value={values.description}
            onChange={(e) => onChange("description", e.target.value)}
            placeholder="Handles payments"
          />
        </label>
        <div className="form-actions">
          <button type="submit" disabled={saving}>
            {saving ? "Menyimpan..." : isEditing ? "Update" : "Create"}
          </button>
          {isEditing && (
            <button
              type="button"
              className="button-secondary"
              onClick={onCancel}
              disabled={saving}
            >
              Batal
            </button>
          )}
        </div>
      </form>
    </section>
  );
}
