import { apiFetch, readErrorMessage } from "../../lib/http";

export type RuntimeJob = {
  id: string;
  service_id: string;
  deployment_id: string;
  service_name: string;
  environment_slug: string;
  action: string;
  status: string;
  worker_id: string | null;
  error_message: string;
  output: string;
  created_at: string;
  updated_at: string;
};

export async function stopService(
  serviceId: string,
  environment = "dev",
) : Promise<RuntimeJob> {
  const res = await apiFetch(`/api/v1/services/${serviceId}/runtime/stop`, {
    method: "POST",
    body: JSON.stringify({ environment }),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Stop failed: ${res.status}`));
  }
  return res.json();
}

export async function startService(
  serviceId: string,
  environment = "dev",
) : Promise<RuntimeJob> {
  const res = await apiFetch(`/api/v1/services/${serviceId}/runtime/start`, {
    method: "POST",
    body: JSON.stringify({ environment }),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Start failed: ${res.status}`));
  }
  return res.json();
}

export async function requestLogs(
  serviceId: string,
  environment = "dev",
): Promise<RuntimeJob> {
  const res = await apiFetch(`/api/v1/services/${serviceId}/runtime/logs`, {
    method: "POST",
    body: JSON.stringify({ environment }),
  });
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, `Logs failed: ${res.status}`));
  }
  return res.json();
}

export async function getRuntimeJob(jobId: string): Promise<RuntimeJob> {
  const res = await apiFetch(`/api/v1/runtime/${jobId}`);
  if (!res.ok) {
    throw new Error(
      await readErrorMessage(res, `Get job failed: ${res.status}`),
    );
  }
  return res.json();
}