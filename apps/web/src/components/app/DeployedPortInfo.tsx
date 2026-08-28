type DeployedPortInfoProps = {
  hostPort: number;
  containerPort: number;
  /** HTTP apps: link host port to /healthz. Catalog apps: plain text only. */
  linkHealthz?: boolean;
};

export function containerPortForSource(
  sourceType: string,
  serviceContainerPort: number,
): number {
  if (sourceType === "catalog_app" && serviceContainerPort > 0) {
    return serviceContainerPort;
  }
  return 8080;
}

export function DeployedPortInfo({
  hostPort,
  containerPort,
  linkHealthz = false,
}: DeployedPortInfoProps) {
  const mapping = `${hostPort} → ${containerPort}`;
  const title = `External (host) ${hostPort} → internal (container) ${containerPort}`;

  if (linkHealthz) {
    return (
      <a
        href={`http://localhost:${hostPort}/healthz`}
        target="_blank"
        rel="noreferrer"
        className="font-mono text-[11px] text-muted-foreground underline-offset-2 hover:text-foreground hover:underline"
        title={`${title} — open /healthz`}
        onClick={(e) => e.stopPropagation()}
      >
        {mapping}
      </a>
    );
  }

  return (
    <span
      className="font-mono text-[11px] text-muted-foreground"
      title={title}
    >
      {mapping}
    </span>
  );
}
