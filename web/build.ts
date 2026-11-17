import { rm } from "node:fs/promises";
import tailwind from "bun-plugin-tailwind";

const outdir = "dist";

await rm(outdir, { recursive: true, force: true });
const result = await Bun.build({
  entrypoints: ["./index.html"],
  outdir,
  minify: true,
  sourcemap: "linked",
  plugins: [tailwind],
  define: { "process.env.NODE_ENV": JSON.stringify("production") },
  // Absolute URLs: the Go server answers deep SPA routes with index.html.
  publicPath: "/",
  // The Go server caches /assets/ as immutable, so every hashed file goes there.
  naming: {
    entry: "[name].[ext]",
    chunk: "assets/[name]-[hash].[ext]",
    asset: "assets/[name]-[hash].[ext]",
  },
});
if (!result.success) {
  for (const log of result.logs) console.error(log);
  process.exit(1);
}
for (const output of result.outputs) {
  console.log(`${output.path.slice(process.cwd().length + 1)}  ${output.size} B`);
}
