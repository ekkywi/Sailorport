import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { Link, Navigate } from "react-router-dom";
import { AlertCircle, Loader2 } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { login } from "@/features/auth/api";
import { PasswordField } from "@/features/auth/PasswordField";
import {
  authFieldClass,
  authLabelClass,
  authSubmitClass,
} from "@/features/auth/styles";
import { createSetupAdmin, getSetupStatus } from "./api";

type SetupPageProps = {
  onSuccess: () => void;
  /** Dipanggil jika API bilang setup sudah selesai (hindari loop redirect di gate). */
  onAlreadySetup?: () => void;
  showSignInLink?: boolean;
};

export function SetupPage({
  onSuccess,
  onAlreadySetup,
  showSignInLink = true,
}: SetupPageProps) {
  const [checking, setChecking] = useState(true);
  const [needsSetup, setNeedsSetup] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;
    async function load() {
      try {
        const st = await getSetupStatus();
        if (!cancelled) {
          setNeedsSetup(st.needs_setup);
        }
      } catch (err) {
        if (!cancelled) {
          setError(
            err instanceof Error ? err.message : "Failed to load setup status",
          );
        }
      } finally {
        if (!cancelled) {
          setChecking(false);
        }
      }
    }
    void load();
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    if (checking || error || needsSetup) {
      return;
    }
    onAlreadySetup?.();
  }, [checking, error, needsSetup, onAlreadySetup]);

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = new FormData(e.currentTarget);
    const name = String(form.get("name") ?? "");
    const email = String(form.get("email") ?? "");
    const password = String(form.get("password") ?? "");
    const confirmPassword = String(form.get("confirmPassword") ?? "");

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    setLoading(true);
    setError("");
    try {
      await createSetupAdmin(email, password, name);
      await login(email, password);
      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Setup failed");
    } finally {
      setLoading(false);
    }
  }

  if (checking) {
    return (
      <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
        <Loader2 className="size-4 animate-spin" />
        Checking installation...
      </div>
    );
  }

  // Instalasi sudah punya admin.
  if (!needsSetup && !error) {
    if (onAlreadySetup) {
      return (
        <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
          <Loader2 className="size-4 animate-spin" />
          Opening sign in...
        </div>
      );
    }
    return <Navigate to="/login" replace />;
  }

  return (
    <div className="space-y-6">
      {error ? (
        <Alert variant="destructive">
          <AlertCircle />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      {needsSetup ? (
        <form onSubmit={onSubmit} className="space-y-5">
          <div className="space-y-1.5">
            <Label htmlFor="setup-name" className={authLabelClass}>
              Name
            </Label>
            <Input
              id="setup-name"
              name="name"
              autoComplete="name"
              autoFocus
              placeholder="Your name"
              disabled={loading}
              aria-invalid={error ? true : undefined}
              className={authFieldClass}
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="setup-email" className={authLabelClass}>
              Email
            </Label>
            <Input
              id="setup-email"
              name="email"
              type="email"
              autoComplete="email"
              required
              placeholder="you@company.com"
              disabled={loading}
              aria-invalid={error ? true : undefined}
              className={authFieldClass}
            />
          </div>

          <PasswordField
            id="setup-password"
            name="password"
            label="Password"
            autoComplete="new-password"
            disabled={loading}
            invalid={Boolean(error)}
          />

          <PasswordField
            id="setup-confirm-password"
            name="confirmPassword"
            label="Confirm password"
            autoComplete="new-password"
            placeholder="Re-enter your password"
            disabled={loading}
            invalid={Boolean(error)}
          />

          <Button type="submit" className={authSubmitClass} disabled={loading}>
            {loading ? (
              <>
                <Loader2 className="size-3.5 animate-spin" />
                Creating admin...
              </>
            ) : (
              "Create admin account"
            )}
          </Button>
        </form>
      ) : null}

      {showSignInLink ? (
        <p className="text-center text-[13px] text-muted-foreground">
          Already set up?{" "}
          <Link
            to="/login"
            className="font-medium text-foreground transition-colors hover:text-foreground/80"
          >
            Sign in
          </Link>
        </p>
      ) : null}
    </div>
  );
}
