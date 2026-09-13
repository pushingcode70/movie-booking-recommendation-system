/* explicit movie details controller: local movies (#/movies/:id) vs tmdb movies (#/tmdb/:tmdbId) */

// -------------------------------------------------------------
// 1. local published movie details (#/movies/:id)
// uses local postgresql primary key (movies.id)
// -------------------------------------------------------------
async function renderLocalMovieDetails(container, params) {
  const localId = Number(params.id);
  const isLoggedIn = !!API.getToken();

  container.innerHTML = '<div style="text-align: center; padding: 3rem; color: var(--text-muted);">Loading movie details...</div>';

  try {
    // strictly call get /movies/:id (local primary key)
    const localMovie = await API.get(`/movies/${localId}`);
    if (!localMovie || !localMovie.id) {
      throw new Error('Local movie not found');
    }

    const tmdbId = localMovie.tmdb_id;
    let tmdbMovie = null;

    // enrich metadata from tmdb if tmdb_id is present
    if (tmdbId) {
      try {
        tmdbMovie = await API.get(`/tmdb/movie/${tmdbId}`);
      } catch (e) {}
    }

    // extract director & cast
    let director = localMovie.director || '';
    let cast = localMovie.cast || '';

    if (tmdbMovie && tmdbMovie.credits) {
      if (!director && tmdbMovie.credits.crew) {
        const dObj = tmdbMovie.credits.crew.find(c => c.job === 'Director');
        if (dObj) director = dObj.name;
      }
      if (!cast && tmdbMovie.credits.cast) {
        cast = tmdbMovie.credits.cast.slice(0, 5).map(c => c.name).join(', ');
      }
    }

    // consolidated display fields
    const title = localMovie.title || (tmdbMovie && tmdbMovie.title) || 'Movie Details';
    const overview = localMovie.description || localMovie.overview || (tmdbMovie && tmdbMovie.overview) || 'No overview available.';
    const releaseDate = localMovie.release_date || (tmdbMovie && tmdbMovie.release_date) || '';
    const releaseYear = releaseDate ? releaseDate.split('-')[0] : '';
    const runtime = localMovie.duration || localMovie.runtime || (tmdbMovie && tmdbMovie.runtime) || 0;
    const language = (localMovie.language || (tmdbMovie && tmdbMovie.original_language) || 'EN').toUpperCase();
    const voteAverage = tmdbMovie && tmdbMovie.vote_average ? tmdbMovie.vote_average.toFixed(1) : null;

    // posters & backdrops
    const rawPoster = localMovie.poster_path || (tmdbMovie && tmdbMovie.poster_path);
    const poster = rawPoster
      ? (rawPoster.startsWith('http') ? rawPoster : `https://image.tmdb.org/t/p/w500${rawPoster}`)
      : 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?q=80&w=1000';

    const rawBackdrop = localMovie.backdrop_path || (tmdbMovie && tmdbMovie.backdrop_path);
    const backdrop = rawBackdrop
      ? (rawBackdrop.startsWith('http') ? rawBackdrop : `https://image.tmdb.org/t/p/w1280${rawBackdrop}`)
      : null;

    // genres
    let genres = [];
    if (localMovie.genres && localMovie.genres.length > 0) {
      genres = localMovie.genres;
    } else if (tmdbMovie && tmdbMovie.genres) {
      genres = tmdbMovie.genres;
    }

    // wishlist & watched status
    let inWishlist = false;
    let watchedItem = null;

    if (isLoggedIn && tmdbId) {
      try {
        const wishlist = await API.get('/wishlist');
        inWishlist = (wishlist || []).some(w => w.tmdb_id === tmdbId || (w.movie && w.movie.tmdb_id === tmdbId));
      } catch (e) {}

      try {
        const watched = await API.get('/watched');
        watchedItem = (watched || []).find(w => w.tmdb_id === tmdbId || (w.movie && w.movie.tmdb_id === tmdbId));
      } catch (e) {}
    }

    let inWatched = !!watchedItem;
    // strict mutual exclusivity: a movie can stay in only one list at once
    if (inWatched) {
      inWishlist = false;
    }

    // fetch scheduled shows for this local movie
    let shows = [];
    try {
      const allShows = await API.get('/shows');
      shows = (allShows || []).filter(s => s.movie_id === localMovie.id);
    } catch (e) {
      console.error('Failed to load shows:', e);
    }

    await renderMovieDetailsView(container, {
      title,
      overview,
      releaseDate,
      releaseYear,
      runtime,
      language,
      voteAverage,
      poster,
      backdrop,
      genres,
      director,
      cast,
      shows,
      isLoggedIn,
      tmdbId,
      inWishlist,
      inWatched,
      watchedItem,
      isLocal: true,
      localId: localMovie.id
    });

  } catch (err) {
    container.innerHTML = `<div class="card" style="color: var(--brand-primary); text-align: center;">${err.message || 'Local movie not found.'}</div>`;
  }
}

