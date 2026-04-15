// Minimal service worker required for PWA install prompt.
// No caching strategy — just satisfies the browser requirement.
self.addEventListener('fetch', () => {})
