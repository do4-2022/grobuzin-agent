import { createServer } from "node:http";

//@ts-ignore
import { run } from "./function";

const server = createServer((req, res) => {
  // Call the function
  const result = run(req);

  res.writeHead(200, { "Content-Type": "application/json" });

  // There may be a problem with binary responses
  res.end(JSON.stringify(result));
});

// starts a simple http server locally on port 3000
server.listen(3000, "127.0.0.1", () => {
  console.log("Listening on 127.0.0.1:3000");
});
