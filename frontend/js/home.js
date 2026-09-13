/* home & recommendation controller with single favorite genre shelf, 30-movie limit, >= 90 min filter, and clean headings */

function renderMovieCard(movie) {
  if (!movie) return '';

  const duration = movie.duration || movie.runtime || 0;
  if (duration > 0 && duration < 90) {
    return ''; // strictly do not show any movie under 90 minutes
  }

  const poster = movie.poster_path
    ? (movie.poster_path.startsWith('http') ? movie.poster_path : `https://image.tmdb.org/t/p/w500${movie.poster_path}`)
    : 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?q=80&w=1000';

  const durationStr = duration ? `${duration} mins` : '';
  const langStr = movie.language ? movie.language.toUpperCase() : (movie.original_language ? movie.original_language.toUpperCase() : 'EN');
  
  // explicit route disambiguation:
  // local published movie -> is_local === true -> #/movies/<local_id>
  // external tmdb movie   -> is_local === false -> #/tmdb/<tmdb_id>
  const tmdbId = movie.tmdb_id || movie.TMDBID || movie.id;
  const route = (movie.is_local && movie.id) ? `#/movies/${movie.id}` : `#/tmdb/${tmdbId}`;

  return `
    <a href="${route}" class="movie-card">
      <img src="${poster}" alt="${movie.title}" class="movie-poster" loading="lazy" />
      <div class="movie-info">
        <div class="movie-title">${movie.title || 'Untitled Movie'}</div>
        <div class="movie-meta">
          ${langStr} ${durationStr ? '• ' + durationStr : ''}
        </div>
      </div>
    </a>
  `;
}

// resolves raw recommendation items to full movie data objects with explicit is_local flags & 90+ min filter
async function resolveRecommendationMovies(rawItems) {
  if (!rawItems || rawItems.length === 0) return [];

  // resolve each raw item containing tmdb_id or movie object
  const resolved = await Promise.all(
    rawItems.slice(0, 40).map(async (item) => {
      const tmdbId = item.tmdb_id || item.TMDBID || (item.movie && (item.movie.tmdb_id || item.movie.TMDBID)) || item.id;
      if (!tmdbId) return null;

      // 1. check if movie is published in local database
      try {
        const localMovie = await API.get(`/movies/tmdb/${tmdbId}`);
        if (localMovie && localMovie.id) {
          const dur = localMovie.duration || localMovie.runtime || 0;
          if (dur > 0 && dur < 90) return null; // filter movies < 90 mins
          return { ...localMovie, is_local: true };
        }
      } catch (e) {}

      // 2. fallback to tmdb movie details endpoint (external catalog item)
      try {
        const tmdbMovie = await API.get(`/tmdb/movie/${tmdbId}`);
        if (tmdbMovie && (tmdbMovie.title || tmdbMovie.name)) {
          const dur = tmdbMovie.runtime || tmdbMovie.duration || 0;
          if (dur > 0 && dur < 90) return null; // filter movies < 90 mins

          return {
            is_local: false,
            tmdb_id: tmdbMovie.id,
            title: tmdbMovie.title || tmdbMovie.name,
            overview: tmdbMovie.overview,
            poster_path: tmdbMovie.poster_path,
            backdrop_path: tmdbMovie.backdrop_path,
            duration: tmdbMovie.runtime,
            language: tmdbMovie.original_language || 'EN',
            release_date: tmdbMovie.release_date
          };
        }
      } catch (e) {}

      return null;
    })
  );

  return resolved.filter(m => m !== null && ((m.duration || m.runtime || 0) >= 90 || !(m.duration || m.runtime)));
}

