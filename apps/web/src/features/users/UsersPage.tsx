import { useCallback, useEffect, useState, type FormEvent } from "react";
import { Plus, RefreshCw, Users as UsersIcon } from "lucide-react";
import {
  DataPanel,
  EmptyState,
  ErrorBanner,
  StatusBadge,
  Toolbar,
  formatRelativeTime,
  useToast,
} from "@/components/app";
import {
  AlertDialog,
  AlertDialogClose,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { me } from "@/features/auth/api";
import type { AuthUser } from "@/features/auth/types";
import {
  createUser,
  listUsers,
  resetUserPassword,
  setUserDisabled,
  updateUserRole,
} from "./api";
import type { User, UserRole } from "./type";

const ROLES: UserRole[] = ["admin", "developer", "viewer"];

const emptyForm = {
  email: "",
  name: "",
  password: "",
  role: "developer" as UserRole,
};

export function UsersPage() {
  const { toast } = useToast();
  const [currentUser, setCurrentUser] = useState<AuthUser | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [savingId, setSavingId] = useState<string | null>(null);
  const [createOpen, setCreateOpen] = useState(false);
  const [form, setForm] = useState(emptyForm);
  const [creating, setCreating] = useState(false);
  const [formError, setFormError] = useState("");
  const [disableTarget, setDisableTarget] = useState<User | null>(null);
  const [togglingId, setTogglingId] = useState<string | null>(null);
  const [passwordTarget, setPasswordTarget] = useState<User | null>(null);
  const [newPassword, setNewPassword] = useState("");
  const [resetting, setResetting] = useState(false);
  const [passwordError, setPasswordError] = useState("");

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

  function closeCreate() {
    if (creating) return;
    setCreateOpen(false);
    setForm(emptyForm);
    setFormError("");
  }

  async function onCreateSubmit(e: FormEvent) {
    e.preventDefault();
    setCreating(true);
    setFormError("");
    try {
      const created = await createUser({
        email: form.email.trim(),
        name: form.name.trim(),
        password: form.password,
        role: form.role,
      });
      toast(`Created ${created.email} — share the temporary password securely`);
      setCreateOpen(false);
      setForm(emptyForm);
      await load();
    } catch (err) {
      setFormError(err instanceof Error ? err.message : "Failed to create user");
    } finally {
      setCreating(false);
    }
  }

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

  async function confirmDisable() {
    if (!disableTarget) return;
    const target = disableTarget;
    setTogglingId(target.id);
    setError("");
    try {
      const updated = await setUserDisabled(target.id, true);
      setUsers((prev) => prev.map((u) => (u.id === updated.id ? updated : u)));
      toast(`Disabled ${updated.email}`);
      setDisableTarget(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to disable user");
      setDisableTarget(null);
    } finally {
      setTogglingId(null);
    }
  }

  async function onEnable(user: User) {
    setTogglingId(user.id);
    setError("");
    try {
      const updated = await setUserDisabled(user.id, false);
      setUsers((prev) => prev.map((u) => (u.id === updated.id ? updated : u)));
      toast(`Enabled ${updated.email}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to enable user");
    } finally {
      setTogglingId(null);
    }
  }

  function closePasswordDialog() {
    if (resetting) return;
    setPasswordTarget(null);
    setNewPassword("");
    setPasswordError("");
  }

  async function onResetPasswordSubmit(e: FormEvent) {
    e.preventDefault();
    if (!passwordTarget) return;
    setResetting(true);
    setPasswordError("");
    try {
      await resetUserPassword(passwordTarget.id, newPassword);
      toast(`Password reset for ${passwordTarget.email} — share it securely`);
      setPasswordTarget(null);
      setNewPassword("");
      setPasswordError("");
    } catch (err) {
      setPasswordError(
        err instanceof Error ? err.message : "Failed to reset password",
      );
    } finally {
      setResetting(false);
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
          <>
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
            <Button
              type="button"
              size="sm"
              className="h-8 gap-1.5 text-[13px]"
              onClick={() => {
                setFormError("");
                setCreateOpen(true);
              }}
            >
              <Plus className="size-3.5" />
              Create user
            </Button>
          </>
        }
      />
      {error ? <ErrorBanner message={error} onRetry={() => void load()} /> : null}
      {!loading && users.length === 0 && !error ? (
        <DataPanel>
          <EmptyState
            icon={UsersIcon}
            title="No users"
            description="Create a user here, or they can self-register at /register."
            action={
              <Button
                type="button"
                size="sm"
                className="h-8 text-[13px]"
                onClick={() => setCreateOpen(true)}
              >
                Create user
              </Button>
            }
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
                  <th className="px-4 py-2.5">Status</th>
                  <th className="px-4 py-2.5">Joined</th>
                  <th className="px-4 py-2.5">
                    <span className="sr-only">Actions</span>
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {users.map((u) => {
                  const isSelf = currentUser?.id === u.id;
                  return (
                    <tr
                      key={u.id}
                      className={cn(
                        "hover:bg-muted/35",
                        u.disabled && "opacity-60",
                      )}
                    >
                      <td className="px-4 py-3">
                        <p className="font-medium">{u.name || u.email}</p>
                        <p className="text-[12px] text-muted-foreground">
                          {u.email}
                        </p>
                      </td>
                      <td className="px-4 py-3">
                        {isSelf ? (
                          <StatusBadge status={u.role} />
                        ) : (
                          <select
                            className="h-8 rounded-md border border-border bg-background px-2 text-[12px] capitalize"
                            value={u.role}
                            disabled={savingId === u.id || u.disabled}
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
                      <td className="px-4 py-3">
                        <StatusBadge
                          status={u.disabled ? "disabled" : "active"}
                        />
                      </td>
                      <td className="px-4 py-3 text-muted-foreground">
                        {formatRelativeTime(u.created_at)}
                      </td>
                      <td className="px-4 py-3">
                        {isSelf ? (
                          <span className="text-[12px] text-muted-foreground">
                            You
                          </span>
                        ) : (
                          <div className="flex flex-wrap items-center justify-end gap-1.5">
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              className="h-8 text-[13px]"
                              disabled={resetting}
                              onClick={() => {
                                setPasswordError("");
                                setNewPassword("");
                                setPasswordTarget(u);
                              }}
                            >
                              Reset password
                            </Button>
                            {u.disabled ? (
                              <Button
                                type="button"
                                variant="outline"
                                size="sm"
                                className="h-8 text-[13px]"
                                disabled={togglingId === u.id}
                                onClick={() => void onEnable(u)}
                              >
                                {togglingId === u.id ? "Enabling…" : "Enable"}
                              </Button>
                            ) : (
                              <Button
                                type="button"
                                variant="outline"
                                size="sm"
                                className="h-8 text-[13px]"
                                disabled={togglingId === u.id}
                                onClick={() => setDisableTarget(u)}
                              >
                                Disable
                              </Button>
                            )}
                          </div>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </DataPanel>
      ) : null}

      <Dialog
        open={createOpen}
        onOpenChange={(open) => {
          if (!open) closeCreate();
          else setCreateOpen(true);
        }}
      >
        <DialogContent className="sm:max-w-md" showCloseButton={!creating}>
          <DialogHeader>
            <DialogTitle>Create user</DialogTitle>
            <DialogDescription>
              Create an account and share the temporary password with the
              teammate. They can sign in at /login.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={onCreateSubmit} className="space-y-4">
            {formError ? (
              <p className="text-[13px] text-destructive">{formError}</p>
            ) : null}
            <div className="space-y-1.5">
              <Label htmlFor="create-email" className="text-[12px] text-muted-foreground">
                Email
              </Label>
              <Input
                id="create-email"
                type="email"
                required
                autoFocus
                autoComplete="off"
                disabled={creating}
                value={form.email}
                onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
                className="h-8 text-[13px]"
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="create-name" className="text-[12px] text-muted-foreground">
                Name
              </Label>
              <Input
                id="create-name"
                autoComplete="off"
                disabled={creating}
                placeholder="Optional"
                value={form.name}
                onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
                className="h-8 text-[13px]"
              />
            </div>
            <div className="space-y-1.5">
              <Label
                htmlFor="create-password"
                className="text-[12px] text-muted-foreground"
              >
                Temporary password
              </Label>
              <Input
                id="create-password"
                type="password"
                required
                minLength={8}
                autoComplete="new-password"
                disabled={creating}
                value={form.password}
                onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))}
                className="h-8 text-[13px]"
              />
              <p className="text-[11px] text-muted-foreground">At least 8 characters.</p>
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="create-role" className="text-[12px] text-muted-foreground">
                Role
              </Label>
              <select
                id="create-role"
                className="h-8 w-full rounded-md border border-border bg-background px-2 text-[13px] capitalize"
                disabled={creating}
                value={form.role}
                onChange={(e) =>
                  setForm((f) => ({ ...f, role: e.target.value as UserRole }))
                }
              >
                {ROLES.map((r) => (
                  <option key={r} value={r}>
                    {r}
                  </option>
                ))}
              </select>
            </div>
            <div className="flex justify-end gap-2 pt-1">
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="h-8 text-[13px]"
                disabled={creating}
                onClick={closeCreate}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                size="sm"
                className="h-8 text-[13px]"
                disabled={creating}
              >
                {creating ? "Creating…" : "Create user"}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <Dialog
        open={passwordTarget !== null}
        onOpenChange={(open) => {
          if (!open) closePasswordDialog();
        }}
      >
        <DialogContent className="sm:max-w-md" showCloseButton={!resetting}>
          <DialogHeader>
            <DialogTitle>Reset password</DialogTitle>
            <DialogDescription>
              Set a temporary password for{" "}
              <span className="font-medium text-foreground">
                {passwordTarget?.email}
              </span>
              . Share it securely — they sign in at /login.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={onResetPasswordSubmit} className="space-y-4">
            {passwordError ? (
              <p className="text-[13px] text-destructive">{passwordError}</p>
            ) : null}
            <div className="space-y-1.5">
              <Label
                htmlFor="reset-password"
                className="text-[12px] text-muted-foreground"
              >
                New temporary password
              </Label>
              <Input
                id="reset-password"
                type="password"
                required
                minLength={8}
                autoFocus
                autoComplete="new-password"
                disabled={resetting}
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                className="h-8 text-[13px]"
              />
              <p className="text-[11px] text-muted-foreground">
                At least 8 characters.
              </p>
            </div>
            <div className="flex justify-end gap-2 pt-1">
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="h-8 text-[13px]"
                disabled={resetting}
                onClick={closePasswordDialog}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                size="sm"
                className="h-8 text-[13px]"
                disabled={resetting}
              >
                {resetting ? "Saving…" : "Reset password"}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <AlertDialog
        open={disableTarget !== null}
        onOpenChange={(open) => {
          if (!open && togglingId === null) setDisableTarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Disable user?</AlertDialogTitle>
            <AlertDialogDescription>
              <span className="font-medium text-foreground">
                {disableTarget?.email}
              </span>{" "}
              will not be able to sign in until enabled again.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogClose
              render={
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="h-8 text-[13px]"
                  disabled={togglingId !== null}
                />
              }
            >
              Cancel
            </AlertDialogClose>
            <Button
              type="button"
              variant="destructive"
              size="sm"
              className="h-8 text-[13px]"
              disabled={togglingId !== null}
              onClick={() => void confirmDisable()}
            >
              {togglingId ? "Disabling…" : "Disable"}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
