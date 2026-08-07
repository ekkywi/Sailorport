import { useEffect, useState } from "react";
import { Moon, Sun } from "lucide-react";
import { cn } from "@/lib/utils";
import { getStoredTheme, toggleTheme, type Theme } from "@/lib/theme";

export function ThemeToggle({ className }: { className?: string }) {
  const [theme, setTheme] = useState<Theme>("light");

  useEffect(() => {
    setTheme(getStoredTheme());
  }, []);

  return (
    <button
      type="button"
      className={cn(
        "inline-flex size-8 items-center justify-center rounded-md border-0 bg-transparent p-0 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground",
        className,
      )}
      aria-label={theme === "dark" ? "Switch to light mode" : "Switch to dark mode"}
      onClick={() => setTheme(toggleTheme())}
    >
      {theme === "dark" ? <Sun className="size-4" /> : <Moon className="size-4" />}
    </button>
  );
}
