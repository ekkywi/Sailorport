import { apiFetch, readErrorMessage } from "../../lib/http";
import type { AuthUser } from "../auth/types";
import type { SetupStatus } from "./types";

export async function getSetupStatus(): Promise<SetupStatus> {
  const res = await apiFetch("/api/v1/setup/status");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, "Failed to load setup status"));
  }
  return res.json();
}

export async function createSetupAdmin(
  email: string,
  password: string,
  name: string,
): Promise<AuthUser> {
  const res = await apiFetch("/api/v1/setup/admin", {
    method: "POST",
    body: JSON.stringify({ email, password, name }),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, "Setup failed"));
  }
  return res.json();
}
