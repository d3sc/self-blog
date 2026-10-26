import { describe, it, expect, vi, afterEach } from "vitest";
import { api } from "../src/lib/api";
afterEach(() => vi.unstubAllGlobals());
describe("API client", () => {
  it("sends the backend CSRF cookie on mutations", async () => {
    vi.stubGlobal("document", { cookie: "blog_csrf=token%20value" });
    const fetch = vi
      .fn()
      .mockResolvedValue(new Response(null, { status: 204 }));
    vi.stubGlobal("fetch", fetch);
    await api.logout();
    expect(
      new Headers(fetch.mock.calls[0][1].headers).get("X-CSRF-Token"),
    ).toBe("token value");
    expect(fetch.mock.calls[0][1].credentials).toBe("include");
  });
});
