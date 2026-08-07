import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { authFieldClass, authLabelClass } from "./styles";

type PasswordFieldProps = {
  id: string;
  name?: string;
  label?: string;
  disabled?: boolean;
  autoComplete?: string;
  invalid?: boolean;
  placeholder?: string;
  className?: string;
};

export function PasswordField({
  id,
  name = "password",
  label = "Password",
  disabled,
  autoComplete = "current-password",
  invalid,
  placeholder = "At least 8 characters",
  className,
}: PasswordFieldProps) {
  const [show, setShow] = useState(false);

  return (
    <div className={cn("space-y-1.5", className)}>
      <Label htmlFor={id} className={authLabelClass}>
        {label}
      </Label>
      <div className="relative">
        <Input
          id={id}
          name={name}
          type={show ? "text" : "password"}
          autoComplete={autoComplete}
          required
          minLength={8}
          placeholder={placeholder}
          disabled={disabled}
          aria-invalid={invalid ? true : undefined}
          className={cn(authFieldClass, "pr-10")}
        />
        <button
          type="button"
          className="absolute top-1/2 right-2.5 -translate-y-1/2 bg-transparent p-0 text-muted-foreground/80 transition-colors hover:text-foreground disabled:opacity-50"
          onClick={() => setShow((v) => !v)}
          disabled={disabled}
          aria-label={show ? "Hide password" : "Show password"}
          tabIndex={-1}
        >
          {show ? <EyeOff className="size-3.5" /> : <Eye className="size-3.5" />}
        </button>
      </div>
    </div>
  );
}
