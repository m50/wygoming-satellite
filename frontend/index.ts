import htmx from 'htmx.org';

// Init HTMX
declare global {
  interface Window { htmx: typeof htmx }
}
window.htmx = htmx

// Load HTMX Extensions
require('htmx.org/dist/ext/ws');

