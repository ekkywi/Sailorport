import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { listTemplates, scaffoldService } from "./api";
import type { TemplateManifest } from "./types";

type CreateServiceFormProps = {
  onSuccess: (workspacePath: string) => void;
  onRegisterExisting?: () => void;
};

export function CreateServiceForm({
  onSuccess,
  onRegisterExisting,
}: CreateServiceFormProps) {
  const [templates, setTemplates] = useState<TemplateManifest[]>([]);
  const [templateId, setTemplateId] = useState("");
  const [name, setName] = useState("");
  const [owner, setOwner] = useState("");
  const [description, setDescription] = useState("");
  const [loadingTemplates, setLoadingTemplates] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

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
        setError(err instanceof Error ? err.message : "Failed to load templates");
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
    try {
      const result = await scaffoldService({
        template_id: templateId,
        name,
        owner,
        description,
      });
      onSuccess(result.service.workspace_path);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create service");
    } finally {
      setSaving(false);
    }
  }

  if (loadingTemplates) {
    return (
      <p className="py-6 text-center text-[13px] text-muted-foreground">
        Loading templates…
      </p>
    );
  }

  if (templates.length === 0) {
    return (
      <div className="space-y-4 py-2">
        <p className="text-[13px] text-muted-foreground">
          No templates configured. You can register an existing service without
          generating a workspace.
        </p>
        {error ? <p className="text-[13px] text-destructive">{error}</p> : null}
        {onRegisterExisting ? (
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 text-[13px]"
            onClick={onRegisterExisting}
          >
            Register existing
          </Button>
        ) : null}
      </div>
    );
  }

  return (
    <form onSubmit={onSubmit} className="grid gap-3 sm:grid-cols-2">
      {error ? (
        <p className="text-[13px] text-destructive sm:col-span-2">{error}</p>
      ) : null}

      <div className="space-y-1.5 sm:col-span-2">
        <Label htmlFor="create-template" className="text-[12px] text-muted-foreground">
          Template
        </Label>
        <select
          id="create-template"
          value={templateId}
          onChange={(e) => setTemplateId(e.target.value)}
          required
          className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 text-[13px] shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
        >
          {templates.map((t) => (
            <option key={t.id} value={t.id}>
              {t.name} ({t.id})
            </option>
          ))}
        </select>
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="create-name" className="text-[12px] text-muted-foreground">
          Name
        </Label>
        <Input
          id="create-name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          placeholder="payments-api"
          className="h-9 text-[13px]"
        />
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="create-owner" className="text-[12px] text-muted-foreground">
          Owner
        </Label>
        <Input
          id="create-owner"
          value={owner}
          onChange={(e) => setOwner(e.target.value)}
          placeholder="platform-team"
          className="h-9 text-[13px]"
        />
      </div>

      <div className="space-y-1.5 sm:col-span-2">
        <Label
          htmlFor="create-description"
          className="text-[12px] text-muted-foreground"
        >
          Description
        </Label>
        <Input
          id="create-description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Optional — defaults from template"
          className="h-9 text-[13px]"
        />
      </div>

      <div className="flex flex-wrap items-center gap-3 pt-1 sm:col-span-2">
        <Button type="submit" size="sm" className="h-8 text-[13px]" disabled={saving}>
          {saving ? "Creating…" : "Create service"}
        </Button>
        {onRegisterExisting ? (
          <button
            type="button"
            className="text-[12px] text-muted-foreground transition-colors hover:text-foreground"
            onClick={onRegisterExisting}
            disabled={saving}
          >
            Register existing instead
          </button>
        ) : null}
      </div>
    </form>
  );
}
