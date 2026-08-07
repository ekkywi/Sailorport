import type { ReactNode } from "react";
import { Link } from "react-router-dom";
import { BrandMark } from "@/components/app";
import { ThemeToggle } from "@/components/ThemeToggle";

type AuthLayoutProps = {
  mode: "login" | "register";
  children: ReactNode;
};

export function AuthLayout({ mode, children }: AuthLayoutProps) {
  const headline =
    mode === "login" ? "Sign in to Sailorport" : "Create your Sailorport account";
  const support =
    mode === "login"
      ? "Welcome back. Please enter your details."
      : "Get started with your team's developer port.";

  return (
    <div className="auth-harbour relative flex min-h-svh flex-col">
      <div className="absolute top-5 right-5 z-10 sm:top-6 sm:right-6">
        <ThemeToggle />
      </div>

      <main className="flex flex-1 items-center justify-center px-6 py-20">
        <div className="w-full max-w-[352px]">
          <div className="mb-10 flex flex-col items-center text-center">
            <Link
              to="/login"
              className="mb-7 inline-flex items-center gap-2 text-foreground"
            >
              <BrandMark />
              <span className="text-[13px] font-medium tracking-[-0.01em]">
                Sailorport
              </span>
            </Link>
            <h1 className="text-[22px] leading-7 font-semibold tracking-[-0.025em] text-foreground">
              {headline}
            </h1>
            <p className="mt-2 max-w-[28ch] text-[13px] leading-5 text-muted-foreground">
              {support}
            </p>
          </div>
          <div className="text-left">{children}</div>
        </div>
      </main>
    </div>
  );
}
