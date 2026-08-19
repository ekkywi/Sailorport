import { apiFetch, readErrorMessage } from "../../lib/http";
import type { AuditEvent } from "./types";

export async function listAuditEvents(limit = 50): Promise<AuditEvent[]> {
  const res = await apiFetch(`/api/v1/audit?limit=${limit}`);
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `List audit failed: ${res.status}`));
  }
  return res.json();
}
