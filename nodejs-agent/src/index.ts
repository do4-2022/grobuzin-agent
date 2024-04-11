import { createServer } from "node:http";

//@ts-ignore
import { run } from "./function";

const server = createServer(async (req, res) => {
  try {
    // Call the function
    const result = await run(req);

    res.writeHead(200, { "Content-Type": "application/json" });

    // There may be a problem with binary responses
    res.end(JSON.stringify(result));
  } catch (e) {
    console.error(e);
    res.writeHead(500, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ message: JSON.stringify(e) }));
  }
});

// starts a simple http server locally on port 3000
server.listen(3000, "127.0.0.1", () => {
  console.log("Listening on 127.0.0.1:3000");
});
