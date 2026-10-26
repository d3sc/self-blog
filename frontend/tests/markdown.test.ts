import { describe, it, expect } from "vitest";
import { renderMarkdown } from "../src/lib/markdown";
describe("safe Markdown", () => {
  it("renders editorial formatting", () =>
    expect(renderMarkdown("# Heading\n\n**Strong**")).toContain(
      "<strong>Strong</strong>",
    ));
  it("removes scripts and event handlers", () => {
    const html = renderMarkdown(
      '<script>alert(1)</script><img src=x onerror="alert(1)">',
    );
    expect(html).not.toMatch(/<script|onerror/i);
  });
  it("removes javascript URLs", () =>
    expect(renderMarkdown("[bad](javascript:alert%281%29)")).not.toContain(
      'href="javascript:',
    ));
});
