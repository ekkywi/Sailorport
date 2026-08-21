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
  auto_deploy_enabled: boolean;
  auto_deploy_environment: string;
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
};

export type ServiceFormValues = {
  name: string;
  description: string;
  owner: string;
};

export type GitServiceFormValues = {
  name: string;
  description: string;
  owner: string;
  repo_url: string;
  branch: string;
  dockerfile_path: string;
};

/** True if the service can be deployed (scaffold workspace or linked Git repo). */
export function canDeployService(svc: Pick<Service, "workspace_path" | "source_type" | "repo_url">): boolean {
  if (svc.workspace_path?.trim()) return true;
  return svc.source_type === "git" && Boolean(svc.repo_url?.trim());
}
