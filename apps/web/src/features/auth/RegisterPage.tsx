import type { FormEvent } from "react";
import { useState } from "react";
import { Link } from "react-router-dom";
import { AlertCircle, Loader2 } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { login, register } from "./api";
import { PasswordField } from "./PasswordField";
import { authFieldClass, authLabelClass, authSubmitClass } from "./styles";

type RegisterPageProps = {
  onSuccess: () => void;
};

export function RegisterPage({ onSuccess }: RegisterPageProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

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

  return (
    <div className="space-y-6">
      {error && (
        <Alert variant="destructive">
          <AlertCircle />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

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
            aria-invalid={error ? true : undefined}
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
            aria-invalid={error ? true : undefined}
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
