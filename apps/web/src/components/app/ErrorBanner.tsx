import { AlertCircle } from "lucide-react";
import {
  Alert,
  AlertAction,
  AlertDescription,
} from "@/components/ui/alert";
import { Button } from "@/components/ui/button";

type ErrorBannerProps = {
  message: string;
  onRetry?: () => void;
  retryLabel?: string;
};

export function ErrorBanner({
  message,
  onRetry,
  retryLabel = "Retry",
}: ErrorBannerProps) {
  return (
    <Alert variant="destructive" className="pr-20">
      <AlertCircle />
      <AlertDescription className="text-[13px]">{message}</AlertDescription>
      {onRetry ? (
        <AlertAction>
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-7 text-[12px]"
            onClick={onRetry}
          >
            {retryLabel}
          </Button>
        </AlertAction>
      ) : null}
    </Alert>
  );
}
