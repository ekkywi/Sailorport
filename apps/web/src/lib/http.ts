const TOKEN_KEY = "sailorport_token";

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

export async function readErrorMessage(res: Response, fallback: string): Promise<string> {
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

/** Paths where 401 means bad credentials / closed gate — not an expired session. */
function isAuthChallengePath(path: string): boolean {
  const p = path.split("?")[0];
  return (
    p === "/api/v1/auth/login" ||
    p === "/api/v1/auth/register" ||
    p === "/api/v1/auth/registration-status" ||
    p.startsWith("/api/v1/setup/")
  );
}

function redirectToLogin() {
  const path = window.location.pathname;
  if (path === "/login" || path === "/register" || path === "/setup") {
    return;
  }
  window.location.assign("/login");
}

export async function apiFetch(path: string, init: RequestInit = {}): Promise<Response> {
  const headers = new Headers(init.headers);
  if (!headers.has("Content-Type") && init.body) {
    headers.set("Content-Type", "application/json");
  }
  const token = getToken();
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  const res = await fetch(path, { ...init, headers });

  // Pass C (C-H1): drop stale session and return to login on authenticated 401.
  if (res.status === 401 && token && !isAuthChallengePath(path)) {
    clearToken();
    redirectToLogin();
  }

  return res;
}
