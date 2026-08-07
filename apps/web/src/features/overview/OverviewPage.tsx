import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { ArrowUpRight, Boxes, Server } from "lucide-react";
import {
  DataPanel,
  EmptyState,
  StatusDot,
  skeletonClass,
} from "@/components/app";
import { listServices } from "../catalog/api";
import type { Service } from "../catalog/types";
import { listWorkers } from "../workers/api";
import type { Worker } from "../workers/types";

type OverviewData = {
  services: Service[];
  workers: Worker[];
};

function MetricSkeleton() {
  return <div className={skeletonClass("mt-2 h-7 w-16")} />;
}

function RowSkeleton() {
  return (
    <div className="space-y-2 px-4 py-3">
      <div className={skeletonClass("h-3.5 w-1/3")} />
      <div className={skeletonClass("h-3 w-1/2")} />
    </div>
  );
}

export function OverviewPage() {
  const [data, setData] = useState<OverviewData | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      setLoading(true);
      setError("");
      try {
        const [services, workers] = await Promise.all([
          listServices(),
          listWorkers(),
        ]);
        if (cancelled) return;
        setData({ services, workers });
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "Failed to load overview");
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    void load();
    return () => {
      cancelled = true;
    };
  }, []);

  const services = data?.services ?? [];
  const workers = data?.workers ?? [];
  const online = workers.filter((w) => w.status === "online").length;
  const offline = workers.filter((w) => w.status !== "online").length;
  const recentServices = services.slice(0, 5);
  const recentWorkers = workers.slice(0, 5);

  return (
    <div className="space-y-6">
      {error ? (
        <p className="text-[13px] text-destructive">{error}</p>
      ) : null}

      <div className="grid gap-3 sm:grid-cols-3">
        <Link
          to="/catalog"
          className="group rounded-lg border border-border bg-card/80 p-4 transition-colors hover:border-foreground/15 hover:bg-accent/30"
        >
          <div className="flex items-center justify-between">
            <span className="text-[12px] text-muted-foreground">Services</span>
            <ArrowUpRight className="size-3.5 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100" />
          </div>
          {loading ? (
            <MetricSkeleton />
          ) : (
            <p className="mt-1.5 text-[22px] font-semibold tracking-[-0.03em] tabular-nums">
              {services.length}
            </p>
          )}
        </Link>

        <Link
          to="/worker"
          className="group rounded-lg border border-border bg-card/80 p-4 transition-colors hover:border-foreground/15 hover:bg-accent/30"
        >
          <div className="flex items-center justify-between">
            <span className="text-[12px] text-muted-foreground">Workers online</span>
            <ArrowUpRight className="size-3.5 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100" />
          </div>
          {loading ? (
            <MetricSkeleton />
          ) : (
            <p className="mt-1.5 text-[22px] font-semibold tracking-[-0.03em] tabular-nums">
              {online}
              <span className="text-[13px] font-medium text-muted-foreground">
                {" "}
                / {workers.length}
              </span>
            </p>
          )}
        </Link>

        <div className="rounded-lg border border-border bg-card/80 p-4">
          <span className="text-[12px] text-muted-foreground">Offline / other</span>
          {loading ? (
            <MetricSkeleton />
          ) : (
            <p className="mt-1.5 text-[22px] font-semibold tracking-[-0.03em] tabular-nums">
              {offline}
            </p>
          )}
        </div>
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        <DataPanel>
          <div className="flex items-center justify-between border-b border-border px-4 py-2.5">
            <p className="text-[13px] font-medium tracking-[-0.01em]">Recent services</p>
            <Link
              to="/catalog"
              className="text-[12px] text-muted-foreground transition-colors hover:text-foreground"
            >
              View all
            </Link>
          </div>
          {loading ? (
            <>
              <RowSkeleton />
              <RowSkeleton />
              <RowSkeleton />
            </>
          ) : recentServices.length === 0 ? (
            <EmptyState
              icon={Boxes}
              title="No services yet"
              description="Create a service from a template in Catalog."
              action={
                <Link
                  to="/catalog"
                  className="text-[12px] font-medium text-primary hover:underline"
                >
                  Open Catalog
                </Link>
              }
              className="py-10"
            />
          ) : (
            <ul className="divide-y divide-border">
              {recentServices.map((svc) => (
                <li key={svc.id} className="px-4 py-2.5">
                  <p className="truncate text-[13px] font-medium tracking-[-0.01em]">
                    {svc.name}
                  </p>
                  <p className="truncate text-[12px] text-muted-foreground">
                    {svc.owner || "Unassigned"}
                    {svc.description ? ` · ${svc.description}` : ""}
                  </p>
                </li>
              ))}
            </ul>
          )}
        </DataPanel>

        <DataPanel>
          <div className="flex items-center justify-between border-b border-border px-4 py-2.5">
            <p className="text-[13px] font-medium tracking-[-0.01em]">Workers</p>
            <Link
              to="/worker"
              className="text-[12px] text-muted-foreground transition-colors hover:text-foreground"
            >
              View all
            </Link>
          </div>
          {loading ? (
            <>
              <RowSkeleton />
              <RowSkeleton />
              <RowSkeleton />
            </>
          ) : recentWorkers.length === 0 ? (
            <EmptyState
              icon={Server}
              title="No workers yet"
              description="Register an agent to see nodes appear here."
              action={
                <Link
                  to="/worker"
                  className="text-[12px] font-medium text-primary hover:underline"
                >
                  Open Workers
                </Link>
              }
              className="py-10"
            />
          ) : (
            <ul className="divide-y divide-border">
              {recentWorkers.map((w) => (
                <li
                  key={w.id}
                  className="flex items-center justify-between gap-3 px-4 py-2.5"
                >
                  <div className="min-w-0">
                    <p className="truncate text-[13px] font-medium tracking-[-0.01em]">
                      {w.name}
                    </p>
                    <p className="truncate text-[12px] text-muted-foreground">
                      {w.hostname || "—"}
                    </p>
                  </div>
                  <span className="inline-flex shrink-0 items-center gap-1.5 text-[11px] capitalize text-muted-foreground">
                    <StatusDot status={w.status} />
                    {w.status}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </DataPanel>
      </div>
    </div>
  );
}
