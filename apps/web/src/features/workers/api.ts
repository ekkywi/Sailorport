import { apiFetch, readErrorMessage } from "../../lib/http";
import type { Worker } from "./types";

export async function listWorkers(): Promise<Worker[]> {
  const res = await apiFetch("/api/v1/workers");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal list workers: ${res.status}`));
  }
  return res.json();
}
