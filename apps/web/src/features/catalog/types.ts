export type LatestDeployment = {
  id: string;
  service_id: string;
  environment_id?: string;
  environment_slug?: string;
  target_worker_id?: string | null;
  worker_id?: string | null;
  status: string;
  port: number | null;
  error_message: string;
  created_at: string;
  updated_at: string;
};

export type Service = {
  id: string;
  name: string;
  description: string;
  owner: string;
  template_id: string;
  workspace_path: string;
  source_type: string;
  repo_url: string;
  branch: string;
  dockerfile_path: string;
  webhook_secret: string;
  webhook_secret_set: boolean;
  auto_deploy_enabled: boolean;
  auto_deploy_environment: string;
  catalog_app_id: string;
  image: string;
  container_port: number;
  created_at: string;
  updated_at: string;
  latest_deployment?: LatestDeployment | null;
  env_deployments?: Record<string, LatestDeployment>;
};

export type CreateServiceInput = {
  name: string;
  description: string;
  owner: string;
  source_type?: string;
  repo_url?: string;
  branch?: string;
  dockerfile_path?: string;
  webhook_secret?: string;
  auto_deploy_enabled?: boolean;
  auto_deploy_environment?: string;
  catalog_app_id?: string;
  image?: string;
  container_port?: number;
};

export type UpdateServiceInput = {
  name: string;
  description: string;
  owner: string;
  source_type?: string;
  repo_url?: string;
  branch?: string;
  dockerfile_path?: string;
  webhook_secret?: string;
  auto_deploy_enabled?: boolean;
  auto_deploy_environment?: string;
  catalog_app_id?: string;
  image?: string;
  container_port?: number;
};

export type ServiceFormValues = {
  name: string;
  description: string;
  owner: string;
  webhook_secret: string;
  webhook_secret_set: boolean;
  auto_deploy_enabled: boolean;
  auto_deploy_environment: string;
};

export type GitServiceFormValues = {
  name: string;
  description: string;
  owner: string;
  repo_url: string;
  branch: string;
  dockerfile_path: string;
};

export type CatalogApp = {
  id: string;
  name: string;
  description: string;
  image: string;
  container_port: number;
  tags: string[];
}

export type CatalogAppFormValues = {
  catalog_app_id: string;
  name: string;
  description: string;
  owner: string;
}

export const emptyServiceForm: ServiceFormValues = {
  name: "",
  description: "",
  owner: "",
  webhook_secret: "",
  webhook_secret_set: false,
  auto_deploy_enabled: false,
  auto_deploy_environment: "staging",
};

export function serviceToFormValues(svc: Service): ServiceFormValues {
  return {
    name: svc.name,
    description: svc.description,
    owner: svc.owner,
    webhook_secret: "",
    webhook_secret_set: Boolean(svc.webhook_secret_set),
    auto_deploy_enabled: Boolean(svc.auto_deploy_enabled),
    auto_deploy_environment: svc.auto_deploy_environment || "staging",
  };
}

export function formValuesToUpdateInput(
  values: ServiceFormValues,
): UpdateServiceInput {
  const input: UpdateServiceInput = {
    name: values.name.trim(),
    description: values.description.trim(),
    owner: values.owner.trim(),
    auto_deploy_enabled: values.auto_deploy_enabled,
    auto_deploy_environment: values.auto_deploy_environment.trim() || "staging",
  };
  const secret = values.webhook_secret.trim();
  if (secret) {
    input.webhook_secret = secret;
  }
  return input;
}

export function generateWebhookSecret(bytes = 32): string {
  const arr = new Uint8Array(bytes);
  crypto.getRandomValues(arr);
  return Array.from(arr, (b) => b.toString(16).padStart(2, "0")).join("");
}

export function canDeployService(
  svc: Pick<Service, "workspace_path" | "source_type" | "repo_url" | "image">,
): boolean {
  if (svc.workspace_path?.trim()) return true;
  if (svc.source_type === "git" && Boolean(svc.repo_url?.trim())) return true;
  if (svc.source_type === "catalog_app" && Boolean(svc.image?.trim())) return true;
  return false;
}
