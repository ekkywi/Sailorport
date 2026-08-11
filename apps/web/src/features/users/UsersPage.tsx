import { useCallback, useEffect, useState } from "react";
import { RefreshCw, Users as UsersIcon } from "lucide-react";
import {
  DataPanel,
  EmptyState,
  ErrorBanner,
  StatusBadge,
  Toolbar,
  formatRelativeTime,
  useToast,
} from "@/components/app";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { me } from "@/features/auth/api";
import type { AuthUser } from "@/features/auth/types";
import { listUsers, updateUserRole } from "./api";
import type { User, UserRole } from "./type";

const ROLES: UserRole[] = ["admin", "developer", "viewer"];

export function UsersPage() {
  const { toast } = useToast();
  const [currentUser, setCurrentUser] = useState<AuthUser | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [savingId, setSavingId] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const [meUser, list] = await Promise.all([me(), listUsers()]);
      setCurrentUser(meUser);
      setUsers(list);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load users");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  async function onRoleChange(userId: string, role: UserRole) {
    setSavingId(userId);
    setError("");
    try {
      const updated = await updateUserRole(userId, role);
      setUsers((prev) => prev.map((u) => (u.id === updated.id ? updated : u)));
      toast(`Role updated for ${updated.email}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update role");
    } finally {
      setSavingId(null);
    }
  }

  const meta =
    loading && users.length === 0
      ? "Loading…"
      : `${users.length} user${users.length === 1 ? "" : "s"}`;

  return (
    <div className="space-y-4">
      <Toolbar
        meta={meta}
        actions={
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 gap-1.5 text-[13px]"
            onClick={() => void load()}
            disabled={loading}
          >
            <RefreshCw className={cn("size-3.5", loading && "animate-spin")} />
            Refresh
          </Button>
        }
      />
      {error ? <ErrorBanner message={error} onRetry={() => void load()} /> : null}
      {!loading && users.length === 0 && !error ? (
        <DataPanel>
          <EmptyState
            icon={UsersIcon}
            title="No users"
            description="Register accounts via /register."
            className="py-16"
          />
        </DataPanel>
      ) : null}
      {users.length > 0 ? (
        <DataPanel>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-[13px]">
              <thead>
                <tr className="border-b border-border text-[11px] font-medium uppercase tracking-wide text-muted-foreground">
                  <th className="px-4 py-2.5">User</th>
                  <th className="px-4 py-2.5">Role</th>
                  <th className="px-4 py-2.5">Joined</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {users.map((u) => {
                  const isSelf = currentUser?.id === u.id;
                  return (
                    <tr key={u.id} className="hover:bg-muted/35">
                      <td className="px-4 py-3">
                        <p className="font-medium">{u.name || u.email}</p>
                        <p className="text-[12px] text-muted-foreground">{u.email}</p>
                      </td>
                      <td className="px-4 py-3">
                        {isSelf ? (
                          <StatusBadge status={u.role} />
                        ) : (
                          <select
                            className="h-8 rounded-md border border-border bg-background px-2 text-[12px] capitalize"
                            value={u.role}
                            disabled={savingId === u.id}
                            onChange={(e) =>
                              void onRoleChange(u.id, e.target.value as UserRole)
                            }
                          >
                            {ROLES.map((r) => (
                              <option key={r} value={r}>
                                {r}
                              </option>
                            ))}
                          </select>
                        )}
                      </td>
                      <td className="px-4 py-3 text-muted-foreground">
                        {formatRelativeTime(u.created_at)}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </DataPanel>
      ) : null}
    </div>
  );
}
