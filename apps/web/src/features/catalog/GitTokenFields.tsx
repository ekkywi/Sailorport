import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { ServiceFormValues } from "./types";

type GitTokenFieldsProps = {
  values: Pick<ServiceFormValues, "git_token" | "git_token_set">;
  disabled?: boolean;
  onChange: (field: "git_token", value: string) => void;
};

export function GitTokenFields({
  values,
  disabled,
  onChange,
}: GitTokenFieldsProps) {
  const status = values.git_token.trim()
    ? "New token will be saved (it won’t be shown again)."
    : values.git_token_set
      ? "A token is already set. Leave blank to keep it, or paste a new one."
      : "No token yet — public HTTPS clone only. Required for private repos.";

  return (
    <div className="space-y-3 border-t border-border/60 pt-3 sm:col-span-2">
      <div>
        <p className="text-[13px] font-medium text-foreground">
          Git access token
        </p>
        <p className="mt-0.5 text-[12px] text-muted-foreground">
          Optional HTTPS personal access token (read-only) for private
          repositories. Prefer a deploy-scoped token over your personal push
          credentials.
        </p>
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="svc-git-token" className="text-[12px] text-muted-foreground">
          Token
        </Label>
        <Input
          id="svc-git-token"
          type="password"
          autoComplete="new-password"
          value={values.git_token}
          disabled={disabled}
          onChange={(e) => onChange("git_token", e.target.value)}
          placeholder="ghp_… or glpat_…"
          className="h-9 font-mono text-[13px]"
        />
        <p className="text-[11px] text-muted-foreground">{status}</p>
      </div>
    </div>
  );
}
