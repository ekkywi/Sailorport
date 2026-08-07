const STORAGE_KEY = "sailorport-theme";

export type Theme = "light" | "dark";

export function getStoredTheme(): Theme {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === "light" || stored === "dark") return stored;
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? "dark"
    : "light";
}

export function applyTheme(theme: Theme) {
  document.documentElement.classList.toggle("dark", theme === "dark");
  localStorage.setItem(STORAGE_KEY, theme);
}

export function toggleTheme(): Theme {
  const next: Theme = document.documentElement.classList.contains("dark")
    ? "light"
    : "dark";
  applyTheme(next);
  return next;
}

export function initTheme() {
  applyTheme(getStoredTheme());
}
