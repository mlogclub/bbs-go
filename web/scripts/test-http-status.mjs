// Run after pnpm build:ssr. Exercise the built SSR handler with a controlled API.
import assert from "node:assert/strict"
import { createRequestHandler } from "react-router"
process.env.BBSGO_SERVER_URL = "http://upstream.test"
const build = await import("../build/server/index.js")
const handler = createRequestHandler(build, "production")
let upstream = "ok"
let currentUser = null
const originalFetch = globalThis.fetch
const originalError = console.error
console.error = () => {}
globalThis.fetch = async (input) => {
  const pathname = new URL(typeof input === "string" ? input : input.url)
    .pathname
  if (upstream === "install") return new Response("", { status: 428 })
  if (upstream === "root503") return new Response("", { status: 503 })
  if (pathname === "/api/config/configs")
    return Response.json({ success: true, data: {} })
  if (pathname === "/api/user/current")
    return Response.json({ success: true, data: currentUser })
  if (upstream === "http503") return new Response("", { status: 503 })
  if (upstream === "http404") return new Response("", { status: 404 })
  if (upstream === "network") throw new TypeError("fetch failed")
  if (upstream === "business404")
    return Response.json({
      success: false,
      errorCode: 404,
      message: "Resource missing",
    })
  return Response.json({
    success: true,
    data: { id: 1, title: "Test article", content: "", tags: [] },
  })
}
try {
  for (const [path, mode, status, method = "GET"] of [
    ["/missing-route", "ok", 404],
    ["/missing-route", "ok", 404, "HEAD"],
    ["/assets/missing.js", "ok", 404],
    ["/article/1", "ok", 200],
    ["/article/1", "business404", 404],
    ["/topic/1", "business404", 404],
    ["/user/1", "http404", 404],
    ["/article/1", "http503", 503],
    ["/article/1", "network", 500],
    ["/user/signin", "root503", 503],
    ["/redirect", "ok", 400],
    ["/redirect?url=https://example.com", "ok", 200],
    ["/install", "install", 200],
  ]) {
    upstream = mode
    const response = await handler(
      new Request("http://localhost" + path, { method })
    )
    const body = await response.text()
    assert.equal(response.status, status, `${method} ${path} (${mode})`)
    if (method === "HEAD") assert.equal(body, "")
    console.log(`PASS ${method} ${path} (${mode}): ${status}`)
  }
  upstream = "ok"
  currentUser = { id: 1, permissions: ["dashboard.view"] }
  for (const path of [
    "/dashboard/articles",
    "/dashboard/categories",
    "/dashboard/settings",
  ]) {
    const response = await handler(new Request("http://localhost" + path))
    await response.text()
    assert.equal(response.status, 403, path)
    console.log(`PASS ${path}: 403`)
  }
} finally {
  globalThis.fetch = originalFetch
  console.error = originalError
}
