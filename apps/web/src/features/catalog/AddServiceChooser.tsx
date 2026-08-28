import { FolderGit2, LayoutTemplate, Package, Tag } from "lucide-react";

type AddServiceChooserProps = {
  onGit: () => void;
  onCatalog?: () => void;
  onTemplate: () => void;
  onRegister: () => void;
};

function ChoiceCard({
  title,
  description,
  icon: Icon,
  onClick,
}: {
  title: string;
  description: string;
  icon: typeof FolderGit2;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className="flex w-full items-start gap-3 rounded-lg border border-border px-3.5 py-3 text-left transition-colors hover:bg-muted/50"
    >
      <span className="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground">
        <Icon className="size-4" />
      </span>
      <span className="min-w-0 flex-1">
        <span className="block text-[13px] font-medium tracking-[-0.01em]">
          {title}
        </span>
        <span className="mt-0.5 block text-[12px] text-muted-foreground">
          {description}
        </span>
      </span>
    </button>
  );
}

export function AddServiceChooser({
  onGit,
  onCatalog,
  onTemplate,
  onRegister,
}: AddServiceChooserProps) {
  return (
    <div className="space-y-2">
      <ChoiceCard
        title="From Git"
        description="Link a public repo with a Dockerfile. Deploy clones and builds on the agent."
        icon={FolderGit2}
        onClick={onGit}
      />
      {onCatalog ? (
        <ChoiceCard
          title="From catalog"
          description="Platform image (Postgres, …). Deploy pulls and runs — no git or build."
          icon={Package}
          onClick={onCatalog}
        />
      ) : null}
      <ChoiceCard
        title="From template"
        description="Scaffold a workspace from a golden-path template (e.g. go-api)."
        icon={LayoutTemplate}
        onClick={onTemplate}
      />
      <ChoiceCard
        title="Register only"
        description="Catalog metadata only — no workspace and no deploy until you link a source."
        icon={Tag}
        onClick={onRegister}
      />
    </div>
  );
}
