export type Worker = {
  id: string;
  name: string;
  hostname: string;
  status: "online" | "offline" | "draining" | string;
  labels: Record<string, unknown>;
  last_seen_at: string | null;
  created_at: string;
  updated_at: string;
};
