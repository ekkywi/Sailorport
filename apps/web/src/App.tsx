import { useEffect, useState } from "react";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { CatalogPage } from "./features/catalog/CatalogPage";
import { LoginPage } from "./features/auth/LoginPage";
import { logout, me } from "./features/auth/api";
import type { AuthUser } from "./features/auth/types";
import { ScaffoldPanel } from "./features/scaffold/ScaffoldPanel";
import { AuthLayout } from "./layouts/AuthLayout";
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
      <div className="flex min-h-svh items-center justify-center bg-background">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Loader2 className="size-4 animate-spin" />
          Loading session...
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <AuthLayout>
        <LoginPage onSuccess={() => void loadSession()} />
      </AuthLayout>
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
          <Button type="button" variant="outline" onClick={handleLogout}>
            Logout
          </Button>
        </div>
      </header>
      <ScaffoldPanel onSuccess={() => setCatalogTick((n) => n + 1)} />
      <CatalogPage refreshToken={catalogTick} />
    </div>
  );
}

export default App;
