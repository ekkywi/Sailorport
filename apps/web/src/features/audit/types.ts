export type AuditEvent = {
  id: string;
  at: string;
  actor_id: string;
  actor_email: string;
  action: string;
  resource_type: string;
  resource_id: string;
  resource_name: string;
  payload: Record<string, unknown>;
};
