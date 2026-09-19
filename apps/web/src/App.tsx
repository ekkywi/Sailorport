import { useCallback, useEffect, useState, type ReactNode } from "react";
import {
  Navigate,
  Route,
  Routes,
  useNavigate,
} from "react-router-dom";
import { Loader2 } from "lucide-react";
import { CatalogPage } from "./features/catalog/CatalogPage";
import { LoginPage } from "./features/auth/LoginPage";
import { logout, me } from "./features/auth/api";
import type { AuthUser } from "./features/auth/types";
import { OverviewPage } from "./features/overview/OverviewPage";
import { AuditPage } from "./features/audit/AuditPage";
import { getSetupStatus } from "./features/setup/api";
import { SetupPage } from "./features/setup/SetupPage";
import { UsersPage } from "./features/users/UsersPage";
import { SettingsPage } from "./features/settings/SettingsPage";
import { WorkersPage } from "./features/workers/WorkersPage";
import { AppShell } from "./layouts/AppShell";
import { AuthLayout } from "./layouts/AuthLayout";
import { getToken } from "./lib/http";
import { isAdmin } from "./lib/rbac";

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

function SetupRoute({
  onSuccess,
  onAlreadySetup,
  showSignInLink = true,
}: {
  onSuccess: () => void;
  onAlreadySetup?: () => void;
  showSignInLink?: boolean;
}) {
  const navigate = useNavigate();
  return (
    <AuthLayout mode="setup">
      <SetupPage
        showSignInLink={showSignInLink}
        onAlreadySetup={onAlreadySetup}
        onSuccess={() => {
          onSuccess();
          void navigate("/overview", { replace: true });
        }}
      />
    </AuthLayout>
  );
}

function App() {
  const [setupChecking, setSetupChecking] = useState(true);
  const [needsSetup, setNeedsSetup] = useState(false);
  const [setupError, setSetupError] = useState("");

  const loadSetup = useCallback(async () => {
    setSetupChecking(true);
    setSetupError("");
    try {
      const st = await getSetupStatus();
      setNeedsSetup(st.needs_setup);
      setSetupError("");
    } catch (err) {
      setSetupError(
        err instanceof Error ? err.message : "Failed to load setup status",
      );
    } finally {
      setSetupChecking(false);
    }
  }, []);

  useEffect(() => {
    void loadSetup();
  }, [loadSetup]);

  const clearSetupGate = useCallback(() => {
    setNeedsSetup(false);
  }, []);

  if (setupChecking) {
    return (
      <div className="flex min-h-svh items-center justify-center bg-background">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Loader2 className="size-4 animate-spin" />
          Checking installation...
        </div>
      </div>
    );
  }

  if (setupError) {
    return (
      <div className="flex min-h-svh flex-col items-center justify-center gap-3 bg-background px-6">
        <p className="text-sm text-destructive">{setupError}</p>
        <button
          type="button"
          className="text-sm font-medium text-foreground underline"
          onClick={() => void loadSetup()}
        >
          Retry
        </button>
      </div>
    );
  }

  // Instalasi kosong: hanya wizard setup (login tidak bisa diakses dulu).
  if (needsSetup) {
    return (
      <Routes>
        <Route
          path="/setup"
          element={
            <SetupRoute
              showSignInLink={false}
              onAlreadySetup={clearSetupGate}
              onSuccess={clearSetupGate}
            />
          }
        />
        <Route path="*" element={<Navigate to="/setup" replace />} />
      </Routes>
    );
  }

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
                path="/setup"
                element={<SetupRoute onSuccess={() => void reload()} />}
              />
              <Route path="*" element={<Navigate to="/login" replace />} />
            </Routes>
          );
        }

        return (
          <AppShell user={user} onLogout={signOut}>
            <Routes>
              <Route path="/overview" element={<OverviewPage />} />
              <Route
                path="/catalog"
                element={<CatalogPage currentUser={user} />}
              />
              <Route path="/worker" element={<WorkersPage />} />
              <Route
                path="/users"
                element={
                  isAdmin(user.role) ? (
                    <UsersPage />
                  ) : (
                    <Navigate to="/overview" replace />
                  )
                }
              />
              <Route
                path="/settings"
                element={
                  isAdmin(user.role) ? (
                    <SettingsPage />
                  ) : (
                    <Navigate to="/overview" replace />
                  )
                }
              />
              <Route
                path="/audit"
                element={
                  isAdmin(user.role) ? (
                    <AuditPage />
                  ) : (
                    <Navigate to="/overview" replace />
                  )
                }
              />
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
