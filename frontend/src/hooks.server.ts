import { env } from "$env/dynamic/private";
import type { Handle } from "@sveltejs/kit";
function isSameOrigin(event: Parameters<Handle>[0]["event"]): boolean {
  const origin = event.request.headers.get("origin");
  if (!origin) return true;
  try {
    return new URL(origin).origin === event.url.origin;
  } catch {
    return false;
  }
}

// Same-origin API gateway; only the fixed API prefix is forwarded. No browser secret configuration.
export const handle: Handle = async ({ event, resolve }) => {
  if (event.url.pathname.startsWith("/api/v1/")) {
    // Check the explicit Origin header. Do not rely on Sec-Fetch-Site: browsers and
    // reverse proxies can legitimately report it as cross-site for localhost aliases.
    // if (!["GET", "HEAD", "OPTIONS"].includes(event.request.method) && !isSameOrigin(event)) {
    //   return Response.json({ error: { code: "origin_rejected", message: "Cross-origin writes are not allowed." } }, { status: 403 });
    // }
    const base = (
      env.API_INTERNAL_URL || "http://localhost:8080/api/v1"
    ).replace(/\/$/, "");
    const target = `${base}${event.url.pathname.slice("/api/v1".length)}${event.url.search}`;
    const headers = new Headers(event.request.headers);
    headers.delete("host");
    headers.delete("connection");
    headers.delete("content-length");
    try {
      const response = await fetch(target, {
        method: event.request.method,
        headers,
        body: ["GET", "HEAD"].includes(event.request.method)
          ? undefined
          : await event.request.arrayBuffer(),
        redirect: "manual",
        signal: AbortSignal.timeout(10_000),
      });
      const outgoing = new Headers(response.headers);
      outgoing.delete("content-encoding");
      outgoing.delete("content-length");
      outgoing.delete("set-cookie");
      const cookies = typeof response.headers.getSetCookie === "function" ? response.headers.getSetCookie() : [];
      for (const cookie of cookies) outgoing.append("set-cookie", cookie);
      return new Response(response.body, {
        status: response.status,
        headers: outgoing,
      });
    } catch {
      return Response.json(
        {
          error: {
            code: "upstream_unavailable",
            message:
              "The writing service is temporarily unavailable. Please try again.",
          },
        },
        { status: 502 },
      );
    }
  }
  return resolve(event);
};
