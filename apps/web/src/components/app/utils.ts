import { cn } from "@/lib/utils";

/** Format an ISO timestamp as a short relative string (e.g. "2m ago"). */
export function formatRelativeTime(value: string | null): string {
  if (!value) return "Never";
  const then = new Date(value).getTime();
  if (Number.isNaN(then)) return value;

  const seconds = Math.round((Date.now() - then) / 1000);
  if (seconds < 0) return "just now";
  if (seconds < 60) return `${seconds}s ago`;
  const minutes = Math.round(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.round(minutes / 60);
  if (hours < 48) return `${hours}h ago`;
  const days = Math.round(hours / 24);
  if (days < 30) return `${days}d ago`;
  return new Date(value).toLocaleDateString();
}

export function formatAbsoluteTime(value: string | null): string {
  if (!value) return "";
  try {
    return new Date(value).toLocaleString();
  } catch {
    return value;
  }
}

export function userInitials(nameOrEmail: string): string {
  const trimmed = nameOrEmail.trim();
  if (!trimmed) return "?";
  if (trimmed.includes("@")) {
    return trimmed.slice(0, 2).toUpperCase();
  }
  const parts = trimmed.split(/\s+/).filter(Boolean);
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

export function truncateMiddle(value: string, max = 28): string {
  if (value.length <= max) return value;
  const keep = Math.floor((max - 1) / 2);
  return `${value.slice(0, keep)}…${value.slice(-keep)}`;
}

export function labelEntries(
  labels: Record<string, unknown> | null | undefined,
): { key: string; value: string }[] {
  if (!labels) return [];
  return Object.entries(labels).map(([key, value]) => ({
    key,
    value: value == null ? "" : String(value),
  }));
}

export function skeletonClass(className?: string) {
  return cn("animate-pulse rounded-md bg-muted", className);
}
