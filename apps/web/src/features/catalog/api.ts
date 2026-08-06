import type { CreateServiceInput, Service, UpdateServiceInput } from "./types";

async function readErrorMessage(res: Response, fallback: string): Promise<string> {
  const text = await res.text();
  if (!text) {
    return fallback;
  }
  try {
    const body = JSON.parse(text) as { error?: string };
    if (body.error) {
      return body.error;
    }
  } catch {
    // response was plain text
  }
  return text;
}

export async function listServices(): Promise<Service[]> {
  const res = await fetch("/api/v1/services");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal list services: ${res.status}`));
  }
  return res.json();
}

export async function createService(input: CreateServiceInput): Promise<Service> {
  const res = await fetch("/api/v1/services", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
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
  const res = await fetch(`/api/v1/services/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal update: ${res.status}`));
  }
  return res.json();
}

export async function deleteService(id: string): Promise<void> {
  const res = await fetch(`/api/v1/services/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal delete: ${res.status}`));
  }
}
