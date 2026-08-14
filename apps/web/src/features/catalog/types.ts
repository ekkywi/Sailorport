export type LatestDeployment = {
  id: string;
  service_id: string;
  environment_id?: string;
  environment_slug?: string;
  status: string;
  port: number | null;
  error_message: string;
  created_at: string;
  updated_at: string;
}

export type Service = {
  id: string;
  name: string;
  description: string;
  owner: string;
  template_id: string;
  workspace_path: string;
  created_at: string;
  updated_at: string;
  latest_deployment?: LatestDeployment | null;
};

export type CreateServiceInput = {
  name: string;
  description: string;
  owner: string;
};

export type UpdateServiceInput = {
  name: string;
  description: string;
  owner: string;
};

export type ServiceFormValues = {
  name: string;
  description: string;
  owner: string;
};
