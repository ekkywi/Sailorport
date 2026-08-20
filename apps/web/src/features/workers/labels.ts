import type { Worker } from "./types";

export function workerLabelString(
  labels: Record<string, unknown> | null | undefined,
  key: string,
): string {
  if (!labels) return "";
  const v = labels[key];
  if (v == null) return "";
  return String(v).trim();
}

export function workerEnvironments(worker: Pick<Worker, "labels">): string[] {
  const raw = workerLabelString(worker.labels, "environments");
  if (!raw) return [];

  const seen = new Set<string>();
  const out: string[] = [];
  for (const part of raw.split(",")) {
    const slug = part.trim().toLowerCase();
    if (!slug || seen.has(slug)) continue;
    seen.add(slug);
    out.push(slug);
  }
  return out;
}

export function workerAllowsEnvironment(
  worker: Pick<Worker, "labels">,
  envSlug: string,
): boolean {
  const normalized = envSlug.trim().toLowerCase() || "dev";
  const allowed = workerEnvironments(worker);
  if (allowed.length === 0) return true;
  return allowed.includes(normalized);
}

export function workerTier(worker: Pick<Worker, "labels">): string {
  return workerLabelString(worker.labels, "tier");
}

export function formatWorkerEnvironments(worker: Pick<Worker, "labels">): string {
  const envs = workerEnvironments(worker);
  if (envs.length === 0) return "All";
  return envs.join(", ");
}

export function workerExtraLabels(
  labels: Record<string, unknown> | null | undefined,
): { key: string; value: string }[] {
  if (!labels) return [];
  const hidden = new Set(["role", "tier", "environments"]);
  return Object.entries(labels)
    .filter(([key]) => !hidden.has(key))
    .map(([key, value]) => ({
      key,
      value: value == null ? "" : String(value),
    }));
}
