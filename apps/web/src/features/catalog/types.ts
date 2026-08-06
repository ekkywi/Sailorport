export type Service = {
  id: string;
  name: string;
  description: string;
  owner: string;
  created_at: string;
  updated_at: string;
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
