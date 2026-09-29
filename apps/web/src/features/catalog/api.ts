import { apiFetch, readErrorMessage } from "../../lib/http";
import type { CatalogApp, CreateServiceInput, Service, UpdateServiceInput } from "./types";

/** Pass C (C-M2): never keep secrets in client state even if API mis-redacts. */
function redactService(svc: Service): Service {
  return { ...svc, webhook_secret: "", git_token: "" };
}

function redactServices(list: Service[]): Service[] {
  return list.map(redactService);
}

export async function listServices(): Promise<Service[]> {
  const res = await apiFetch("/api/v1/services");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to list services: ${res.status}`));
  }
  const data = (await res.json()) as Service[];
  return redactServices(data);
}

export async function createService(input: CreateServiceInput): Promise<Service> {
  const res = await apiFetch("/api/v1/services", {
    method: "POST",
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to create service: ${res.status}`));
  }
  return redactService((await res.json()) as Service);
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
  return redactService((await res.json()) as Service);
}

export async function deleteService(id: string): Promise<void> {
  const res = await apiFetch(`/api/v1/services/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to delete service: ${res.status}`));
  }
}

export async function transferService(
  id: string,
  input: { email: string },
): Promise<Service> {
  const res = await apiFetch(`/api/v1/services/${id}/transfer`, {
    method: "POST",
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to transfer: ${res.status}`));
  }
  return redactService((await res.json()) as Service);
}

export async function listCatalogApps(): Promise<CatalogApp[]> {
  const res = await apiFetch("/api/v1/catalog-apps");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to list catalog apps: ${res.status}`));
  }
  return res.json();
}
