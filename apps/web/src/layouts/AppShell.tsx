import { useEffect, useState, type ReactNode } from "react";
import { NavLink, useLocation } from "react-router-dom";
import {
  Boxes,
  ChevronDown,
  ClipboardList,
  LayoutDashboard,
  LogOut,
  ChevronLeft,
  ChevronRight,
  Menu,
  Server,
  Users,
  X,
  type LucideIcon,
} from "lucide-react";
import { BrandMark, SectionLabel, userInitials } from "@/components/app";
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
import { isAdmin } from "@/lib/rbac";
import { cn } from "@/lib/utils";

type AppShellProps = {
  user: AuthUser;
  onLogout: () => void;
  children: ReactNode;
};

const SIDEBAR_COLLAPSED_KEY = "sailorport.sidebar.collapsed";

type NavItem = {
  to: string;
  label: string;
  icon: LucideIcon;
  end?: boolean;
  adminOnly?: boolean;
};

type NavSection = {
  id: string;
  label: string;
  items: NavItem[];
};

const navSections: NavSection[] = [
  {
    id: "workspace",
    label: "Workspace",
    items: [
      { to: "/overview", label: "Overview", icon: LayoutDashboard, end: true },
    ],
  },
  {
    id: "platform",
    label: "Platform",
    items: [
      { to: "/catalog", label: "Catalog", icon: Boxes },
      { to: "/worker", label: "Workers", icon: Server },
    ],
  },
  {
    id: "admin",
    label: "Administration",
    items: [
      { to: "/users", label: "Users", icon: Users, adminOnly: true },
      { to: "/audit", label: "Audit", icon: ClipboardList, adminOnly: true },
    ],
  },
];

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
  "/users": {
    title: "Users",
    description: "Manage accounts and roles",
  },
  "/audit": {
    title: "Audit log",
    description: "Administrative actions and catalog changes",
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

function SidebarNavLink({
  to,
  label,
  icon: Icon,
  end,
  collapsed,
}: {
  to: string;
  label: string;
  icon: LucideIcon;
  end?: boolean;
  collapsed: boolean;
}) {
  return (
    <NavLink
      to={to}
      end={end ?? false}
      title={collapsed ? label : undefined}
      aria-label={collapsed ? label : undefined}
      className={({ isActive }) =>
        cn(
          "group relative flex items-center rounded-md text-[13px] font-medium tracking-[-0.01em] outline-none transition-colors",
          "focus-visible:ring-2 focus-visible:ring-sidebar-ring/50",
          collapsed ? "h-9 justify-center px-0" : "h-8 gap-2.5 px-2.5",
          isActive
            ? "bg-sidebar-accent text-sidebar-accent-foreground"
            : "text-sidebar-foreground/65 hover:bg-sidebar-accent/60 hover:text-sidebar-accent-foreground",
        )
      }
    >
      {({ isActive }) => (
        <>
          {!collapsed && isActive ? (
            <span
              className="absolute inset-y-1.5 left-0 w-[2px] rounded-full bg-sidebar-primary"
              aria-hidden
            />
          ) : null}
          <Icon
            className={cn(
              "size-4 shrink-0 transition-colors",
              isActive
                ? "text-sidebar-primary"
                : "text-muted-foreground group-hover:text-sidebar-accent-foreground",
            )}
            strokeWidth={isActive ? 2.25 : 2}
          />
          <span className={cn("min-w-0 truncate", collapsed && "sr-only")}>
            {label}
          </span>
        </>
      )}
    </NavLink>
  );
}

function SidebarNav({
  collapsed,
  user,
  onCloseMobile,
}: {
  collapsed: boolean;
  user: AuthUser;
  onCloseMobile?: () => void;
}) {
  const sections = navSections
    .map((section) => ({
      ...section,
      items: section.items.filter(
        (item) => !item.adminOnly || isAdmin(user.role),
      ),
    }))
    .filter((section) => section.items.length > 0);

  return (
    <div className="flex h-full flex-col">
      <div
        className={cn(
          "flex h-12 shrink-0 items-center border-b border-sidebar-border",
          collapsed ? "justify-center px-2" : "gap-2.5 px-3",
        )}
      >
        <BrandMark />
        {!collapsed ? (
          <div className="min-w-0 flex-1 overflow-hidden">
            <p className="truncate text-[13px] font-semibold tracking-[-0.02em] text-sidebar-foreground">
              Sailorport
            </p>
            <p className="truncate text-[11px] text-muted-foreground">
              Control plane
            </p>
          </div>
        ) : null}

        {onCloseMobile ? (
          <button
            type="button"
            className="ml-auto inline-flex size-7 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground md:hidden"
            aria-label="Close menu"
            onClick={onCloseMobile}
          >
            <X className="size-4" />
          </button>
        ) : null}
      </div>

      <nav
        className={cn(
          "flex min-h-0 flex-1 flex-col overflow-y-auto py-3",
          collapsed ? "gap-2.5 px-1.5" : "gap-5 px-2",
        )}
        aria-label="Main"
      >
        {sections.map((section, index) => (
          <div key={section.id} className="flex flex-col gap-0.5">
            {collapsed ? (
              <>
                {index > 0 ? (
                  <div
                    className="mx-1.5 mb-1.5 h-px bg-sidebar-border"
                    aria-hidden
                  />
                ) : null}
                <span className="sr-only">{section.label}</span>
              </>
            ) : (
              <SectionLabel className="mb-1 px-2.5">
                {section.label}
              </SectionLabel>
            )}
            {section.items.map((item) => (
              <SidebarNavLink
                key={item.to}
                to={item.to}
                label={item.label}
                icon={item.icon}
                end={item.end}
                collapsed={collapsed}
              />
            ))}
          </div>
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
          collapsed ? "w-14" : "w-[232px]",
        )}
      >
        <SidebarNav collapsed={collapsed} user={user} />
        <button
          type="button"
          className={cn(
            "absolute top-3.5 -right-3 z-50 inline-flex size-6 items-center justify-center rounded-full",
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
          <aside className="absolute inset-y-0 left-0 w-[260px] border-r border-sidebar-border bg-sidebar shadow-lg">
            <SidebarNav
              collapsed={false}
              user={user}
              onCloseMobile={() => setMobileOpen(false)}
            />
          </aside>
        </div>
      ) : null}

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="sticky top-0 z-30 flex h-12 items-center gap-3 border-b border-border bg-background/85 px-4 backdrop-blur-md sm:px-6 lg:px-8">
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

        <main className="app-harbour flex-1 px-4 py-5 sm:px-6 sm:py-6 lg:px-8">
          <div className="w-full min-w-0">{children}</div>
        </main>
      </div>
    </div>
  );
}
