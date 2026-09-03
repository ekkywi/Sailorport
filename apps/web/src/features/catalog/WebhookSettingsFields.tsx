import { useState } from "react";
import { Check, Copy, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { Environment } from "../environments/types";
import {
  catalogOptionClassName,
  catalogSelectClassName,
} from "./selectClassName";
import {
  generateWebhookSecret,
  type ServiceFormValues,
} from "./types";

type WebhookSettingsFieldsProps = {
  values: ServiceFormValues;
  environments: Environment[];
  disabled?: boolean;
  onChange: (field: keyof ServiceFormValues, value: string | boolean) => void;
};

const FALLBACK_ENVS = ["dev", "staging", "prod"] as const;

export function WebhookSettingsFields({
  values,
  environments,
  disabled,
  onChange,
}: WebhookSettingsFieldsProps) {
  const [copied, setCopied] = useState(false);

  const envOptions =
    environments.length > 0
      ? environments.map((e) => e.slug)
      : [...FALLBACK_ENVS];

  function handleGenerate() {
    const secret = generateWebhookSecret();
    onChange("webhook_secret", secret);
    setCopied(false);
  }

  async function handleCopy() {
    const text = values.webhook_secret.trim();
    if (!text) return;
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 2000);
    } catch {
      setCopied(false);
    }
  }

  const secretStatus = values.webhook_secret.trim()
    ? "New secret will be saved (copy it now — it won’t be shown again)."
    : values.webhook_secret_set
      ? "A secret is already set. Leave blank to keep it, or generate a new one."
      : "No secret yet. Generate one before enabling auto-deploy.";

  return (
    <div className="space-y-3 border-t border-border/60 pt-3 sm:col-span-2">
      <div>
        <p className="text-[13px] font-medium text-foreground">
          GitHub webhook
        </p>
        <p className="mt-0.5 text-[12px] text-muted-foreground">
          Push to the tracked branch can create a deployment when auto-deploy is
          on. URL:{" "}
          <code className="text-[11px]">POST /api/v1/webhooks/github</code>
        </p>
      </div>

      <div className="flex items-start gap-2.5">
        <input
          id="svc-auto-deploy"
          type="checkbox"
          className="mt-1 size-4 rounded border-border accent-primary"
          checked={values.auto_deploy_enabled}
          disabled={disabled}
          onChange={(e) => onChange("auto_deploy_enabled", e.target.checked)}
        />
        <div className="min-w-0">
          <Label
            htmlFor="svc-auto-deploy"
            className="text-[13px] font-medium text-foreground"
          >
            Auto-deploy on push
          </Label>
          <p className="text-[11px] text-muted-foreground">
            Requires a webhook secret and matching branch on the service.
          </p>
        </div>
      </div>

      <div className="space-y-1.5">
        <Label
          htmlFor="svc-auto-deploy-env"
          className="text-[12px] text-muted-foreground"
        >
          Auto-deploy environment
        </Label>
        <select
          id="svc-auto-deploy-env"
          value={values.auto_deploy_environment || "staging"}
          disabled={disabled}
          onChange={(e) =>
            onChange("auto_deploy_environment", e.target.value)
          }
          className={`${catalogSelectClassName} disabled:cursor-not-allowed disabled:opacity-50`}
        >
          {envOptions.map((slug) => (
            <option key={slug} value={slug} className={catalogOptionClassName}>
              {slug}
            </option>
          ))}
        </select>
      </div>

      <div className="space-y-1.5">
        <Label
          htmlFor="svc-webhook-secret"
          className="text-[12px] text-muted-foreground"
        >
          Webhook secret
        </Label>
        <div className="flex flex-col gap-2 sm:flex-row">
          <Input
            id="svc-webhook-secret"
            type="text"
            value={values.webhook_secret}
            disabled={disabled}
            onChange={(e) => onChange("webhook_secret", e.target.value)}
            placeholder={
              values.webhook_secret_set
                ? "••••••••  (unchanged if empty)"
                : "Generate or paste secret"
            }
            className="h-9 flex-1 font-mono text-[12px]"
            autoComplete="off"
          />
          <div className="flex gap-2">
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="h-9 text-[12px]"
              disabled={disabled}
              onClick={handleGenerate}
            >
              <RefreshCw className="size-3.5" />
              Generate
            </Button>
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="h-9 text-[12px]"
              disabled={disabled || !values.webhook_secret.trim()}
              onClick={() => void handleCopy()}
            >
              {copied ? (
                <Check className="size-3.5" />
              ) : (
                <Copy className="size-3.5" />
              )}
              {copied ? "Copied" : "Copy"}
            </Button>
          </div>
        </div>
        <p className="text-[11px] text-muted-foreground">{secretStatus}</p>
      </div>
    </div>
  );
}
