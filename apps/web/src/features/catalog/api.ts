import { apiFetch, readErrorMessage } from "../../lib/http";
import type { CatalogApp, CreateServiceInput, Service, UpdateServiceInput } from "./types";

export async function listServices(): Promise<Service[]> {
  const res = await apiFetch("/api/v1/services");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to list services: ${res.status}`));
  }
  return res.json();
}

export async function createService(input: CreateServiceInput): Promise<Service> {
  const res = await apiFetch("/api/v1/services", {
    method: "POST",
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to create service: ${res.status}`));
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
    throw new Error(await readErrorMessage(res, `Failed to update service: ${res.status}`));
  }
  return res.json();
}

export async function deleteService(id: string): Promise<void> {
  const res = await apiFetch(`/api/v1/services/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to delete service: ${res.status}`));
  }
}

export async function listCatalogApps(): Promise<CatalogApp[]> {
  const res = await apiFetch("/api/v1/catalog-apps");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to list catalog apps: ${res.status}`));
  }
  return res.json();
}
