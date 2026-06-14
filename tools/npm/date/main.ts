import { binaryTimeNow } from "../../../dist/index.js";

function main() {
  const now = binaryTimeNow();
  console.log("Binary time now:", now);
  console.log("Binary time now (string):", now.toString());
}

main();
