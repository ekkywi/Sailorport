import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { createService, listServices } from "./api";
import type { Service } from "./types";
import "./App.css";

function App() {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [owner, setOwner] = useState("");
  const [saving, setSaving] = useState(false);

  async function load() {
    setLoading(true);
    setError("");

    try {
      const data = await listServices();
      setServices(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gagal memuat catalog")
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void load();
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError("");

    try {
      await createService({ name, description, owner });
      setName("");
      setDescription("");
      setOwner("");
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gagal membuat service");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="page">
      <header>
        <h1>Sailorport</h1>
        <p>Software Catalog</p>
      </header>

      <section className="card">
        <h2>Tambah service</h2>
        <form onSubmit={onSubmit} className="form">
          <label>
            Name
            <input value={name} onChange={(e) => setName(e.target.value)} required placeholder="payments-api" />
          </label>
          <label>
            Owner
            <input value={owner} onChange={(e) => setOwner(e.target.value)} placeholder="platform-team" />
          </label>
          <label>
            Description
            <input value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Handles payments" />
          </label>
          <button type="submit" disabled={saving}>
            {saving ? "Menyimpan..." : "Create"}
          </button>
        </form>
      </section>

      <section className="card">
        <div className="row">
          <h2>Services</h2>
          <button type="button" onClick={() => void load()} disabled={loading}>
            Refresh
          </button>
        </div>

        {error && <p className="error">{error}</p>}
        {loading && <p>Memuat...</p>}

        {!loading && services.length === 0 && (
          <p className="muted">Belum ada service. Buat yang pertama di atas.</p>
        )}

        <ul className="list">
          {services.map((svc) => (
            <li key={svc.id}>
              <strong>{svc.name}</strong>
              <span className="muted">
                {svc.owner || "no owner"} · {svc.description || "no description"}
              </span>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}

export default App;