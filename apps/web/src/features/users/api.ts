import { apiFetch, readErrorMessage } from "../../lib/http";
import type { User, UserRole } from "./type";

export async function listUsers(): Promise<User[]> {
    const res = await apiFetch("/api/v1/users");
    if (!res.ok) {
        throw new Error(await readErrorMessage(res, `Failed to list users: ${res.status}`));
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