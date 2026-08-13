import { apiFetch, readErrorMessage } from "../../lib/http";
import type { CreateUserInput, User, UserRole } from "./type";

export async function listUsers(): Promise<User[]> {
  const res = await apiFetch("/api/v1/users");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to list users: ${res.status}`));
  }
  return res.json();
}

export async function createUser(input: CreateUserInput): Promise<User> {
  const res = await apiFetch("/api/v1/users", {
    method: "POST",
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to create user: ${res.status}`));
  }
  return res.json();
}

export async function updateUserRole(id: string, role: UserRole): Promise<User> {
  const res = await apiFetch(`/api/v1/users/${id}`, {
    method: "PATCH",
    body: JSON.stringify({ role }),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Failed to update role: ${res.status}`));
  }
  return res.json();
}

export async function setUserDisabled(
  id: string,
  disabled: boolean,
): Promise<User> {
  const res = await apiFetch(`/api/v1/users/${id}`, {
    method: "PATCH",
    body: JSON.stringify({ disabled }),
  });
  if (!res.ok) {
    throw new Error(
      await readErrorMessage(res, `Failed to update user: ${res.status}`),
    );
  }
  return res.json();
}

export async function resetUserPassword(
  id: string,
  password: string,
): Promise<void> {
  const res = await apiFetch(`/api/v1/users/${id}/password`, {
    method: "POST",
    body: JSON.stringify({ password }),
  });
  if (!res.ok) {
    throw new Error(
      await readErrorMessage(res, `Failed to reset password: ${res.status}`),
    );
  }
}