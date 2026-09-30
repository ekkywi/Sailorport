import {
  buildDeployUrl,
  defaultDeployHost,
  deployHttpLinksEnabled,
} from "@/features/catalog/deployUrls";

type DeployedPortInfoProps = {
  hostPort: number;
  containerPort: number;
  sourceType: string;
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
  sourceType,
}: DeployedPortInfoProps) {
  const mapping = `${hostPort} → ${containerPort}`;
  const title = `External (host) ${hostPort} → internal (container) ${containerPort}`;
  const showLinks = deployHttpLinksEnabled(sourceType);

  const mappingEl = (
    <span className="font-mono text-[11px] text-muted-foreground" title={title}>
      {mapping}
    </span>
  );

  if (!showLinks) {
    return mappingEl;
  }

  const host = defaultDeployHost();
  const openHref = buildDeployUrl({ host, hostPort, kind: "open" });
  const healthHref = buildDeployUrl({ host, hostPort, kind: "healthz" });

  const linkClass =
    "text-[11px] text-muted-foreground underline-offset-2 hover:text-foreground hover:underline";

  return (
    <span className="inline-flex min-w-0 flex-wrap items-center gap-x-2 gap-y-0.5">
      {mappingEl}
      <a
        href={openHref}
        target="_blank"
        rel="noreferrer"
        className={linkClass}
        title={`Open app — ${openHref}`}
        onClick={(e) => e.stopPropagation()}
      >
        Open
      </a>
      <a
        href={healthHref}
        target="_blank"
        rel="noreferrer"
        className={linkClass}
        title={`Health check — ${healthHref}`}
        onClick={(e) => e.stopPropagation()}
      >
        Health
      </a>
    </span>
  );
}