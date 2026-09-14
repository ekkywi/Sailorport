import type { FormEvent } from "react";
import { useState } from "react";
import { AlertCircle, Loader2 } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { login } from "./api";
import { PasswordField } from "./PasswordField";
import { authFieldClass, authLabelClass, authSubmitClass } from "./styles";

type LoginPageProps = {
  onSuccess: () => void;
};

export function LoginPage({ onSuccess }: LoginPageProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = new FormData(e.currentTarget);
    const email = String(form.get("email") ?? "");
    const password = String(form.get("password") ?? "");

    setLoading(true);
    setError("");
    try {
      await login(email, password);
      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      {error ? (
        <Alert variant="destructive">
          <AlertCircle />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      <form onSubmit={onSubmit} className="space-y-5">
        <div className="space-y-1.5">
          <Label htmlFor="login-email" className={authLabelClass}>
            Email
          </Label>
          <Input
            id="login-email"
            name="email"
            type="email"
            autoComplete="email"
            autoFocus
            required
            placeholder="you@company.com"
            disabled={loading}
            aria-invalid={error ? true : undefined}
            className={authFieldClass}
          />
        </div>

        <PasswordField
          id="login-password"
          autoComplete="current-password"
          disabled={loading}
          invalid={Boolean(error)}
        />

        <Button type="submit" className={authSubmitClass} disabled={loading}>
          {loading ? (
            <>
              <Loader2 className="size-3.5 animate-spin" />
              Signing in...
            </>
          ) : (
            "Continue"
          )}
        </Button>
      </form>
    </div>
  );
}
