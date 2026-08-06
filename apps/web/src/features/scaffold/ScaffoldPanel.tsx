import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { listTemplates, scaffoldService } from "./api";
import type { TemplateManifest } from "./types";

type ScaffoldPanelProps = {
  onSuccess: () => void;
};

export function ScaffoldPanel({ onSuccess }: ScaffoldPanelProps) {
  const [templates, setTemplates] = useState<TemplateManifest[]>([]);
  const [templateId, setTemplateId] = useState("");
  const [name, setName] = useState("");
  const [owner, setOwner] = useState("");
  const [description, setDescription] = useState("");
  const [loadingTemplates, setLoadingTemplates] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [successPath, setSuccessPath] = useState("");

  useEffect(() => {
    async function load() {
      setLoadingTemplates(true);
      setError("");
      try {
        const data = await listTemplates();
        setTemplates(data);
        if (data.length > 0) {
          setTemplateId(data[0].id);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Gagal memuat templates");
      } finally {
        setLoadingTemplates(false);
      }
    }
    void load();
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError("");
    setSuccessPath("");
    try {
      const result = await scaffoldService({
        template_id: templateId,
        name,
        owner,
        description,
      });
      setSuccessPath(result.service.workspace_path);
      setName("");
      setOwner("");
      setDescription("");
      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gagal scaffold");
    } finally {
      setSaving(false);
    }
  }

  return (
    <section className="card">
      <h2>Create from template</h2>
      <p className="muted">
        Golden path: pilih template, Sailorport generate folder + daftar ke catalog.
      </p>

      {loadingTemplates && <p>Memuat templates...</p>}
      {error && <p className="error">{error}</p>}
      {successPath && (
        <p className="success">
          Berhasil. Workspace: <code>{successPath}</code>
        </p>
      )}

      {!loadingTemplates && templates.length === 0 && (
        <p className="muted">Tidak ada template. Cek folder `templates/` dan env SAILORPORT_TEMPLATES.</p>
      )}

      {templates.length > 0 && (
        <form onSubmit={onSubmit} className="form">
          <label>
            Template
            <select
              value={templateId}
              onChange={(e) => setTemplateId(e.target.value)}
              required
            >
              {templates.map((t) => (
                <option key={t.id} value={t.id}>
                  {t.name} ({t.id})
                </option>
              ))}
            </select>
          </label>
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
              placeholder="Opsional — default dari template"
            />
          </label>
          <button type="submit" disabled={saving}>
            {saving ? "Scaffolding..." : "Scaffold"}
          </button>
        </form>
      )}
    </section>
  );
}
