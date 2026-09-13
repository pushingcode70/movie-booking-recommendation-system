/* Theatres & Shows Controller */

async function renderTheatres(container, params) {
  const theatreId = params ? params.id : null;

  if (theatreId) {
    await renderTheatreShowtimes(container, theatreId);
    return;
  }

  container.innerHTML = '<div style="text-align: center; padding: 3rem; color: var(--text-muted);">Loading theatres...</div>';

  try {
    const theatres = await API.get('/theatres') || [];
    const screens = await API.get('/screens') || [];

    container.innerHTML = `
      <div style="max-width: 1100px; margin: 0 auto; padding-top: 1rem;">
        <h1 class="page-title" style="font-size: 2rem; margin-bottom: 0.5rem; color: #ffffff;">Theatres & Multiplexes</h1>
        <p style="color: var(--text-muted); font-size: 0.875rem; margin-bottom: 2rem;">Select a theatre to view available movies and showtimes.</p>

        ${theatres.length > 0 ? `
          <div class="grid" style="grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.5rem;">
            ${theatres.map(t => {
              const tScreens = screens.filter(s => s.theatre_id === t.id);
              return `
                <div class="card theatre-card-item" data-id="${t.id}" style="display: flex; flex-direction: column; justify-content: space-between; background-color: #121212; border: 1px solid var(--border-dark); border-radius: 12px; padding: 1.5rem; cursor: pointer; transition: transform 0.2s, border-color 0.2s;">
                  <div>
                    <h2 style="font-size: 1.25rem; font-weight: 800; color: #ffffff; margin-bottom: 0.35rem;">${t.name}</h2>
                    <div style="font-size: 0.8125rem; color: var(--text-muted); margin-bottom: 1rem;">📍 ${t.location || 'Location not specified'}</div>
                    <div style="font-size: 0.78rem; color: var(--accent-emerald); font-weight: 600;">
                      ${tScreens.length} Screen(s) Available
                    </div>
                  </div>
                  <div style="margin-top: 1.25rem; font-size: 0.8125rem; color: var(--brand-primary); font-weight: 700;">
                    View Shows &rarr;
                  </div>
                </div>
              `;
            }).join('')}
          </div>
        ` : `
          <div class="card" style="text-align: center; padding: 3rem; color: var(--text-muted);">
            No theatres currently registered.
          </div>
        `}
      </div>
    `;

    // Click listener to select theatre
    container.querySelectorAll('.theatre-card-item').forEach(card => {
      card.addEventListener('click', () => {
        const id = card.dataset.id;
        window.location.hash = `#/theatres/${id}`;
      });
    });

  } catch (err) {
    container.innerHTML = `<div class="card" style="color: var(--brand-primary);">${err.message || 'Failed to load theatres.'}</div>`;
  }
}

