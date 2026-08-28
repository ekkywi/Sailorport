import { apiFetch, readErrorMessage } from "../../lib/http";
import type { CatalogApp, CreateServiceInput, Service, UpdateServiceInput } from "./types";

export async function listServices(): Promise<Service[]> {
  const res = await apiFetch("/api/v1/services");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal list services: ${res.status}`));
  }
  return res.json();
}

export async function createService(input: CreateServiceInput): Promise<Service> {
  const res = await apiFetch("/api/v1/services", {
    method: "POST",
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal create: ${res.status}`));
  }
  return res.json();
}

export async function updateService(
  id: string,
  input: UpdateServiceInput,
): Promise<Service> {
  const res = await apiFetch(`/api/v1/services/${id}`, {
    method: "PUT",
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal update: ${res.status}`));
  }
  return res.json();
}

export async function deleteService(id: string): Promise<void> {
  const res = await apiFetch(`/api/v1/services/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal delete: ${res.status}`));
  }
}

export async function listCatalogApps(): Promise<CatalogApp[]> {
  const res = await apiFetch("/api/v1/catalog-apps");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal list catalog apps: ${res.status}`));
  }
  return res.json();
}
