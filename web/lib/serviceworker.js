const CACHE_NAME = 'files-md-v1';
const urlsToCache = [
    '/',
    '/favicon.ico',
    '/manifest.json',
    '/app.css',
    '/lib/normalize.css',
    '/lib/sidebar.css',
    '/lib/codemirror.css',
    '/lib/hypermd.css',
    '/lib/theme-light.css',
    '/lib/theme-dark.css',
    '/chat.css',
    '/lib/sidebar.js',
    '/lib/codemirror.js',
    '/lib/core.js',
    '/lib/markdown.js',
    '/lib/hypermd.js',
    '/lib/keymap.js',
    '/lib/click.js',
    '/lib/hide-token.js',
    '/lib/fold.js',
    '/lib/fold-image.js',
    '/lib/fold-link.js',
    '/lib/table-align.js',
    '/lib/autocomplete-link.js',
    '/lib/show-hint.js',
    '/lib/autoscroll.js',
    '/lib/codemirror-go.js',
    '/lib/codemirror-php.js',
    '/lib/codemirror-shell.js',
    '/lib/similarity.js',
    '/welcome.js',
    '/files.js',
    '/wasm_exec.js',
    '/app.js',
    '/wasm.js',
    '/chat.js',
    '/modals.js',
];

self.addEventListener('install', event => {
    event.waitUntil(
        caches.open(CACHE_NAME)
            .then(cache => {
                // Cache each file individually to find the problem
                const cachePromises = urlsToCache.map(url => {
                    console.log('Trying to cache:', url);
                    return cache.add(url + window.COMMIT_HASH)
                        .then(() => console.log('✓ Cached:', url))
                        .catch(err => console.error('✗ Failed to cache:', url, err));
                });
                return Promise.allSettled(cachePromises); // Won't fail if one fails
            })
    );
});

self.addEventListener('fetch', event => {
    console.log('intercepting');
    event.respondWith(
        caches.match(event.request)
            .then(response => {
                return response || fetch(event.request);
            })
    );
});