async function renderHome(container) {
  const isLoggedIn = !!API.getToken();

  container.innerHTML = `
    <div style="display: flex; flex-direction: column; gap: 2.5rem;">
      <!-- search card -->
      <div class="search-card">
        <form id="form-search-recommendation">
          <div class="search-input-wrapper">
            <input type="text" id="rec-prompt" class="search-input" placeholder="Search movies, themes, or custom prompt (e.g. dark detective movies)..." />
            <button type="submit" class="btn btn-primary" style="padding: 0.75rem 1.5rem;">Find Movies</button>
          </div>

          <!-- genre chips list -->
          <div id="rec-genre-chips" class="genre-chips">
            <div style="font-size: 0.75rem; color: var(--text-muted);">Loading genres...</div>
          </div>
        </form>
      </div>

      <!-- main recommendation / search results section (up to 30 movies) -->
      <div>
        <h2 id="results-heading" class="section-title">Recommended For You</h2>
        <div id="home-movies-grid" class="grid">
          <div style="font-size: 0.8125rem; color: var(--text-muted);">Loading movies...</div>
        </div>
      </div>

      ${isLoggedIn ? `
        <!-- single favorite genres shelf (up to 30 movies) -->
        <div>
          <h2 class="section-title">Based on Genres You Like</h2>
          <div id="shelf-fav-genres" class="grid">
            <div style="font-size: 0.8125rem; color: var(--text-muted);">Loading recommendations based on your favorite genres...</div>
          </div>
        </div>
      ` : ''}
    </div>
  `;

  // fetch genres for chips
  let selectedGenreIds = [];
  try {
    const genres = await API.get('/genres');
    const genreContainer = document.getElementById('rec-genre-chips');
    if (genres && genres.length > 0) {
      genreContainer.innerHTML = genres.map(g => {
        const tmdbId = g.tmdb_id || g.id;
        return `<div class="genre-chip" data-tmdb-id="${tmdbId}">${g.name}</div>`;
      }).join('');

      genreContainer.addEventListener('click', (e) => {
        const chip = e.target.closest('.genre-chip');
        if (!chip) return;
        const tmdbId = Number(chip.dataset.tmdbId);
        if (selectedGenreIds.includes(tmdbId)) {
          selectedGenreIds = selectedGenreIds.filter(id => id !== tmdbId);
          chip.classList.remove('active');
        } else {
          selectedGenreIds.push(tmdbId);
          chip.classList.add('active');
        }
      });
    }
  } catch (err) {
    console.error('Failed to load genres:', err);
  }

  // load initial recommended / published movies
  const gridEl = document.getElementById('home-movies-grid');
  try {
    const movies = await API.get('/movies');
    if (movies && movies.length > 0) {
      const localMovies = movies
        .filter(m => !m.duration || m.duration >= 90)
        .map(m => ({ ...m, is_local: true }));

      gridEl.innerHTML = localMovies.slice(0, 30).map(renderMovieCard).filter(c => c !== '').join('');
    } else {
      gridEl.innerHTML = '<div style="font-size: 0.8125rem; color: var(--text-muted);">No movies currently published.</div>';
    }
  } catch (err) {
    gridEl.innerHTML = `<div style="font-size: 0.8125rem; color: var(--brand-primary);">Failed to load movies.</div>`;
  }

  // handle search / recommendation form submit
  document.getElementById('form-search-recommendation').addEventListener('submit', async (e) => {
    e.preventDefault();
    const prompt = document.getElementById('rec-prompt').value.trim();
    const headingEl = document.getElementById('results-heading');
    
    // clean section heading without repeating long prompt text
    headingEl.textContent = prompt ? 'Search Results' : 'Recommended For You';
    gridEl.innerHTML = '<div style="font-size: 0.8125rem; color: var(--text-muted);">Searching recommendations...</div>';

    try {
      let rawMovies = [];
      if (prompt && selectedGenreIds.length > 0) {
        rawMovies = await API.post('/recommendations/prompt-genre', { prompt, genre_ids: selectedGenreIds });
      } else if (prompt) {
        rawMovies = await API.post('/recommendations/custom', { prompt });
      } else if (selectedGenreIds.length > 0) {
        rawMovies = await API.post('/recommendations/genre', { genre_ids: selectedGenreIds });
      } else {
        const pubMovies = await API.get('/movies');
        rawMovies = (pubMovies || []).map(m => ({ ...m, is_local: true }));
      }

      // resolve metadata, explicit is_local flag, and filter out movies < 90 minutes
      const resolvedMovies = await resolveRecommendationMovies(rawMovies);

      if (resolvedMovies && resolvedMovies.length > 0) {
        const renderedCards = resolvedMovies.map(renderMovieCard).filter(c => c !== '');
        gridEl.innerHTML = renderedCards.slice(0, 30).join('');
      } else {
        gridEl.innerHTML = '<div style="font-size: 0.8125rem; color: var(--text-muted);">No matching movies found.</div>';
      }
    } catch (err) {
      let errMsg = err.message || 'Search failed.';
      if (errMsg.includes('8001') || errMsg.includes('connection refused') || errMsg.includes('embed')) {
        errMsg = 'The prompt vector embedding service (port 8001) is currently offline. Please start Python embedding service or select genre chips to discover movies.';
      }
      gridEl.innerHTML = `<div style="color: var(--brand-primary); font-size: 0.8125rem; background-color: #1a080a; padding: 1rem; border-radius: 8px; border: 1px solid #4a1518;">${errMsg}</div>`;
    }
  });

  // load single favorite genres recommendation shelf for logged-in user
  if (isLoggedIn) {
    const favShelfEl = document.getElementById('shelf-fav-genres');
    try {
      const userFavGenres = await API.get('/users/me/genres');
      if (userFavGenres && userFavGenres.length > 0) {
        const genreIds = userFavGenres.map(g => g.tmdb_id || g.TMDBID || g.id);
        const rawRecMovies = await API.post('/recommendations/genre', { genre_ids: genreIds });
        const resolvedRecMovies = await resolveRecommendationMovies(rawRecMovies);
        if (resolvedRecMovies && resolvedRecMovies.length > 0) {
          const renderedFavCards = resolvedRecMovies.map(renderMovieCard).filter(c => c !== '');
          favShelfEl.innerHTML = renderedFavCards.slice(0, 30).join('');
        } else {
          favShelfEl.innerHTML = '<div style="font-size: 0.8125rem; color: var(--text-muted);">No recommendations found for your selected genres.</div>';
        }
      } else {
        favShelfEl.innerHTML = '<div style="font-size: 0.8125rem; color: var(--text-muted);">Select your favorite genres in Profile to personalize recommendations here.</div>';
      }
    } catch (err) {
      if (favShelfEl) favShelfEl.innerHTML = '<div style="font-size: 0.8125rem; color: var(--text-muted);">Could not load genre recommendations.</div>';
    }
  }
}

Router.addRoute('#/', renderHome);
