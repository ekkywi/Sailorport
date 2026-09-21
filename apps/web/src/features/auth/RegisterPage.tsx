import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { Link, Navigate } from "react-router-dom";
import { AlertCircle, Loader2 } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { getRegistrationStatus, login, register } from "./api";
import { PasswordField } from "./PasswordField";
import { authFieldClass, authLabelClass, authSubmitClass } from "./styles";

type RegisterPageProps = {
  onSuccess: () => void;
};

export function RegisterPage({ onSuccess }: RegisterPageProps) {
  const [checking, setChecking] = useState(true);
  const [registrationOpen, setRegistrationOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;
    async function load() {
      try {
        const st = await getRegistrationStatus();
        if (!cancelled) {
          setRegistrationOpen(st.registration_open);
        }
      } catch (err) {
        if (!cancelled) {
          setError(
            err instanceof Error
              ? err.message
              : "Failed to load registration status",
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
      await register(email, password, name);
      await login(email, password);
      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
    } finally {
      setLoading(false);
    }
  }

  if (checking) {
    return (
      <div className="flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground">
        <Loader2 className="size-4 animate-spin" />
        Checking registration…
      </div>
    );
  }

  if (!registrationOpen && !error) {
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

      {registrationOpen ? (
        <form onSubmit={onSubmit} className="space-y-5">
          <div className="space-y-1.5">
            <Label htmlFor="register-name" className={authLabelClass}>
              Name
            </Label>
            <Input
              id="register-name"
              name="name"
              autoComplete="name"
              autoFocus
              placeholder="Your name"
              disabled={loading}
              className={authFieldClass}
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="register-email" className={authLabelClass}>
              Email
            </Label>
            <Input
              id="register-email"
              name="email"
              type="email"
              autoComplete="email"
              required
              placeholder="you@company.com"
              disabled={loading}
              className={authFieldClass}
            />
          </div>

          <PasswordField
            id="register-password"
            name="password"
            label="Password"
            autoComplete="new-password"
            disabled={loading}
            invalid={Boolean(error)}
          />

          <PasswordField
            id="register-confirm-password"
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
                Creating...
              </>
            ) : (
              "Create account"
            )}
          </Button>
        </form>
      ) : null}

      <p className="text-center text-[13px] text-muted-foreground">
        Already have an account?{" "}
        <Link
          to="/login"
          className="font-medium text-foreground transition-colors hover:text-foreground/80"
        >
          Sign in
        </Link>
      </p>
    </div>
  );
}
