import { useEffect, useState, type ReactNode } from "react";
import {
  Navigate,
  Route,
  Routes,
  useNavigate,
} from "react-router-dom";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/ThemeToggle";
import { CatalogPage } from "./features/catalog/CatalogPage";
import { LoginPage } from "./features/auth/LoginPage";
import { RegisterPage } from "./features/auth/RegisterPage";
import { logout, me } from "./features/auth/api";
import type { AuthUser } from "./features/auth/types";
import { ScaffoldPanel } from "./features/scaffold/ScaffoldPanel";
import { AuthLayout } from "./layouts/AuthLayout";
import { getToken } from "./lib/http";
import "./styles/app.css";

function SessionGate({
  children,
}: {
  children: (ctx: {
    user: AuthUser | null;
    checking: boolean;
    reload: () => Promise<void>;
    signOut: () => void;
  }) => ReactNode;
}) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [checking, setChecking] = useState(true);

  async function reload() {
    setChecking(true);
    if (!getToken()) {
      setUser(null);
      setChecking(false);
      return;
    }
    try {
      setUser(await me());
    } catch {
      setUser(null);
    } finally {
      setChecking(false);
    }
  }

  useEffect(() => {
    void reload();
  }, []);

  function signOut() {
    logout();
    setUser(null);
  }

  return <>{children({ user, checking, reload, signOut })}</>;
}

function Dashboard({
  user,
  onLogout,
}: {
  user: AuthUser;
  onLogout: () => void;
}) {
  const [catalogTick, setCatalogTick] = useState(0);

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
          <div className="row" style={{ gap: "0.5rem" }}>
            <ThemeToggle />
            <Button type="button" variant="outline" onClick={onLogout}>
              Logout
            </Button>
          </div>
        </div>
      </header>
      <ScaffoldPanel onSuccess={() => setCatalogTick((n) => n + 1)} />
      <CatalogPage refreshToken={catalogTick} />
    </div>
  );
}

function LoginRoute({ onSuccess }: { onSuccess: () => void }) {
  const navigate = useNavigate();
  return (
    <AuthLayout mode="login">
      <LoginPage
        onSuccess={() => {
          onSuccess();
          void navigate("/", { replace: true });
        }}
      />
    </AuthLayout>
  );
}

function RegisterRoute({ onSuccess }: { onSuccess: () => void }) {
  const navigate = useNavigate();
  return (
    <AuthLayout mode="register">
      <RegisterPage
        onSuccess={() => {
          onSuccess();
          void navigate("/", { replace: true });
        }}
      />
    </AuthLayout>
  );
}

function App() {
  return (
    <SessionGate>
      {({ user, checking, reload, signOut }) => {
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
            <Routes>
              <Route
                path="/login"
                element={<LoginRoute onSuccess={() => void reload()} />}
              />
              <Route
                path="/register"
                element={<RegisterRoute onSuccess={() => void reload()} />}
              />
              <Route path="*" element={<Navigate to="/login" replace />} />
            </Routes>
          );
        }

        return (
          <Routes>
            <Route
              path="/"
              element={<Dashboard user={user} onLogout={signOut} />}
            />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        );
      }}
    </SessionGate>
  );
}

export default App;
