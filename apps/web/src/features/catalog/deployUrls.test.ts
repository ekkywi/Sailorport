import { expect, it } from "vitest";
import {
    buildDeployUrl,
    canShowDeployLinks,
    normalizeDeployHost,
    defaultDeployHost,
} from "./deployUrls";

it("builds open and health urls", () => {
    expect(buildDeployUrl({ host: "localhost", hostPort: 18080, kind: "open" }))
        .toBe("http://localhost:18080");
    expect(buildDeployUrl({ host: "localhost", hostPort: 18080, kind: "healthz"}))
        .toBe("http://localhost:18080/healthz");
});

it("normalizes host", () => {
    expect(normalizeDeployHost("")).toBe("localhost");
    expect(normalizeDeployHost("http://vps.example")).toBe("vps.example");
});

it("gates links by status/source", () => {
    expect(canShowDeployLinks("running", 18080, "git")).toBe(true);
    expect(canShowDeployLinks("failed", 18080, "git")).toBe(false);
    expect(canShowDeployLinks("running", null, "catalog_app")).toBe(false);
});

it("defaultDeployHost falls back and normalizes", () => {
    expect(defaultDeployHost("")).toBe("localhost");
    expect(defaultDeployHost("http://vps.example")).toBe("vps.example");
    expect(defaultDeployHost("  worker-1  ")).toBe("worker-1");
});
