import { Relayer } from "./relayer"

export async function startRelayer(url: string) {
  const relayer = new Relayer(url);
  process.on("SIGINT", async function () {
    await relayer.stop();
  });
  process.on("SIGTERM", async function () {
    await relayer.stop();
  });
  process.on("SIGQUIT", async function () {
    await relayer.stop();
  });

  await relayer.start();

  return relayer;
}