import { apiFetch, readErrorMessage } from "../../lib/http";
import type { AppSettings } from "./types";

export async function getSettings(): Promise<AppSettings> {
    const res = await apiFetch("/api/v1/settings");
    if (!res.ok) {
        throw new Error(await readErrorMessage(res, "Failed to load settings"));
    }
    return res.json();
}

export async function updateSettings(
    registrationOpen: boolean,
): Promise<AppSettings> {
    const res = await apiFetch("/api/v1/settings", {
        method: "PATCH",
        body: JSON.stringify({ registration_open: registrationOpen }),
    });
    if (!res.ok) {
        throw new Error(await readErrorMessage(res, "Failed to update settings"));
    }
    return res.json()
}