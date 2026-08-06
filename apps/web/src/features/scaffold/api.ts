import type { ScaffoldInput, ScaffoldResult, TemplateManifest } from "./types";

async function readErrorMessage(res: Response, fallback: string): Promise<string> {
  const text = await res.text();
  if (!text) {
    return fallback;
  }
  try {
    const body = JSON.parse(text) as { error?: string };
    if (body.error) {
      return body.error;
    }
  } catch {
    // plain text
  }
  return text;
}

export async function listTemplates(): Promise<TemplateManifest[]> {
  const res = await fetch("/api/v1/templates");
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal list templates: ${res.status}`));
  }
  return res.json();
}

export async function scaffoldService(input: ScaffoldInput): Promise<ScaffoldResult> {
  const res = await fetch("/api/v1/scaffold", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Gagal scaffold: ${res.status}`));
  }
  return res.json();
}
