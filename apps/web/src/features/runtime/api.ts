import { apiFetch, readErrorMessage } from "../../lib/http";

export type RuntimeJob = {
  id: string;
  service_id: string;
  deployment_id: string;
  service_name: string;
  environment_slug: string;
  action: string;
  status: string;
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