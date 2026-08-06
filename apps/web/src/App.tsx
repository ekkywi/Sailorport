import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import {
  createService,
  deleteService,
  listServices,
  updateService,
} from "./api";
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
  const [editingId, setEditingId] = useState<string | null>(null);

  async function load() {
    setLoading(true);
    setError("");
    try {
      const data = await listServices();
      setServices(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gagal memuat catalog");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void load();
  }, []);

  function resetForm() {
    setName("");
    setDescription("");
    setOwner("");
    setEditingId(null);
  }

  function startEdit(svc: Service) {
    setEditingId(svc.id);
    setName(svc.name);
    setDescription(svc.description);
    setOwner(svc.owner);
    setError("");
  }

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError("");
    try {
      if (editingId) {
        await updateService(editingId, { name, description, owner });
      } else {
        await createService({ name, description, owner });
      }
      resetForm();
      await load();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : editingId
            ? "Gagal mengupdate service"
            : "Gagal membuat service",
      );
    } finally {
      setSaving(false);
    }
  }

  async function onDelete(svc: Service) {
    const ok = window.confirm(`Hapus service "${svc.name}"?`);
    if (!ok) {
      return;
    }
    setError("");
    try {
      await deleteService(svc.id);
      if (editingId === svc.id) {
        resetForm();
      }
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gagal menghapus service");
    }
  }

  return (
    <div className="page">
      <header>
        <h1>Sailorport</h1>
        <p>Software catalog</p>
      </header>

      <section className="card">
        <h2>{editingId ? "Edit service" : "Tambah service"}</h2>
        <form onSubmit={onSubmit} className="form">
          <label>
            Name
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              placeholder="payments-api"
            />
          </label>
          <label>
            Owner
            <input
              value={owner}
              onChange={(e) => setOwner(e.target.value)}
              placeholder="platform-team"
            />
          </label>
          <label>
            Description
            <input
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Handles payments"
            />
          </label>
          <div className="form-actions">
            <button type="submit" disabled={saving}>
              {saving ? "Menyimpan..." : editingId ? "Update" : "Create"}
            </button>
            {editingId && (
              <button
                type="button"
                className="button-secondary"
                onClick={resetForm}
                disabled={saving}
              >
                Batal
              </button>
            )}
          </div>
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
              <div className="item-main">
                <strong>{svc.name}</strong>
                <span className="muted">
                  {svc.owner || "no owner"} ·{" "}
                  {svc.description || "no description"}
                </span>
              </div>
              <div className="item-actions">
                <button
                  type="button"
                  className="button-secondary"
                  onClick={() => startEdit(svc)}
                >
                  Edit
                </button>
                <button
                  type="button"
                  className="button-danger"
                  onClick={() => void onDelete(svc)}
                >
                  Hapus
                </button>
              </div>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}

export default App;