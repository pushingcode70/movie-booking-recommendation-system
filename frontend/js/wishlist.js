/* wishlist & watched views controller with explicit route links & 90+ min filter */

async function renderWishlist(container) {
  if (!API.getToken()) {
    window.location.hash = '#/login';
    return;
  }

  container.innerHTML = '<div style="text-align: center; padding: 3rem; color: var(--text-muted);">Loading wishlist...</div>';

  try {
    // get /wishlist returns array of wishlist item response
    const items = await API.get('/wishlist');
    const filteredItems = (items || []).filter(item => {
      const m = item.movie || item;
      const dur = m.duration || m.runtime || 0;
      return dur === 0 || dur >= 90; // filter out movies under 90 mins
    });

    container.innerHTML = `
      <h1 class="page-title">My Wishlist</h1>
      ${filteredItems && filteredItems.length > 0 ? `
        <div class="grid">
          ${filteredItems.map(item => {
            const m = item.movie || item;
            const tmdbId = item.tmdb_id || m.tmdb_id || m.TMDBID || m.id;
            const title = m.title || m.Title || 'Untitled';
            const posterPath = m.poster_path || m.PosterPath;

            const poster = posterPath
              ? (posterPath.startsWith('http') ? posterPath : `https://image.tmdb.org/t/p/w500${posterPath}`)
              : 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?q=80&w=1000';

            const durationStr = m.duration ? `${m.duration} mins` : (m.runtime ? `${m.runtime} mins` : '');
            const langStr = (m.language || 'EN').toUpperCase();

            // wishlist items represent tmdb items, navigate explicitly to #/tmdb/<tmdb_id>
            return `
              <a href="#/tmdb/${tmdbId}" class="movie-card">
                <img src="${poster}" alt="${title}" class="movie-poster" loading="lazy" />
                <div class="movie-info">
                  <div class="movie-title">${title}</div>
                  <div class="movie-meta">
                    ${langStr} ${durationStr ? '• ' + durationStr : ''}
                  </div>
                </div>
              </a>
            `;
          }).join('')}
        </div>
      ` : `
        <div class="card" style="text-align: center; padding: 3rem; color: var(--text-muted);">
          Your wishlist is currently empty.
        </div>
      `}
    `;
  } catch (err) {
    container.innerHTML = `<div class="card" style="color: var(--accent);">${err.message || 'Failed to load wishlist.'}</div>`;
  }
}

