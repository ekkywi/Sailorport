import { useCallback, useEffect, useRef, useState } from "react";
import { Loader2, RefreshCw } from "lucide-react";
import { ErrorBanner } from "@/components/app";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import { getRuntimeJob, requestLogs } from "../runtime/api";

type LogsDialogProps = {
  serviceId: string | null;
  serviceName: string;
  environment: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function LogsDialog({
  serviceId,
  serviceName,
  environment,
  open,
  onOpenChange,
}: LogsDialogProps) {
  const [output, setOutput] = useState("");
  const [status, setStatus] = useState<"idle" | "loading" | "done" | "failed">(
    "idle",
  );
  const [error, setError] = useState("");
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const stopPolling = useCallback(() => {
    if (pollRef.current) {
      clearInterval(pollRef.current);
      pollRef.current = null;
    }
  }, []);

  const fetchLogs = useCallback(async () => {
    if (!serviceId) return;
    setStatus("loading");
    setError("");
    setOutput("");
    stopPolling();

    try {
      const job = await requestLogs(serviceId, environment);

      pollRef.current = setInterval(async () => {
        try {
          const current = await getRuntimeJob(job.id);
          if (current.status === "done") {
            setOutput(current.output || "(no output)");
            setStatus("done");
            stopPolling();
          } else if (current.status === "failed") {
            setError(current.error_message || "Agent failed to fetch logs");
            setStatus("failed");
            stopPolling();
          }
        } catch {
          setError("Failed to poll job status");
          setStatus("failed");
          stopPolling();
        }
      }, 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to request logs");
      setStatus("failed");
    }
  }, [serviceId, environment, stopPolling]);

  useEffect(() => {
    if (open && serviceId) {
      void fetchLogs();
    }
    if (!open) {
      stopPolling();
      setStatus("idle");
      setOutput("");
      setError("");
    }
    return stopPolling;
  }, [open, serviceId, environment, fetchLogs, stopPolling]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>
            Logs — {serviceName} ({environment})
          </DialogTitle>
          <DialogDescription>
            Container logs from the agent node.
          </DialogDescription>
        </DialogHeader>

        <div className="flex justify-end">
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 gap-1.5 text-[13px]"
            onClick={() => void fetchLogs()}
            disabled={status === "loading"}
          >
            <RefreshCw
              className={cn("size-3.5", status === "loading" && "animate-spin")}
            />
            Refresh
          </Button>
        </div>

        {error ? (
          <ErrorBanner message={error} onRetry={() => void fetchLogs()} />
        ) : null}

        {status === "loading" ? (
          <div className="flex items-center justify-center gap-2 py-10 text-[13px] text-muted-foreground">
            <Loader2 className="size-4 animate-spin" />
            Waiting for agent…
          </div>
        ) : null}

        {status === "done" || output ? (
          <pre className="max-h-[60vh] overflow-auto rounded-lg border border-border bg-muted/40 p-3 font-mono text-[11px] leading-relaxed text-foreground">
            {output}
          </pre>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
