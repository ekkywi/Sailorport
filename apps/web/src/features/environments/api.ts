import { apiFetch, readErrorMessage } from "../../lib/http";
import type { Environment } from "./types";

export async function listEnvironments(): Promise<Environment[]> {
  const res = await apiFetch("/api/v1/environments");
  if (!res.ok) {
    throw new Error(
      await readErrorMessage(res, `Failed to list environments: ${res.status}`),
    );
  }
  return res.json();
}
