import { useEffect, useState, type ReactNode } from "react";
import { NavLink, useLocation } from "react-router-dom";
import {
  Boxes,
  LayoutDashboard,
  LogOut,
  Menu,
  Server,
  X,
} from "lucide-react";
import {
  BrandMark,
  SectionLabel,
  userInitials,
} from "@/components/app";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/ThemeToggle";
import type { AuthUser } from "@/features/auth/types";
import { cn } from "@/lib/utils";

type AppShellProps = {
  user: AuthUser;
  onLogout: () => void;
  children: ReactNode;
};

const navItems = [
  { to: "/overview", label: "Overview", icon: LayoutDashboard, end: true },
  { to: "/catalog", label: "Catalog", icon: Boxes },
  { to: "/worker", label: "Workers", icon: Server },
] as const;

const pageMeta: Record<string, { title: string; description: string }> = {
  "/overview": {
    title: "Overview",
    description: "Harbour summary and quick links",
  },
  "/catalog": {
    title: "Catalog",
    description: "Create and manage harbour services",
  },
  "/worker": {
    title: "Workers",
    description: "Registered agent nodes and heartbeats",
  },
};

function navClass({ isActive }: { isActive: boolean }) {
  return cn(
    "group flex items-center gap-2 rounded-md px-2 py-[6px] text-[13px] font-medium tracking-[-0.01em] transition-colors",
    isActive
      ? "bg-sidebar-accent text-sidebar-accent-foreground"
      : "text-sidebar-foreground/70 hover:bg-sidebar-accent/70 hover:text-sidebar-accent-foreground",
  );
}

export function AppShell({ user, onLogout, children }: AppShellProps) {
  const location = useLocation();
  const [mobileOpen, setMobileOpen] = useState(false);
  const meta = pageMeta[location.pathname] ?? {
    title: "Sailorport",
    description: "",
  };
  const displayName = user.name || user.email;

  useEffect(() => {
    setMobileOpen(false);
  }, [location.pathname]);

  const sidebar = (
    <div className="flex h-full flex-col">
      <div className="flex h-12 items-center gap-2 px-3">
        <BrandMark />
        <div className="min-w-0">
          <p className="truncate text-[13px] font-medium tracking-[-0.01em] text-sidebar-foreground">
            Sailorport
          </p>
          <p className="truncate text-[11px] text-muted-foreground">Harbour</p>
        </div>
        <button
          type="button"
          className="ml-auto inline-flex size-7 items-center justify-center rounded-md text-muted-foreground hover:bg-sidebar-accent md:hidden"
          aria-label="Close menu"
          onClick={() => setMobileOpen(false)}
        >
          <X className="size-4" />
        </button>
      </div>

      <div className="flex flex-1 flex-col gap-1 px-2 py-2">
        <SectionLabel className="mb-1">Workspace</SectionLabel>
        <nav className="flex flex-col gap-0.5" aria-label="Main">
          {navItems.map(({ to, label, icon: Icon, ...rest }) => (
            <NavLink
              key={to}
              to={to}
              end={"end" in rest ? rest.end : false}
              className={navClass}
            >
              {({ isActive }) => (
                <>
                  <Icon
                    className={cn(
                      "size-4 shrink-0",
                      isActive
                        ? "text-sidebar-accent-foreground"
                        : "text-muted-foreground group-hover:text-sidebar-accent-foreground",
                    )}
                  />
                  {label}
                </>
              )}
            </NavLink>
          ))}
        </nav>
      </div>

      <div className="mt-auto border-t border-sidebar-border px-3 py-3">
        <div className="flex items-center gap-2.5">
          <span className="flex size-7 shrink-0 items-center justify-center rounded-full bg-sidebar-accent text-[11px] font-semibold tracking-tight text-sidebar-accent-foreground">
            {userInitials(displayName)}
          </span>
          <div className="min-w-0 flex-1">
            <p className="truncate text-[13px] font-medium text-sidebar-foreground">
              {displayName}
            </p>
            <span className="mt-0.5 inline-flex rounded bg-muted px-1.5 py-px text-[10px] font-medium capitalize text-muted-foreground">
              {user.role}
            </span>
          </div>
        </div>
      </div>
    </div>
  );

  return (
    <div className="flex min-h-svh bg-background text-foreground">
      <aside className="sticky top-0 hidden h-svh w-[220px] shrink-0 border-r border-sidebar-border bg-sidebar md:block">
        {sidebar}
      </aside>

      {mobileOpen ? (
        <div className="fixed inset-0 z-40 md:hidden">
          <button
            type="button"
            className="absolute inset-0 bg-foreground/20 backdrop-blur-[2px]"
            aria-label="Close menu overlay"
            onClick={() => setMobileOpen(false)}
          />
          <aside className="absolute inset-y-0 left-0 w-[240px] border-r border-sidebar-border bg-sidebar shadow-lg">
            {sidebar}
          </aside>
        </div>
      ) : null}

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="sticky top-0 z-30 flex h-12 items-center gap-3 border-b border-border bg-background/85 px-4 backdrop-blur-md sm:px-6">
          <button
            type="button"
            className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground hover:bg-muted hover:text-foreground md:hidden"
            aria-label="Open menu"
            onClick={() => setMobileOpen(true)}
          >
            <Menu className="size-4" />
          </button>

          <div className="min-w-0 flex-1">
            <h1 className="truncate text-[13px] font-semibold tracking-[-0.02em]">
              {meta.title}
            </h1>
            {meta.description ? (
              <p className="hidden truncate text-[11px] text-muted-foreground sm:block">
                {meta.description}
              </p>
            ) : null}
          </div>

          <div className="flex items-center gap-1">
            <ThemeToggle />
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="h-8 gap-1.5 px-2 text-[13px] text-muted-foreground hover:text-foreground"
              onClick={onLogout}
            >
              <LogOut className="size-3.5" />
              <span className="hidden sm:inline">Log out</span>
            </Button>
          </div>
        </header>

        <main className="app-harbour flex-1 px-4 py-5 sm:px-6 sm:py-6">
          <div className="mx-auto w-full max-w-5xl">{children}</div>
        </main>
      </div>
    </div>
  );
}
