import type { ReactNode } from "react";
import { Anchor, Layers, Ship } from "lucide-react";

type AuthLayoutProps = {
  children: ReactNode;
};

export function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <aside className="relative hidden overflow-hidden bg-gradient-to-br from-teal-700 via-teal-800 to-slate-900 lg:flex lg:flex-col lg:justify-between">
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(255,255,255,0.12),transparent_55%)]" />
        <div className="relative flex flex-1 flex-col justify-between p-10 text-white">
          <div className="flex items-center gap-3">
            <div className="flex size-10 items-center justify-center rounded-xl bg-white/15 ring-1 ring-white/20">
              <Ship className="size-5" />
            </div>
            <div>
              <p className="text-lg font-semibold tracking-tight">Sailorport</p>
              <p className="text-sm text-teal-100/80">Internal Developer Platform</p>
            </div>
          </div>

          <div className="space-y-8">
            <div className="space-y-3">
              <h1 className="max-w-md text-3xl font-semibold leading-tight tracking-tight">
                Golden path untuk tim engineering Anda
              </h1>
              <p className="max-w-md text-sm leading-relaxed text-teal-100/85">
                Kelola software catalog, scaffold service baru, dan deploy dengan
                standar yang konsisten dari satu portal.
              </p>
            </div>

            <ul className="space-y-4 text-sm text-teal-50/90">
              <li className="flex items-start gap-3">
                <Layers className="mt-0.5 size-4 shrink-0" />
                <span>Catalog terpusat untuk semua service dan komponen</span>
              </li>
              <li className="flex items-start gap-3">
                <Anchor className="mt-0.5 size-4 shrink-0" />
                <span>Scaffold dari template golden path yang sudah distandarkan</span>
              </li>
            </ul>
          </div>

          <p className="text-xs text-teal-100/60">
            Self-hosted IDP · Open source · Built for platform teams
          </p>
        </div>
      </aside>

      <main className="flex items-center justify-center bg-background p-6 sm:p-10">
        <div className="w-full max-w-md">{children}</div>
      </main>
    </div>
  );
}
