export type DeployLinkKind = "open" | "healthz";

export type DeployUrlOptions = {
    host: string;
    hostPort: number;
    kind: DeployLinkKind;
}

export function normalizeDeployHost(host: string): string {
    let h = host.trim();
    if (!h) return "localhost";
    h = h.replace(/^https?:\/\//i, "");
    h = h.split("/")[0] ?? h;
    h = h.split(":")[0] ?? h;
    return h || "localhost";
}

export function buildDeployUrl(opts: DeployUrlOptions): string {
    const host = normalizeDeployHost(opts.host);
    const port = opts.hostPort;
    if (!Number.isFinite(port) || port <= 0) {
        throw new Error("invalid host port");
    }
    const path = opts.kind === "healthz" ? "/healthz" : "";
    return `http://${host}:${port}${path}`;
}

export function deployHttpLinksEnabled(sourceType: string): boolean {
    return sourceType === "git" || sourceType === "scaffold";
}

export function canShowDeployLinks(
    status: string,
    hostPort: number | null | undefined,
    sourceType: string,
): boolean {
    if (status !== "running") return false;
    if (hostPort == null || hostPort <= 0) return false;
    return deployHttpLinksEnabled(sourceType);
}

export function defaultDeployHost(
    envHost: string | undefined = import.meta.env.VITE_DEPLOY_HOST as
        | string
        | undefined
): string {
    return normalizeDeployHost(envHost ?? "");
}