/* Render Theatre Showtimes Page */
async function renderTheatreShowtimes(container, theatreId) {
  container.innerHTML = '<div style="text-align: center; padding: 3rem; color: var(--text-muted);">Loading theatre showtimes...</div>';

  try {
    const theatre = await API.get(`/theatres/${theatreId}`);
    const shows = await API.get('/shows') || [];
    const screens = await API.get('/screens') || [];
    const movies = await API.get('/movies') || [];

    // Filter screens for this theatre
    const tScreens = screens.filter(s => s.theatre_id === Number(theatreId));
    const tScreenIds = tScreens.map(s => s.id);

    // Filter shows for these screens
    const tShows = shows.filter(s => tScreenIds.includes(s.screen_id));

    // Extract unique date strings YYYY-MM-DD from backend shows
    const backendDateStrings = new Set();
    tShows.forEach(s => {
      if (s.start_time) {
        try {
          const dateStr = new Date(s.start_time).toISOString().split('T')[0];
          backendDateStrings.add(dateStr);
        } catch (e) {}
      }
    });

    // Generate upcoming dates (today + next 3 days)
    const today = new Date();
    for (let i = 0; i < 4; i++) {
      const d = new Date(today);
      d.setDate(today.getDate() + i);
      backendDateStrings.add(d.toISOString().split('T')[0]);
    }

    const sortedDateStrings = Array.from(backendDateStrings).sort();
    const dates = sortedDateStrings.map(dateStr => {
      const d = new Date(dateStr + 'T00:00:00');
      const isTodayStr = new Date().toISOString().split('T')[0] === dateStr;
      const dayName = isTodayStr ? 'TODAY' : d.toLocaleDateString('en-US', { weekday: 'short' }).toUpperCase();
      const dayNum = d.getDate();
      const monthName = d.toLocaleDateString('en-US', { month: 'short' });
      return {
        fullDate: dateStr,
        dayName,
        dayNum,
        monthName
      };
    });

    let selectedDate = dates.length > 0 ? dates[0].fullDate : new Date().toISOString().split('T')[0];

    function renderShowtimesView() {
      // Filter shows by selected date
      const dateShows = tShows.filter(s => {
        if (!s.start_time) return true;
        const sDate = new Date(s.start_time).toISOString().split('T')[0];
        return sDate === selectedDate;
      });

      // Group shows by movie_id
      const movieMap = {};
      dateShows.forEach(s => {
        if (!movieMap[s.movie_id]) movieMap[s.movie_id] = [];
        movieMap[s.movie_id].push(s);
      });

      container.innerHTML = `
        <div style="max-width: 1000px; margin: 0 auto; padding-top: 1rem;">
          
          <a href="#/shows" style="font-size: 0.8125rem; color: var(--text-muted); text-decoration: none; display: inline-block; margin-bottom: 1rem;">
            &larr; Back to All Theatres
          </a>

          <h1 style="font-size: 2rem; font-weight: 800; color: #ffffff; margin-bottom: 0.25rem;">${theatre.name}</h1>
          <div style="font-size: 0.875rem; color: var(--text-muted); margin-bottom: 2rem;">📍 ${theatre.location || 'Location not specified'}</div>

          <!-- Section Title & Date Selector -->
          <div style="margin-bottom: 2rem;">
            <div style="display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.5rem;">
              <span style="color: #ef4444; font-size: 1.25rem;">📅</span>
              <h2 style="font-size: 1.25rem; font-weight: 800; color: #ffffff; margin: 0;">Select Theatre & Showtimes</h2>
            </div>
            <p style="color: var(--text-muted); font-size: 0.8125rem; margin: 0 0 1.25rem 0;">Select a date below to view active show times for this theatre.</p>

            <!-- Date Selector Pills -->
            <div style="display: flex; gap: 0.75rem; flex-wrap: wrap;">
              ${dates.map(d => `
                <button class="date-pill-btn ${d.fullDate === selectedDate ? 'active-red-pill' : ''}" data-date="${d.fullDate}" style="
                  background-color: ${d.fullDate === selectedDate ? '#ef4444' : '#161616'};
                  color: #ffffff;
                  border: 1px solid ${d.fullDate === selectedDate ? '#ef4444' : 'var(--border-dark)'};
                  border-radius: 10px;
                  padding: 0.6rem 1.25rem;
                  display: flex;
                  flex-direction: column;
                  align-items: center;
                  cursor: pointer;
                  min-width: 75px;
                  transition: background-color 0.2s, border-color 0.2s;
                ">
                  <span style="font-size: 0.65rem; font-weight: 800; letter-spacing: 0.5px; opacity: 0.9;">${d.dayName}</span>
                  <span style="font-size: 1.15rem; font-weight: 900; line-height: 1.2; margin: 0.1rem 0;">${d.dayNum}</span>
                  <span style="font-size: 0.65rem; font-weight: 600; opacity: 0.8;">${d.monthName}</span>
                </button>
              `).join('')}
            </div>
          </div>

          <!-- Movies & Showtimes Cards List -->
          <div style="display: flex; flex-direction: column; gap: 1.5rem;">
            ${Object.keys(movieMap).length > 0 ? Object.keys(movieMap).map(mId => {
              const mShows = movieMap[mId];
              const sampleShow = mShows[0];
              let movie = (sampleShow && sampleShow.movie && sampleShow.movie.title) ? sampleShow.movie : null;
              if (!movie) {
                movie = movies.find(m => m.id === Number(mId) || m.tmdb_id === Number(mId));
              }
              if (!movie || !movie.title) {
                movie = { title: `Movie #${mId}`, poster_path: '' };
              }

              const poster = movie.poster_path
                ? (movie.poster_path.startsWith('http') ? movie.poster_path : `https://image.tmdb.org/t/p/w500${movie.poster_path}`)
                : 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?q=80&w=1000';

              return `
                <div class="card" style="display: flex; gap: 1.5rem; background-color: #121212; border: 1px solid var(--border-dark); border-radius: 12px; padding: 1.25rem; flex-wrap: wrap;">
                  <!-- Movie Poster -->
                  <img src="${poster}" alt="${movie.title}" style="width: 100px; aspect-ratio: 2/3; object-fit: cover; border-radius: 8px; border: 1px solid var(--border-dark);" />

                  <!-- Movie Info & Showtimes -->
                  <div style="flex: 1; min-width: 250px;">
                    <h3 style="font-size: 1.2rem; font-weight: 800; color: #ffffff; margin-bottom: 0.25rem;">${movie.title}</h3>
                    <div style="font-size: 0.78rem; color: var(--text-muted); margin-bottom: 1rem;">
                      ${movie.duration ? movie.duration + ' mins • ' : ''}${movie.language || 'English'}
                    </div>

                    <div style="font-size: 0.75rem; font-weight: 700; color: #e5e5e5; margin-bottom: 0.6rem;">AVAILABLE SHOWTIMES</div>
                    
                    <!-- Showtimes Pills Grid -->
                    <div style="display: flex; gap: 0.75rem; flex-wrap: wrap;">
                      ${mShows.map(s => {
                        const scr = tScreens.find(sc => sc.id === s.screen_id);
                        const sName = scr ? (scr.screen_number ? `Screen #${scr.screen_number}` : scr.name || `Screen #${scr.id}`) : `Screen #${s.screen_id}`;
                        const timeStr = s.start_time ? new Date(s.start_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : 'Showtime';

                        return `
                          <a href="#/book/${s.id}" style="
                            background-color: #1a1a1a;
                            border: 1px solid var(--accent-emerald);
                            color: #ffffff;
                            padding: 0.5rem 0.85rem;
                            border-radius: 8px;
                            text-decoration: none;
                            display: flex;
                            flex-direction: column;
                            align-items: center;
                            transition: transform 0.15s, background-color 0.15s;
                          ">
                            <span style="font-size: 0.85rem; font-weight: 800; color: var(--accent-emerald);">${timeStr}</span>
                            <span style="font-size: 0.7rem; color: var(--text-muted); margin-top: 0.15rem;">${sName} • Rs.${s.price}</span>
                          </a>
                        `;
                      }).join('')}
                    </div>
                  </div>
                </div>
              `;
            }).join('') : `
              <div class="card" style="text-align: center; padding: 3rem; color: var(--text-muted); background-color: #121212; border: 1px dashed var(--border-dark); border-radius: 12px;">
                No showtimes scheduled for this theatre on ${selectedDate}.
              </div>
            `}
          </div>

        </div>
      `;

      // Date pill click listener
      container.querySelectorAll('.date-pill-btn').forEach(btn => {
        btn.addEventListener('click', () => {
          selectedDate = btn.dataset.date;
          renderShowtimesView();
        });
      });
    }

    renderShowtimesView();

  } catch (err) {
    container.innerHTML = `<div class="card" style="color: var(--brand-primary);">${err.message || 'Failed to load theatre showtimes.'}</div>`;
  }
}

Router.addRoute('#/shows', renderTheatres);
Router.addRoute('#/theatres', renderTheatres);
Router.addRoute('#/shows/:id', renderTheatres);
Router.addRoute('#/theatres/:id', renderTheatres);
