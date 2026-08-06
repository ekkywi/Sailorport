import type { CreateServiceInput,UpdateServiceInput, Service } from "./types";

export async function listServices(): Promise<Service[]> {
    const res = await fetch("/api/v1/services");
    if (!res.ok) {
        throw new Error(`Gagal list services: ${res.status}`);
    }
    return res.json();
}

export async function createService(
    input: CreateServiceInput
): Promise<Service> {
    const res = await fetch("/api/v1/services", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(input)
    });

    if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Gagal create: ${res.status}`);
    }
    return res.json();
}

export async function updateService(
    id: string,
    input: UpdateServiceInput,
): Promise<Service> {
    const res = await fetch(`/api/v1/services/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(input),
    });
    
    if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Gagal update: ${res.status}`);
    }
    return res.json();
}

export async function deleteService(id: string): Promise<void> {
    const res = await fetch(`/api/v1/services/${id}`, {
        method: "DELETE",
    });
    
    if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Gagal delete: ${res.status}`);
    }
}