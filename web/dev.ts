import index from "./index.html";

const backend = `http://${process.env.PREVIEW_BACKEND ?? "localhost:8080"}`;

// Serves the dashboard with hot reloading and proxies the API to a running
// `preview serve`.
const server = Bun.serve({
  port: Number(process.env.PORT ?? 5173),
  development: true,
  routes: {
    "/api/*": (req) => {
      const url = new URL(req.url);
      return fetch(backend + url.pathname + url.search, req);
    },
    "/*": index,
  },
});
console.log(`dashboard on ${server.url} (API: ${backend})`);
