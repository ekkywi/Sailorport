import type { FormEvent } from "react";
import { useState } from "react";
import { login, register } from "./api";

type LoginPageProps = {
  onSuccess: () => void;
};

export function LoginPage({ onSuccess }: LoginPageProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [mode, setMode] = useState<"login" | "register">("login");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
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
      setError(err instanceof Error ? err.message : "Gagal auth");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="card">
      <h2>{mode === "login" ? "Login" : "Register"}</h2>
      <p className="muted">Masuk untuk mengakses catalog dan scaffold.</p>

      {error && <p className="error">{error}</p>}

      <form onSubmit={onSubmit} className="form">
        {mode === "register" && (
          <label>
            Name
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Your name"
            />
          </label>
        )}
        <label>
          Email
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            placeholder="dev@example.com"
          />
        </label>
        <label>
          Password
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
            placeholder="min 8 characters"
          />
        </label>
        <div className="form-actions">
          <button type="submit" disabled={loading}>
            {loading ? "Loading..." : mode === "login" ? "Login" : "Register"}
          </button>
          <button
            type="button"
            className="button-secondary"
            onClick={() => {
              setMode(mode === "login" ? "register" : "login");
              setError("");
            }}
            disabled={loading}
          >
            {mode === "login" ? "Buat akun" : "Sudah punya akun?"}
          </button>
        </div>
      </form>
    </section>
  );
}
