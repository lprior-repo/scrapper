import { readFileSync } from "fs";
import { join } from "path";

const PORT = 3000;

// Serve static files and React app
Bun.serve({
  port: PORT,
  fetch(req) {
    const url = new URL(req.url);
    const pathname = url.pathname;

    try {
      // Handle static files
      if (pathname.startsWith("/dist/")) {
        const filePath = join(process.cwd(), pathname);
        const file = readFileSync(filePath);
        
        // Set appropriate content type
        const ext = pathname.split('.').pop();
        const contentType = ext === 'js' ? 'application/javascript' : 
                           ext === 'css' ? 'text/css' : 
                           'application/octet-stream';
        
        return new Response(file, {
          headers: { "Content-Type": contentType },
        });
      }
      
      // Serve index.html for all other routes (SPA)
      const html = readFileSync(join(process.cwd(), "index.html"), "utf-8");
      return new Response(html, {
        headers: { "Content-Type": "text/html" },
      });
    } catch (error) {
      console.error("Error serving file:", error);
      return new Response("File not found", { status: 404 });
    }
  },
});

console.log(`🚀 GitHub Codeowner Visualization running at http://localhost:${PORT}`);
console.log(`📊 View the interactive graph at http://localhost:${PORT}`);