// -------------------------------------------------------------
// 2. tmdb external movie details (#/tmdb/:tmdbId)
// uses external tmdb movie id (tmdb_id)
// -------------------------------------------------------------
async function renderTMDBMovieDetails(container, params) {
  const tmdbId = Number(params.tmdbId);
  const isLoggedIn = !!API.getToken();

  container.innerHTML = '<div style="text-align: center; padding: 3rem; color: var(--text-muted);">Loading TMDB movie details...</div>';

  try {
    // strictly call get /tmdb/movie/:id (external tmdb endpoint)
    const tmdbMovie = await API.get(`/tmdb/movie/${tmdbId}`);
    if (!tmdbMovie || !tmdbMovie.id) {
      throw new Error('TMDB movie details not found');
    }

    // check if this tmdb movie is also published locally (for shows/tickets)
    let localMovie = null;
    try {
      localMovie = await API.get(`/movies/tmdb/${tmdbId}`);
    } catch (e) {}

    // extract director & cast
    let director = localMovie ? localMovie.director : '';
    let cast = localMovie ? localMovie.cast : '';

    if (tmdbMovie.credits) {
      if (!director && tmdbMovie.credits.crew) {
        const dObj = tmdbMovie.credits.crew.find(c => c.job === 'Director');
        if (dObj) director = dObj.name;
      }
      if (!cast && tmdbMovie.credits.cast) {
        cast = tmdbMovie.credits.cast.slice(0, 5).map(c => c.name).join(', ');
      }
    }

    const title = tmdbMovie.title || (localMovie && localMovie.title) || 'Movie Details';
    const overview = tmdbMovie.overview || (localMovie && localMovie.description) || 'No overview available.';
    const releaseDate = tmdbMovie.release_date || (localMovie && localMovie.release_date) || '';
    const releaseYear = releaseDate ? releaseDate.split('-')[0] : '';
    const runtime = tmdbMovie.runtime || (localMovie && localMovie.duration) || 0;
    const language = (tmdbMovie.original_language || (localMovie && localMovie.language) || 'EN').toUpperCase();
    const voteAverage = tmdbMovie.vote_average ? tmdbMovie.vote_average.toFixed(1) : null;

    const rawPoster = tmdbMovie.poster_path || (localMovie && localMovie.poster_path);
    const poster = rawPoster
      ? (rawPoster.startsWith('http') ? rawPoster : `https://image.tmdb.org/t/p/w500${rawPoster}`)
      : 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?q=80&w=1000';

    const rawBackdrop = tmdbMovie.backdrop_path || (localMovie && localMovie.backdrop_path);
    const backdrop = rawBackdrop
      ? (rawBackdrop.startsWith('http') ? rawBackdrop : `https://image.tmdb.org/t/p/w1280${rawBackdrop}`)
      : null;

    const genres = tmdbMovie.genres || (localMovie ? localMovie.genres : []);

    // wishlist & watched status
    let inWishlist = false;
    let watchedItem = null;

    if (isLoggedIn) {
      try {
        const wishlist = await API.get('/wishlist');
        inWishlist = (wishlist || []).some(w => w.tmdb_id === tmdbId || (w.movie && w.movie.tmdb_id === tmdbId));
      } catch (e) {}

      try {
        const watched = await API.get('/watched');
        watchedItem = (watched || []).find(w => w.tmdb_id === tmdbId || (w.movie && w.movie.tmdb_id === tmdbId));
      } catch (e) {}
    }

    let inWatched = !!watchedItem;
    // strict mutual exclusivity: a movie can stay in only one list at once
    if (inWatched) {
      inWishlist = false;
    }

    // scheduled shows if movie is published locally
    let shows = [];
    if (localMovie && localMovie.id) {
      try {
        const allShows = await API.get('/shows');
        shows = (allShows || []).filter(s => s.movie_id === localMovie.id);
      } catch (e) {}
    }

    await renderMovieDetailsView(container, {
      title,
      overview,
      releaseDate,
      releaseYear,
      runtime,
      language,
      voteAverage,
      poster,
      backdrop,
      genres,
      director,
      cast,
      shows,
      isLoggedIn,
      tmdbId,
      inWishlist,
      inWatched,
      watchedItem,
      isLocal: false,
      localId: localMovie ? localMovie.id : null
    });

  } catch (err) {
    container.innerHTML = `<div class="card" style="color: var(--brand-primary); text-align: center;">${err.message || 'TMDB movie details not found.'}</div>`;
  }
}

