import type { FormEvent } from "react";
import { useState } from "react";
import { ChevronDown } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import type { GitServiceFormValues } from "./types";

type GitServiceFormProps = {
  values: GitServiceFormValues;
  saving: boolean;
  error?: string;
  onChange: (field: keyof GitServiceFormValues, value: string) => void;
  onSubmit: (e: FormEvent) => void;
  onCancel: () => void;
  onBack?: () => void;
};

export function GitServiceForm({
  values,
  saving,
  error,
  onChange,
  onSubmit,
  onCancel,
  onBack,
}: GitServiceFormProps) {
  const [advancedOpen, setAdvancedOpen] = useState(false);

  return (
    <form onSubmit={onSubmit} className="grid gap-3 sm:grid-cols-2">
      {error ? (
        <p className="text-[13px] text-destructive sm:col-span-2">{error}</p>
      ) : null}

      <div className="space-y-1.5 sm:col-span-2">
        <Label htmlFor="git-repo" className="text-[12px] text-muted-foreground">
          Repository URL
        </Label>
        <Input
          id="git-repo"
          type="url"
          value={values.repo_url}
          onChange={(e) => onChange("repo_url", e.target.value)}
          required
          placeholder="https://github.com/org/app.git"
          className="h-9 font-mono text-[13px]"
        />
        <p className="text-[11px] text-muted-foreground">
          Public repo with a Dockerfile. Private repos are not supported yet.
        </p>
      </div>

      <div className="space-y-1.5 sm:col-span-2">
        <Label htmlFor="git-name" className="text-[12px] text-muted-foreground">
          Service name
        </Label>
        <Input
          id="git-name"
          value={values.name}
          onChange={(e) => onChange("name", e.target.value)}
          required
          placeholder="payments-api"
          className="h-9 text-[13px]"
        />
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="git-owner" className="text-[12px] text-muted-foreground">
          Owner
        </Label>
        <Input
          id="git-owner"
          value={values.owner}
          onChange={(e) => onChange("owner", e.target.value)}
          placeholder="platform-team"
          className="h-9 text-[13px]"
        />
      </div>

      <div className="space-y-1.5">
        <Label
          htmlFor="git-description"
          className="text-[12px] text-muted-foreground"
        >
          Description
        </Label>
        <Input
          id="git-description"
          value={values.description}
          onChange={(e) => onChange("description", e.target.value)}
          placeholder="Optional"
          className="h-9 text-[13px]"
        />
      </div>

      <div className="sm:col-span-2">
        <button
          type="button"
          className="flex items-center gap-1 text-[12px] text-muted-foreground transition-colors hover:text-foreground"
          onClick={() => setAdvancedOpen((o) => !o)}
        >
          <ChevronDown
            className={cn(
              "size-3.5 transition-transform",
              advancedOpen && "rotate-180",
            )}
          />
          Advanced (branch &amp; Dockerfile)
        </button>
        {advancedOpen ? (
          <div className="mt-2 grid gap-3 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label
                htmlFor="git-branch"
                className="text-[12px] text-muted-foreground"
              >
                Branch
              </Label>
              <Input
                id="git-branch"
                value={values.branch}
                onChange={(e) => onChange("branch", e.target.value)}
                placeholder="main"
                className="h-9 font-mono text-[13px]"
              />
            </div>
            <div className="space-y-1.5">
              <Label
                htmlFor="git-dockerfile"
                className="text-[12px] text-muted-foreground"
              >
                Dockerfile path
              </Label>
              <Input
                id="git-dockerfile"
                value={values.dockerfile_path}
                onChange={(e) => onChange("dockerfile_path", e.target.value)}
                placeholder="Dockerfile"
                className="h-9 font-mono text-[13px]"
              />
            </div>
          </div>
        ) : null}
      </div>

      <div className="flex flex-wrap items-center gap-3 pt-1 sm:col-span-2">
        <Button type="submit" size="sm" className="h-8 text-[13px]" disabled={saving}>
          {saving ? "Adding…" : "Add service"}
        </Button>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="h-8 text-[13px]"
          onClick={onCancel}
          disabled={saving}
        >
          Cancel
        </Button>
        {onBack ? (
          <button
            type="button"
            className="text-[12px] text-muted-foreground transition-colors hover:text-foreground"
            onClick={onBack}
            disabled={saving}
          >
            Back
          </button>
        ) : null}
      </div>
    </form>
  );
}
