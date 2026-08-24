export type Deployment = {
    id: string;
    service_id: string;
    environment_id: string;
    environment_slug: string;
    target_worker_id: string | null;
    worker_id: string | null;
    status: string;
    image_tag: string;
    git_sha: string;
    container_id: string;
    port: number | null;
    error_message: string;
    created_at: string;
    updated_at: string;
};