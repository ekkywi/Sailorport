import { apiFetch, readErrorMessage } from "../../lib/http";
import type { ScaffoldInput, ScaffoldResult, TemplateManifest } from "./types";

export async function listTemplates(): Promise<TemplateManifest[]> {
  const res = await apiFetch("/api/v1/templates");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal list templates: ${res.status}`));
  }
  return res.json();
}

export async function scaffoldService(input: ScaffoldInput): Promise<ScaffoldResult> {
  const res = await apiFetch("/api/v1/scaffold", {
    method: "POST",
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal scaffold: ${res.status}`));
  }
  return res.json();
}
