import { useCallback, useEffect, useState } from "react";
import { Loader2, Settings } from "lucide-react";
import {
    DataPanel,
    ErrorBanner,
    Toolbar,
    useToast,
} from "@/components/app";
import { Label } from "@/components/ui/label";
import { getSettings, updateSettings } from "./api";

export function SettingsPage() {
    const { toast } = useToast();
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState("");
    const [registrationOpen, setRegistrationOpen] = useState(false);

    const load = useCallback(async () => {
        setLoading(true);
        setError("");
        try {
            const st = await getSettings();
            setRegistrationOpen(st.registration_open);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load settings");
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        void load();
    }, [load]);

    async function onToggle(next: boolean) {
        setSaving(true);
        setError("");
        try {
            const st = await updateSettings(next);
            setRegistrationOpen(st.registration_open);
            toast(
                st.registration_open
                    ? "Public registration enabled"
                    : "Public registration disabled",
            );
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to update settings");
            try {
                const st = await getSettings();
                setRegistrationOpen(st.registration_open);
            } catch {
                // Ignore
            }
        } finally {
            setSaving(false);
        }
    }

    return (
        <div className="space-y-4">
            <Toolbar
                meta={
                    loading
                        ? "Loading..."
                        : registrationOpen
                            ? "Public registration is open"
                            : "Public registration is closed"
                }
            />

            {error ? (
                <ErrorBanner message={error} onRetry={() => void load()} />
            ) : null}

            <DataPanel>
                {loading ? (
                    <div className="flex items-center gap-2 px-4 py-8 text-sm text-muted-foreground">
                        <Loader2 className="size-4 animate-spin" />
                        Loading settings..
                    </div>
                ) : (
                    <div className="flex items-start gap-4 px-4 py-5">
                        <div className="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground">
                            <Settings className="size-4" />
                        </div>
                        <div className="min-w-0 flex-1 space-y-1">
                            <Label
                                htmlFor="registration-open"
                                className="text-[13px] font-medium text-foreground"
                            >
                                Allow public registration
                            </Label>
                            <p className="text-[13px] leading-5 text-muted-foreground">
                                When enabled, new users can sign up at /register as developers.
                                The first admin is still created only via /setup. Turn this off
                                to require an admin to create accounts.
                            </p>
                        </div>
                        <input
                            id="registration-open"
                            type="checkbox"
                            className="mt-1 size-4 accent-foreground"
                            checked={registrationOpen}
                            disabled={saving}
                            onChange={(e) => void onToggle(e.target.checked)}
                        />
                    </div>
                )}
            </DataPanel>
        </div>
    );
}