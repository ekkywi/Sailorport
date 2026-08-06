import type { FormEvent } from "react";
import { useState } from "react";
import { AlertCircle, Eye, EyeOff, Loader2, Ship } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { login, register } from "./api";

type LoginPageProps = {
  onSuccess: () => void;
};

type AuthMode = "login" | "register";

function AuthForm({
  mode,
  loading,
  onSubmit,
}: {
  mode: AuthMode;
  loading: boolean;
  onSubmit: (e: FormEvent<HTMLFormElement>) => void;
}) {
  const [showPassword, setShowPassword] = useState(false);

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      {mode === "register" && (
        <div className="space-y-2">
          <Label htmlFor="name">Name</Label>
          <Input
            id="name"
            name="name"
            autoComplete="name"
            placeholder="Your name"
            disabled={loading}
          />
        </div>
      )}

      <div className="space-y-2">
        <Label htmlFor={`${mode}-email`}>Email</Label>
        <Input
          id={`${mode}-email`}
          name="email"
          type="email"
          autoComplete="email"
          required
          placeholder="dev@example.com"
          disabled={loading}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${mode}-password`}>Password</Label>
        <div className="relative">
          <Input
            id={`${mode}-password`}
            name="password"
            type={showPassword ? "text" : "password"}
            autoComplete={mode === "login" ? "current-password" : "new-password"}
            required
            minLength={8}
            placeholder="Min. 8 characters"
            disabled={loading}
            className="pr-10"
          />
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            className="absolute top-1/2 right-1 -translate-y-1/2 text-muted-foreground"
            onClick={() => setShowPassword((v) => !v)}
            disabled={loading}
            aria-label={showPassword ? "Hide password" : "Show password"}
          >
            {showPassword ? <EyeOff /> : <Eye />}
          </Button>
        </div>
      </div>

      <Button type="submit" className="w-full" disabled={loading}>
        {loading ? (
          <>
            <Loader2 className="animate-spin" />
            Processing...
          </>
        ) : mode === "login" ? (
          "Sign in"
        ) : (
          "Create account"
        )}
      </Button>
    </form>
  );
}

export function LoginPage({ onSuccess }: LoginPageProps) {
  const [mode, setMode] = useState<AuthMode>("login");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = new FormData(e.currentTarget);
    const email = String(form.get("email") ?? "");
    const password = String(form.get("password") ?? "");
    const name = String(form.get("name") ?? "");

    setLoading(true);
    setError("");
    try {
      if (mode === "login") {
        await login(email, password);
      } else {
        await register(email, password, name);
        await login(email, password);
      }
      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Authentication failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      <div className="space-y-2 text-center lg:hidden">
        <div className="mx-auto flex size-11 items-center justify-center rounded-xl bg-teal-700 text-white">
          <Ship className="size-5" />
        </div>
        <h1 className="text-2xl font-semibold tracking-tight">Sailorport</h1>
        <p className="text-sm text-muted-foreground">
          Sign in to access catalog and scaffold
        </p>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader className="space-y-1">
          <CardTitle className="text-xl">Welcome back</CardTitle>
          <CardDescription>
            {mode === "login"
              ? "Enter your credentials to continue."
              : "Create an account to get started."}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {error && (
            <Alert variant="destructive">
              <AlertCircle />
              <AlertTitle>Authentication failed</AlertTitle>
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <Tabs
            value={mode}
            onValueChange={(value) => {
              setMode(value as AuthMode);
              setError("");
            }}
          >
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="login">Sign in</TabsTrigger>
              <TabsTrigger value="register">Create account</TabsTrigger>
            </TabsList>
            <TabsContent value="login" className="mt-4">
              <AuthForm mode="login" loading={loading} onSubmit={onSubmit} />
            </TabsContent>
            <TabsContent value="register" className="mt-4">
              <AuthForm mode="register" loading={loading} onSubmit={onSubmit} />
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  );
}
