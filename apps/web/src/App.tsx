import { useEffect, useState } from "react";
import { CatalogPage } from "./features/catalog/CatalogPage";
import { LoginPage } from "./features/auth/LoginPage";
import { logout, me } from "./features/auth/api";
import type { AuthUser } from "./features/auth/types";
import { ScaffoldPanel } from "./features/scaffold/ScaffoldPanel";
import { getToken } from "./lib/http";
import "./styles/app.css";

function App() {
  const [catalogTick, setCatalogTick] = useState(0);
  const [user, setUser] = useState<AuthUser | null>(null);
  const [checking, setChecking] = useState(true);

  async function loadSession() {
    setChecking(true);
    if (!getToken()) {
      setUser(null);
      setChecking(false);
      return;
    }
    try {
      const current = await me();
      setUser(current);
    } catch {
      setUser(null);
    } finally {
      setChecking(false);
    }
  }

  useEffect(() => {
    void loadSession();
  }, []);

  function handleLogout() {
    logout();
    setUser(null);
  }

  if (checking) {
    return (
      <div className="page">
        <p>Memuat...</p>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="page">
        <header>
          <h1>Sailorport</h1>
          <p>Software catalog & golden path</p>
        </header>
        <LoginPage onSuccess={() => void loadSession()} />
      </div>
    );
  }

  return (
    <div className="page">
      <header>
        <div className="row">
          <div>
            <h1>Sailorport</h1>
            <p className="muted">
              {user.name || user.email} · role: {user.role}
            </p>
          </div>
          <button type="button" className="button-secondary" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </header>
      <ScaffoldPanel onSuccess={() => setCatalogTick((n) => n + 1)} />
      <CatalogPage refreshToken={catalogTick} />
    </div>
  );
}

export default App;
