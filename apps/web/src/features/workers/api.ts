import { apiFetch, readErrorMessage } from "../../lib/http";
import type { Worker } from "./types";

export async function listWorkers(): Promise<Worker[]> {
  const res = await apiFetch("/api/v1/workers");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to list workers: ${res.status}`));
  }
  return res.json();
}

export async function updateWorkerLabels(
  id: string,
  body: { tier: string; environments: string },
): Promise<Worker> {
  const res = await apiFetch(`/api/v1/workers/${id}`, {
    method: "PATCH",
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to update labels: ${res.status}`));
  }
  return res.json();
}

export async function decommissionWorker(id: string): Promise<Worker> {
  const res = await apiFetch(`/api/v1/workers/${id}/decommission`, {
    method: "POST",
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to decommission: ${res.status}`));
  }
  return res.json();
}

export async function restoreWorker(id: string): Promise<Worker> {
  const res = await apiFetch(`/api/v1/workers/${id}/restore`, {
    method: "POST",
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to restore: ${res.status}`));
  }
  return res.json();
}