// helper renderer for both local & tmdb movie details
async function renderMovieDetailsView(container, data) {
  const {
    title, overview, releaseDate, releaseYear, runtime, language, voteAverage,
    poster, backdrop, genres, director, cast, shows, isLoggedIn, tmdbId,
    watchedItem
  } = data;

  let inWishlist = !!data.inWishlist;
  let inWatched = !!data.inWatched;

  container.innerHTML = `
    <div style="display: flex; flex-direction: column; gap: 2rem;">
      
      <!-- backdrop banner -->
      ${backdrop ? `
        <div style="width: 100%; height: 220px; border-radius: 12px; overflow: hidden; position: relative; border: 1px solid var(--border-dark);">
          <img src="${backdrop}" alt="${title}" style="width: 100%; height: 100%; object-fit: cover; filter: brightness(0.4);" />
          <div style="position: absolute; bottom: 1.5rem; left: 1.5rem; color: #ffffff;">
            <h1 style="font-size: 2rem; font-weight: 900; line-height: 1.1; margin-bottom: 0.25rem;">${title}</h1>
            <div style="font-size: 0.8125rem; color: #d1d5db;">
              ${releaseYear ? releaseYear + ' • ' : ''}${runtime ? runtime + ' mins • ' : ''}${language}
              ${voteAverage ? ' • ★ ' + voteAverage + '/10' : ''}
            </div>
          </div>
        </div>
      ` : ''}

      <!-- movie details card -->
      <div class="card" style="display: flex; flex-direction: row; gap: 2rem; flex-wrap: wrap;">
        <img src="${poster}" alt="${title}" style="width: 200px; aspect-ratio: 2/3; object-fit: cover; border-radius: 8px; border: 1px solid var(--border-dark);" />
        
        <div style="flex: 1; display: flex; flex-direction: column; gap: 0.75rem; min-width: 280px;">
          ${!backdrop ? `<h1 style="font-size: 2rem; font-weight: 800; color: #ffffff;">${title}</h1>` : ''}
          
          <div style="font-size: 0.84rem; color: var(--text-muted);">
            ${releaseDate ? 'Release Date: ' + releaseDate + ' • ' : ''}${runtime ? runtime + ' mins • ' : ''}${language}
            ${voteAverage ? ' • ★ ' + voteAverage + '/10' : ''}
          </div>

          ${genres && genres.length > 0 ? `
            <div class="genre-chips">
              ${genres.map(g => `<span class="genre-chip active">${g.name}</span>`).join('')}
            </div>
          ` : ''}

          <p style="font-size: 0.9rem; color: var(--text-main); line-height: 1.6; margin-top: 0.25rem;">
            ${overview}
          </p>
          
          ${director ? `<div style="font-size: 0.84rem; color: #e5e5e5;"><strong>Director:</strong> ${director}</div>` : ''}
          ${cast ? `<div style="font-size: 0.84rem; color: #e5e5e5;"><strong>Cast:</strong> ${cast}</div>` : ''}

          ${isLoggedIn && tmdbId ? `
            <div style="display: flex; flex-direction: column; gap: 1rem; margin-top: 1rem;">
              <div style="display: flex; gap: 0.75rem;">
                <button id="btn-toggle-wishlist" class="btn ${inWishlist ? 'btn-danger' : 'btn-secondary'} btn-sm">
                  ${inWishlist ? '✓ Saved in Wishlist' : '+ Add to Wishlist'}
                </button>
                <button id="btn-toggle-watched" class="btn ${inWatched ? 'btn-danger' : 'btn-secondary'} btn-sm">
                  ${inWatched ? '✓ Marked as Watched' : '+ Mark as Watched'}
                </button>
              </div>

              <!-- review / rating editor for watched movie -->
              ${inWatched ? `
                <div class="card" style="padding: 1rem; background-color: #0a0a0a; border: 1px solid var(--border-dark);">
                  <div style="font-size: 0.84rem; font-weight: 700; color: #ffffff; margin-bottom: 0.5rem;">Your Rating & Review</div>
                  <form id="form-edit-review" style="display: flex; flex-direction: column; gap: 0.5rem;">
                    <div style="display: flex; gap: 0.5rem; align-items: center;">
                      <label style="font-size: 0.75rem; color: var(--text-muted);">Rating (1-10):</label>
                      <input type="number" id="input-user-rating" min="1" max="10" value="${watchedItem && watchedItem.rating ? watchedItem.rating : ''}" class="form-input" style="width: 80px; padding: 0.3rem;" placeholder="Rating" />
                    </div>
                    <textarea id="input-user-review" class="form-textarea" rows="2" style="font-size: 0.8125rem;" placeholder="Write your review...">${watchedItem && watchedItem.review ? watchedItem.review : ''}</textarea>
                    <button type="submit" class="btn btn-primary btn-sm" style="align-self: flex-start;">Save Review</button>
                  </form>
                </div>
              ` : ''}
            </div>
          ` : ''}
        </div>
      </div>

      <!-- scheduled shows section -->
      <div class="card">
        <h2 class="section-title">Scheduled Shows</h2>
        ${shows && shows.length > 0 ? `
          <div class="table-container">
            <table>
              <thead>
                <tr>
                  <th>Screen</th>
                  <th>Showtime</th>
                  <th>Price</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                ${shows.map(s => `
                  <tr>
                    <td>Screen #${s.screen_id}</td>
                    <td>${new Date(s.start_time).toLocaleString([], { dateStyle: 'medium', timeStyle: 'short' })}</td>
                    <td style="color: var(--accent-emerald); font-weight: 700;">Rs. ${s.price}</td>
                    <td>
                      <a href="#/book/${s.id}" class="btn btn-primary btn-sm">Book Seats</a>
                    </td>
                  </tr>
                `).join('')}
              </tbody>
            </table>
          </div>
        ` : '<p style="color: var(--text-muted); font-size: 0.8125rem;">No upcoming showtimes scheduled for this movie.</p>'}
      </div>

      <!-- similar movies shelf -->
      <div id="similar-movies-section" class="card" style="display: none;">
        <h2 class="section-title">Similar Movies</h2>
        <div id="similar-movies-grid" class="grid"></div>
      </div>

    </div>
  `;

  // attach wishlist listener with strict re-render & mutual exclusivity
  const wishlistBtn = document.getElementById('btn-toggle-wishlist');
  if (wishlistBtn && tmdbId) {
    wishlistBtn.addEventListener('click', async () => {
      try {
        if (inWishlist) {
          await API.delete(`/wishlist/${tmdbId}`);
        } else {
          await API.post('/wishlist', { tmdb_id: Number(tmdbId) });
        }
        // re-render controller to update backend state & enforce single list exclusivity
        if (data.isLocal) renderLocalMovieDetails(container, { id: data.localId });
        else renderTMDBMovieDetails(container, { tmdbId });
      } catch (err) { alert(err.message || 'Wishlist update failed'); }
    });
  }

  // attach watched listener with strict re-render & mutual exclusivity
  const watchedBtn = document.getElementById('btn-toggle-watched');
  if (watchedBtn && tmdbId) {
    watchedBtn.addEventListener('click', async () => {
      try {
        if (inWatched) {
          await API.delete(`/watched/${tmdbId}`);
        } else {
          await API.post('/watched', { tmdb_id: Number(tmdbId) });
        }
        // re-render controller to update backend state & enforce single list exclusivity
        if (data.isLocal) renderLocalMovieDetails(container, { id: data.localId });
        else renderTMDBMovieDetails(container, { tmdbId });
      } catch (err) { alert(err.message || 'Watched update failed'); }
    });
  }

  // attach edit review listener
  const reviewForm = document.getElementById('form-edit-review');
  if (reviewForm && tmdbId) {
    reviewForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const ratingVal = document.getElementById('input-user-rating').value;
      const reviewVal = document.getElementById('input-user-review').value;

      const body = {};
      if (ratingVal) body.rating = Number(ratingVal);
      if (reviewVal !== undefined) body.review = reviewVal;

      try {
        await API.put(`/watched/${tmdbId}`, body);
        alert('Rating & Review updated successfully!');
      } catch (err) {
        alert(err.message || 'Failed to update review');
      }
    });
  }

  // fetch similar movies shelf
  if (tmdbId) {
    try {
      const similarSection = document.getElementById('similar-movies-section');
      const similarGrid = document.getElementById('similar-movies-grid');
      const similarRes = await API.post('/recommendations/movie', { movie_ids: [Number(tmdbId)] });

      if (similarRes && similarRes.length > 0) {
        const resolvedSimilar = await resolveRecommendationMovies(similarRes);
        if (resolvedSimilar && resolvedSimilar.length > 0) {
          similarGrid.innerHTML = resolvedSimilar.slice(0, 6).map(renderMovieCard).filter(c => c !== '').join('');
          similarSection.style.display = 'block';
        }
      }
    } catch (e) {}
  }
}

// register two explicit, unambiguous routes
Router.addRoute('#/movies/:id', renderLocalMovieDetails);
Router.addRoute('#/tmdb/:tmdbId', renderTMDBMovieDetails);
