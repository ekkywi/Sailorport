export type Deployment = {
    id: string;
    service_id: string;
    environment_id: string;
    environment_slug: string;
    worker_id: string;
    status: string;
    image_tag: string;
    container_id: string;
    port: number | null;
    error_message: string;
    created_at: string;
    updated_at: string;
}