async function renderWatched(container) {
  if (!API.getToken()) {
    window.location.hash = '#/login';
    return;
  }

  container.innerHTML = '<div style="text-align: center; padding: 3rem; color: var(--text-muted);">Loading watched list...</div>';

  try {
    // get /watched returns array of watched movie item response
    const items = await API.get('/watched');
    const filteredItems = (items || []).filter(item => {
      const m = item.movie || item;
      const dur = m.duration || m.runtime || 0;
      return dur === 0 || dur >= 90; // filter out movies under 90 mins
    });

    container.innerHTML = `
      <h1 class="page-title">Watched Movies</h1>
      ${filteredItems && filteredItems.length > 0 ? `
        <div class="grid">
          ${filteredItems.map(item => {
            const m = item.movie || item;
            const title = m.title || m.Title || 'Untitled';
            const posterPath = m.poster_path || m.PosterPath;
            const tmdbId = item.tmdb_id || m.tmdb_id || m.TMDBID || m.id;
            const rating = item.rating || '';
            const review = item.review || '';

            // watched items represent tmdb items, navigate explicitly to #/tmdb/<tmdb_id>
            const route = `#/tmdb/${tmdbId}`;

            const poster = posterPath
              ? (posterPath.startsWith('http') ? posterPath : `https://image.tmdb.org/t/p/w500${posterPath}`)
              : 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?q=80&w=1000';

            return `
              <div class="movie-card watched-poster-card" 
                   data-route="${route}"
                   data-tmdb-id="${tmdbId}"
                   data-title="${title}"
                   data-poster="${poster}"
                   data-rating="${rating}"
                   data-review="${review}"
                   style="cursor: pointer; position: relative;">
                <img src="${poster}" alt="${title}" class="movie-poster" loading="lazy" />
              </div>
            `;
          }).join('')}
        </div>

        <!-- medium sized dialogue box modal -->
        <div id="watched-dialog-overlay" style="display: none; position: fixed; inset: 0; background-color: rgba(0,0,0,0.85); backdrop-filter: blur(4px); z-index: 300; align-items: center; justify-content: center; padding: 1.5rem;">
          <div class="card" style="width: 100%; max-width: 480px; background-color: #121212; border: 1px solid var(--border); border-radius: 12px; box-shadow: 0 25px 50px rgba(0,0,0,0.9); padding: 1.75rem;">
            
            <div style="display: flex; gap: 1.25rem; margin-bottom: 1.25rem; align-items: flex-start;">
              <img id="dialog-poster" src="" alt="Movie" style="width: 80px; aspect-ratio: 2/3; object-fit: cover; border-radius: 6px; border: 1px solid var(--border);" />
              <div style="flex: 1; display: flex; flex-direction: column; gap: 0.25rem;">
                <h2 id="dialog-title" style="font-size: 1.25rem; font-weight: 800; color: #ffffff;">Movie Title</h2>
                <a id="dialog-details-link" href="#" class="btn btn-secondary btn-sm" style="align-self: flex-start; margin-top: 0.5rem;">View Movie Details</a>
              </div>
            </div>

            <form id="form-dialog-watched" style="display: flex; flex-direction: column; gap: 1rem;">
              <input type="hidden" id="dialog-tmdb-id" />
              
              <div class="form-group" style="margin-bottom: 0;">
                <label class="form-label">Your Rating (1 to 10)</label>
                <input type="number" id="dialog-rating" min="1" max="10" class="form-input" placeholder="Enter rating (e.g. 8)" />
              </div>

              <div class="form-group" style="margin-bottom: 0;">
                <label class="form-label">Your Review</label>
                <textarea id="dialog-review" class="form-textarea" rows="3" placeholder="Write your review here..."></textarea>
              </div>

              <div style="display: flex; gap: 0.75rem; justify-content: flex-end; margin-top: 0.5rem;">
                <button type="button" id="btn-dialog-remove" class="btn btn-danger btn-sm" style="margin-right: auto;">Remove</button>
                <button type="button" id="btn-dialog-cancel" class="btn btn-secondary btn-sm">Close</button>
                <button type="submit" class="btn btn-primary btn-sm">Save Changes</button>
              </div>
            </form>

          </div>
        </div>
      ` : `
        <div class="card" style="text-align: center; padding: 3rem; color: var(--text-muted);">
          You haven't marked any movies as watched yet.
        </div>
      `}
    `;

    // attach poster click handlers to open dialogue modal
    const overlay = document.getElementById('watched-dialog-overlay');
    const dialogForm = document.getElementById('form-dialog-watched');
    const cancelBtn = document.getElementById('btn-dialog-cancel');
    const removeBtn = document.getElementById('btn-dialog-remove');

    document.querySelectorAll('.watched-poster-card').forEach(card => {
      card.addEventListener('click', () => {
        const route = card.dataset.route;
        const tmdbId = card.dataset.tmdbId;
        const title = card.dataset.title;
        const poster = card.dataset.poster;
        const rating = card.dataset.rating;
        const review = card.dataset.review;

        document.getElementById('dialog-poster').src = poster;
        document.getElementById('dialog-title').textContent = title;
        document.getElementById('dialog-details-link').href = route;
        document.getElementById('dialog-tmdb-id').value = tmdbId;
        document.getElementById('dialog-rating').value = rating;
        document.getElementById('dialog-review').value = review;

        overlay.style.display = 'flex';
      });
    });

    if (cancelBtn) {
      cancelBtn.addEventListener('click', () => { overlay.style.display = 'none'; });
    }

    if (overlay) {
      overlay.addEventListener('click', (e) => {
        if (e.target === overlay) overlay.style.display = 'none';
      });
    }

    if (removeBtn) {
      removeBtn.addEventListener('click', async () => {
        const tmdbId = document.getElementById('dialog-tmdb-id').value;
        if (!tmdbId) return;
        try {
          await API.delete(`/watched/${tmdbId}`);
          overlay.style.display = 'none';
          renderWatched(container);
        } catch (err) {
          alert(err.message || 'Failed to remove movie.');
        }
      });
    }

    if (dialogForm) {
      dialogForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const tmdbId = document.getElementById('dialog-tmdb-id').value;
        const ratingVal = document.getElementById('dialog-rating').value;
        const reviewVal = document.getElementById('dialog-review').value;

        const body = {};
        if (ratingVal) body.rating = Number(ratingVal);
        if (reviewVal !== undefined) body.review = reviewVal;

        try {
          await API.put(`/watched/${tmdbId}`, body);
          overlay.style.display = 'none';
          renderWatched(container);
        } catch (err) {
          alert(err.message || 'Failed to save review.');
        }
      });
    }

  } catch (err) {
    container.innerHTML = `<div class="card" style="color: var(--accent);">${err.message || 'Failed to load watched list.'}</div>`;
  }
}

Router.addRoute('#/wishlist', renderWishlist);
Router.addRoute('#/watched', renderWatched);
