/* Hash-based SPA Router & Global Search Navigation Controller */
const Router = {
  routes: {},

  addRoute(path, handler) {
    this.routes[path] = handler;
  },

  init() {
    window.addEventListener('hashchange', () => this.handleRoute());
    const start = () => {
      this.handleRoute();
      this.initGlobalSearch();
    };
    if (document.readyState === 'complete' || document.readyState === 'interactive') {
      start();
    } else {
      window.addEventListener('DOMContentLoaded', start);
    }
  },

  initGlobalSearch() {
    const searchInput = document.getElementById('global-search-input');
    const dropdown = document.getElementById('global-search-dropdown');
    if (!searchInput || !dropdown) return;

    let debounceTimer = null;

    searchInput.addEventListener('input', (e) => {
      const query = e.target.value.trim();
      clearTimeout(debounceTimer);

      if (!query || query.length < 2) {
        dropdown.style.display = 'none';
        dropdown.innerHTML = '';
        return;
      }

      debounceTimer = setTimeout(async () => {
        try {
          // Call existing backend TMDB search endpoint: GET /tmdb/search?query=...
          const res = await API.get(`/tmdb/search?query=${encodeURIComponent(query)}`);
          const results = (res && res.results) ? res.results : [];

          if (results.length === 0) {
            dropdown.innerHTML = '<div style="padding: 0.875rem; font-size: 0.8125rem; color: var(--text-muted); text-align: center;">No movies found.</div>';
            dropdown.style.display = 'block';
            return;
          }

          dropdown.innerHTML = results.slice(0, 8).map(movie => {
            const poster = movie.poster_path
              ? (movie.poster_path.startsWith('http') ? movie.poster_path : `https://image.tmdb.org/t/p/w92${movie.poster_path}`)
              : 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?q=80&w=1000';

            const year = movie.release_date ? movie.release_date.split('-')[0] : '';
            const tmdbId = movie.id;

            return `
              <div class="search-dropdown-item" data-tmdb-id="${tmdbId}">
                <img src="${poster}" alt="${movie.title}" class="search-thumb" />
                <div class="search-item-info">
                  <div class="search-item-title">${movie.title}</div>
                  <div class="search-item-meta">${year ? year + ' • ' : ''}TMDB #${tmdbId}</div>
                </div>
              </div>
            `;
          }).join('');

          dropdown.style.display = 'block';
        } catch (err) {
          dropdown.innerHTML = `<div style="padding: 0.875rem; font-size: 0.8125rem; color: var(--brand-primary); text-align: center;">${err.message || 'Search error'}</div>`;
          dropdown.style.display = 'block';
        }
      }, 300);
    });

    // Handle Item Selection
    dropdown.addEventListener('click', (e) => {
      const item = e.target.closest('.search-dropdown-item');
      if (!item) return;
      const tmdbId = item.dataset.tmdbId;
      dropdown.style.display = 'none';
      searchInput.value = '';
      window.location.hash = `#/tmdb/${tmdbId}`;
    });

    // Hide dropdown on outside click
    document.addEventListener('click', (e) => {
      if (!e.target.closest('.global-search-container')) {
        dropdown.style.display = 'none';
      }
    });
  },

  handleRoute() {
    let rawHash = window.location.hash || '#/';

    if (rawHash === '' || rawHash === '#' || rawHash === '#home' || rawHash === '#/home') {
      rawHash = '#/';
    }

    const hashPath = rawHash.split('?')[0];

    this.updateNavbar(hashPath);
    const container = document.getElementById('app');
    if (!container) return;

    // Match exact route (using hashPath stripped of query string parameters)
    if (this.routes[hashPath]) {
      this.routes[hashPath](container);
      return;
    }

    // Match parameterized routes (e.g., #/movies/:id, #/book/:showId, #/bookings/confirm/:id)
    for (const routePath in this.routes) {
      if (routePath.includes(':')) {
        const routeParts = routePath.split('/');
        const hashParts = hashPath.split('/');

        if (routeParts.length === hashParts.length) {
          const params = {};
          let match = true;

          for (let i = 0; i < routeParts.length; i++) {
            if (routeParts[i].startsWith(':')) {
              const paramName = routeParts[i].slice(1);
              params[paramName] = hashParts[i];
            } else if (routeParts[i] !== hashParts[i]) {
              match = false;
              break;
            }
          }

          if (match) {
            this.routes[routePath](container, params);
            return;
          }
        }
      }
    }

    // 404 Fallback
    container.innerHTML = `
      <div class="card" style="text-align: center; padding: 3rem; max-width: 500px; margin: 2rem auto;">
        <h1 style="color: var(--brand-primary); font-size: 2rem; font-weight: 800;">404</h1>
        <p style="color: var(--text-muted); margin-top: 0.5rem;">Page Not Found</p>
        <a href="#/" class="btn btn-primary" style="margin-top: 1.5rem;">Return to Home</a>
      </div>
    `;
  },

  updateNavbar(currentHash) {
    const navLinks = document.getElementById('nav-links');
    if (!navLinks) return;

    const user = API.getUser();
    const isLoggedIn = !!API.getToken();
    const isAdmin = user && (user.role === 'admin' || user.role === 'ADMIN');

    let html = `
      <a href="#/" class="${currentHash === '#/' ? 'active-red' : ''}">Home</a>
      <a href="#/shows" class="${currentHash.startsWith('#/shows') || currentHash.startsWith('#/theatres') ? 'active' : ''}">Shows</a>
    `;

    if (isLoggedIn) {
      html += `
        <a href="#/bookings" class="${currentHash === '#/bookings' ? 'active' : ''}">My Bookings</a>
        <a href="#/wishlist" class="${currentHash === '#/wishlist' ? 'active' : ''}">Wishlist</a>
        <a href="#/watched" class="${currentHash === '#/watched' ? 'active' : ''}">Watched</a>
        <a href="#/profile" class="${currentHash === '#/profile' ? 'active' : ''}">Profile</a>
      `;

      if (isAdmin) {
        html += `<a href="#/admin" class="${currentHash.startsWith('#/admin') ? 'active-red' : ''}" style="color: var(--brand-primary); font-weight: 800;">Admin</a>`;
      }

      html += `<button id="btn-logout" class="btn btn-secondary btn-sm" style="margin-left: 0.5rem;">Logout</button>`;
    } else {
      html += `
        <a href="#/login" class="btn-nav-accent">Login</a>
      `;
    }

    navLinks.innerHTML = html;

    const logoutBtn = document.getElementById('btn-logout');
    if (logoutBtn) {
      logoutBtn.addEventListener('click', () => {
        API.setToken(null);
        window.location.hash = '#/login';
      });
    }
  }
};

Router.init();
