import { useEffect, useState, type ReactNode } from "react";
import {
  Navigate,
  Route,
  Routes,
  useNavigate,
} from "react-router-dom";
import { Loader2 } from "lucide-react";
import { CatalogPage } from "./features/catalog/CatalogPage";
import { LoginPage } from "./features/auth/LoginPage";
import { RegisterPage } from "./features/auth/RegisterPage";
import { logout, me } from "./features/auth/api";
import type { AuthUser } from "./features/auth/types";
import { OverviewPage } from "./features/overview/OverviewPage";
import { WorkersPage } from "./features/workers/WorkersPage";
import { AppShell } from "./layouts/AppShell";
import { AuthLayout } from "./layouts/AuthLayout";
import { getToken } from "./lib/http";

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

function LoginRoute({ onSuccess }: { onSuccess: () => void }) {
  const navigate = useNavigate();
  return (
    <AuthLayout mode="login">
      <LoginPage
        onSuccess={() => {
          onSuccess();
          void navigate("/overview", { replace: true });
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
          void navigate("/overview", { replace: true });
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
          <AppShell user={user} onLogout={signOut}>
            <Routes>
              <Route path="/overview" element={<OverviewPage />} />
              <Route path="/catalog" element={<CatalogPage />} />
              <Route path="/worker" element={<WorkersPage />} />
              <Route path="/" element={<Navigate to="/overview" replace />} />
              <Route path="*" element={<Navigate to="/overview" replace />} />
            </Routes>
          </AppShell>
        );
      }}
    </SessionGate>
  );
}

export default App;
