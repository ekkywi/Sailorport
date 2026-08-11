import { useEffect, useState, type ReactNode } from "react";
import { NavLink, useLocation } from "react-router-dom";
import {
  Boxes,
  ChevronDown,
  LayoutDashboard,
  LogOut,
  ChevronLeft,
  ChevronRight,
  Menu,
  Server,
  X,
} from "lucide-react";
import { BrandMark, userInitials } from "@/components/app";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ThemeToggle } from "@/components/ThemeToggle";
import type { AuthUser } from "@/features/auth/types";
import { cn } from "@/lib/utils";

type AppShellProps = {
  user: AuthUser;
  onLogout: () => void;
  children: ReactNode;
};

const SIDEBAR_COLLAPSED_KEY = "sailorport.sidebar.collapsed";

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

function readCollapsed(): boolean {
  try {
    return localStorage.getItem(SIDEBAR_COLLAPSED_KEY) === "1";
  } catch {
    return false;
  }
}

function writeCollapsed(collapsed: boolean) {
  try {
    localStorage.setItem(SIDEBAR_COLLAPSED_KEY, collapsed ? "1" : "0");
  } catch {
    // ignore quota / private mode
  }
}

function SidebarNav({
  collapsed,
  onCloseMobile,
}: {
  collapsed: boolean;
  onCloseMobile?: () => void;
}) {
  return (
    <div className="flex h-full flex-col">
      <div
        className={cn(
          "flex h-12 shrink-0 items-center border-b border-sidebar-border",
          collapsed ? "justify-center px-2" : "gap-2 px-3",
        )}
      >
        <BrandMark />
        {!collapsed ? (
          <div className="min-w-0 flex-1 overflow-hidden">
            <p className="truncate text-[13px] font-medium tracking-[-0.01em] text-sidebar-foreground">
              Sailorport
            </p>
            <p className="truncate text-[11px] text-muted-foreground">Harbour</p>
          </div>
        ) : null}

        {onCloseMobile ? (
          <button
            type="button"
            className="ml-auto inline-flex size-7 items-center justify-center rounded-md text-muted-foreground hover:bg-sidebar-accent md:hidden"
            aria-label="Close menu"
            onClick={onCloseMobile}
          >
            <X className="size-4" />
          </button>
        ) : null}
      </div>

      <nav
        className={cn(
          "flex flex-1 flex-col gap-1 py-3",
          collapsed ? "px-1.5" : "px-2",
        )}
        aria-label="Main"
      >
        {navItems.map(({ to, label, icon: Icon, ...rest }) => (
          <NavLink
            key={to}
            to={to}
            end={"end" in rest ? rest.end : false}
            title={collapsed ? label : undefined}
            aria-label={collapsed ? label : undefined}
            className={({ isActive }) =>
              cn(
                "group flex items-center rounded-md text-[13px] font-medium tracking-[-0.01em] transition-colors",
                collapsed
                  ? "h-9 justify-center px-0"
                  : "gap-2.5 px-2.5 py-2",
                isActive
                  ? "bg-sidebar-accent text-sidebar-accent-foreground"
                  : "text-sidebar-foreground/70 hover:bg-sidebar-accent/70 hover:text-sidebar-accent-foreground",
              )
            }
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
                <span className={cn("truncate", collapsed && "sr-only")}>
                  {label}
                </span>
              </>
            )}
          </NavLink>
        ))}
      </nav>
    </div>
  );
}

export function AppShell({ user, onLogout, children }: AppShellProps) {
  const location = useLocation();
  const [mobileOpen, setMobileOpen] = useState(false);
  const [collapsed, setCollapsed] = useState(readCollapsed);
  const meta = pageMeta[location.pathname] ?? {
    title: "Sailorport",
    description: "",
  };
  const displayName = user.name || user.email;

  useEffect(() => {
    setMobileOpen(false);
  }, [location.pathname]);

  function toggleCollapsed() {
    setCollapsed((prev) => {
      const next = !prev;
      writeCollapsed(next);
      return next;
    });
  }

  return (
    <div className="flex min-h-svh bg-background text-foreground">
      <aside
        className={cn(
          "relative sticky top-0 z-40 hidden h-svh shrink-0 border-r border-sidebar-border bg-sidebar transition-[width] duration-200 ease-out md:block",
          collapsed ? "w-14" : "w-[220px]",
        )}
      >
        <SidebarNav collapsed={collapsed} />
        <button
          type="button"
          className={cn(
            "absolute top-3 -right-3 z-50 inline-flex size-6 items-center justify-center rounded-full",
            "border border-border bg-background text-muted-foreground shadow-sm",
            "transition-colors hover:bg-muted hover:text-foreground",
            "outline-none focus-visible:ring-2 focus-visible:ring-ring/40",
          )}
          aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
          title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
          onClick={toggleCollapsed}
        >
          {collapsed ? (
            <ChevronRight className="size-3.5" />
          ) : (
            <ChevronLeft className="size-3.5" />
          )}
        </button>
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
            <SidebarNav
              collapsed={false}
              onCloseMobile={() => setMobileOpen(false)}
            />
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

          <div className="flex items-center gap-1.5">
            <ThemeToggle />

            <DropdownMenu>
              <DropdownMenuTrigger
                className={cn(
                  "inline-flex h-8 max-w-[200px] items-center gap-2 rounded-md px-1.5 text-left",
                  "text-muted-foreground transition-colors hover:bg-muted hover:text-foreground",
                  "outline-none focus-visible:ring-2 focus-visible:ring-ring/40",
                  "data-popup-open:bg-muted data-popup-open:text-foreground",
                )}
              >
                <span className="flex size-7 shrink-0 items-center justify-center rounded-full bg-muted text-[11px] font-semibold tracking-tight text-foreground">
                  {userInitials(displayName)}
                </span>
                <span className="hidden min-w-0 sm:block">
                  <span className="block truncate text-[13px] font-medium text-foreground">
                    {displayName}
                  </span>
                </span>
                <ChevronDown className="hidden size-3.5 shrink-0 opacity-60 sm:block" />
              </DropdownMenuTrigger>

              <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuLabel className="space-y-0.5">
                  <p className="truncate text-[13px] font-medium text-foreground">
                    {displayName}
                  </p>
                  <p className="truncate text-[12px]">{user.email}</p>
                  <p className="pt-0.5 text-[11px] font-medium capitalize">
                    {user.role}
                  </p>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  className="text-destructive data-highlighted:bg-destructive/10 data-highlighted:text-destructive"
                  onClick={onLogout}
                >
                  <LogOut />
                  Log out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </header>

        <main className="app-harbour flex-1 px-4 py-5 sm:px-6 sm:py-6">
          <div className="mx-auto w-full max-w-5xl">{children}</div>
        </main>
      </div>
    </div>
  );
}
