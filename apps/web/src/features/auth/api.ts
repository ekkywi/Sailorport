import { apiFetch, clearToken, readErrorMessage, setToken } from "../../lib/http";
import type { AuthUser, LoginResponse } from "./types";

export async function login(email: string, password: string): Promise<LoginResponse> {
  const res = await apiFetch("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, "Login failed"));
  }
  const data = (await res.json()) as LoginResponse;
  setToken(data.token);
  return data;
}

export async function me(): Promise<AuthUser> {
  const res = await apiFetch("/api/v1/auth/me");
  if (!res.ok) {
    clearToken();
    throw new Error(await readErrorMessage(res, "Session invalid"));
  }
  return res.json();
}

export function logout() {
  clearToken();
}
