import { apiFetch, readErrorMessage } from "../../lib/http";
import type { Deployment } from "./types";

export async function createDeployment(
    serviceId: string,
    environment = "dev",
    workerId?: string,
): Promise<Deployment> {
    const body: { environment: string; worker_id?: string } = { environment };
    if (workerId) {
        body.worker_id = workerId;
    }
    const res = await apiFetch(`/api/v1/services/${serviceId}/deployments`, {
        method: "POST",
        body: JSON.stringify(body),
    });
    if (!res.ok) {
        throw new Error (
            await readErrorMessage(res, `Deploy failed: ${res.status}`)
        );
    }
    return res.json();
}

export async function listDeploymentsByService(
    serviceId: string
): Promise<Deployment[]> {
    const res = await apiFetch(`/api/v1/services/${serviceId}/deployments`);
    if (!res.ok) {
        throw new Error (
            await readErrorMessage(res, `Failed to list deployments: ${res.status}`)
        );
    }
    return res.json();
}

export async function redeployDeployment(deploymentId: string): Promise<Deployment> {
    const res = await apiFetch(`/api/v1/deployments/${deploymentId}/redeploy`, {
        method: "POST",
    });
    if (!res.ok) {
        throw new Error(
            await readErrorMessage(res, `Redeploy failed: ${res.status}`)
        );
    }
    return res.json();
}