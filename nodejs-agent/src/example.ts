// example function

import { IncomingMessage } from "http";

export async function run(req: IncomingMessage) {
  console.log("Hello World!");

  // Return a simple JSON response
  return {
    status: 200,
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ message: "Hello World!" }),
  };
}
