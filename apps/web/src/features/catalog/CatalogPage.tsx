import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import {
  createService,
  deleteService,
  listServices,
  updateService,
} from "./api";
import { ServiceForm } from "./ServiceForm";
import { ServiceList } from "./ServiceList";
import type { Service, ServiceFormValues } from "./types";

const emptyForm: ServiceFormValues = {
  name: "",
  description: "",
  owner: "",
};

export function CatalogPage() {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [values, setValues] = useState<ServiceFormValues>(emptyForm);
  const [saving, setSaving] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);

  async function load() {
    setLoading(true);
    setError("");
    try {
      const data = await listServices();
      setServices(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gagal memuat catalog");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void load();
  }, []);

  function resetForm() {
    setValues(emptyForm);
    setEditingId(null);
  }

  function startEdit(svc: Service) {
    setEditingId(svc.id);
    setValues({
      name: svc.name,
      description: svc.description,
      owner: svc.owner,
    });
    setError("");
  }

  function onChange(field: keyof ServiceFormValues, value: string) {
    setValues((prev) => ({ ...prev, [field]: value }));
  }

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError("");
    try {
      if (editingId) {
        await updateService(editingId, values);
      } else {
        await createService(values);
      }
      resetForm();
      await load();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : editingId
            ? "Gagal mengupdate service"
            : "Gagal membuat service",
      );
    } finally {
      setSaving(false);
    }
  }

  async function onDelete(svc: Service) {
    const ok = window.confirm(`Hapus service "${svc.name}"?`);
    if (!ok) {
      return;
    }
    setError("");
    try {
      await deleteService(svc.id);
      if (editingId === svc.id) {
        resetForm();
      }
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gagal menghapus service");
    }
  }

  return (
    <>
      <ServiceForm
        title={editingId ? "Edit service" : "Tambah service"}
        values={values}
        saving={saving}
        isEditing={editingId !== null}
        onChange={onChange}
        onSubmit={onSubmit}
        onCancel={resetForm}
      />
      <ServiceList
        services={services}
        loading={loading}
        error={error}
        onRefresh={() => void load()}
        onEdit={startEdit}
        onDelete={(svc) => void onDelete(svc)}
      />
    </>
  );
}
