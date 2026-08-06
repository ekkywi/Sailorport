export type TemplateManifest = {
  id: string;
  name: string;
  description: string;
  language: string;
};

export type ScaffoldInput = {
  template_id: string;
  name: string;
  owner: string;
  description: string;
};

export type ScaffoldResult = {
  service: {
    id: string;
    name: string;
    description: string;
    owner: string;
    template_id: string;
    workspace_path: string;
    created_at: string;
    updated_at: string;
  };
};
