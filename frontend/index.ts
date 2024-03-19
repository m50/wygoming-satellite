import htmx from 'htmx.org';

// Init HTMX
declare global {
  interface Window { htmx: typeof htmx }
}
window.htmx = htmx;

// Load HTMX Extensions
require('htmx.org/dist/ext/ws');

let ws: WebSocket | null = null;

const WsMessageTypes = {
  Echo: "echo",
  Message: "message",
  Binary: "binary",
  Close: "close",
} as const;

interface WSOpenEvent extends Event {
  detail: {
    elt: HTMLElement;
    event: {
      currentTarget: WebSocket;
      target: WebSocket;
      srcElement: WebSocket;
    };
    socketWrapper: any;
  };
}
htmx.on('htmx:wsOpen', (evt: Event) => {
  const e = evt as unknown as WSOpenEvent;
  ws = e.detail.event.target;
});

htmx.on('htmx:wsClose', (evt: Event) => {
  // console.log(evt);
});

htmx.on('htmx:beforeSwap', (evt: Event) => {
  if (ws) {
    ws.send(JSON.stringify({
      type: WsMessageTypes.Close
    }));
    ws.close()
    ws = null
  }
  // console.log(evt)
})

htmx.onLoad(() => {
  const sidebar = document.getElementById('sidebar');
  const burgerMenu = document.getElementById('sidebar-burger');
  const options = document.getElementById('sidebar-options');
  const content = document.getElementById('content');
  burgerMenu?.addEventListener('click', (evt) => {
    // console.log(evt);
    evt.stopImmediatePropagation()
    const isVisible = options?.classList.contains('flex');
    console.log(isVisible);
    if (isVisible) {
      options?.classList.replace('flex', 'hidden');
      sidebar?.classList.remove('col-span-3');
      content?.classList.remove('hidden');
      content?.classList.add('col-span-3');
    } else {
      options?.classList.replace('hidden', 'flex');
      sidebar?.classList.add('col-span-3');
      content?.classList.add('hidden');
      content?.classList.remove('col-span-3');
    }
  })
});
