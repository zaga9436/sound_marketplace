type WSOptions = {
  token: string;
  onMessage?: (event: MessageEvent) => void;
  onOpen?: () => void;
  onClose?: () => void;
};

export function createWebSocketClient(path: string, options: WSOptions) {
  const url = new URL(path, window.location.origin);
  url.protocol = window.location.protocol === "https:" ? "wss:" : "ws:";

  const socket = new WebSocket(url.toString());

  socket.addEventListener("open", () => {
    socket.send(JSON.stringify({ type: "auth", token: options.token }));
    options.onOpen?.();
  });
  socket.addEventListener("message", (event) => options.onMessage?.(event));
  socket.addEventListener("close", () => options.onClose?.());

  return socket;
}
