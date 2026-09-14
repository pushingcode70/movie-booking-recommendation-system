/* admin dashboard & management controller */

let currentAdminTab = 'theatres';

async function renderAdmin(container) {
  const user = API.getUser();
  if (!API.getToken() || !user || (user.role !== 'admin' && user.role !== 'ADMIN')) {
    container.innerHTML = `
      <div class="card" style="text-align: center; padding: 3rem; max-width: 500px; margin: 2rem auto;">
        <h1 style="color: var(--accent); font-size: 1.5rem; font-weight: 800;">Access Denied</h1>
        <p style="color: var(--muted); margin-top: 0.5rem;">You must be logged in as an Administrator to access this section.</p>
        <a href="#/" class="btn btn-secondary" style="margin-top: 1rem;">Return to Main Site</a>
      </div>
    `;
    return;
  }

  container.innerHTML = `
    <div>
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem;">
        <h1 class="page-title" style="margin-bottom: 0;">Admin Management</h1>
        <a href="#/" class="btn btn-secondary btn-sm">&larr; Back to Main Site</a>
      </div>

      <!-- compact top navigation tabs -->
      <div class="admin-nav-bar">
        <div id="admin-tab-theatres" class="admin-nav-tab ${currentAdminTab === 'theatres' ? 'active' : ''}">
          Theatres & Screens
        </div>
        <div id="admin-tab-shows" class="admin-nav-tab ${currentAdminTab === 'shows' ? 'active' : ''}">
          Shows
        </div>
        <div id="admin-tab-bookings" class="admin-nav-tab ${currentAdminTab === 'bookings' || currentAdminTab === 'payments' ? 'active' : ''}">
          Bookings & Payments
        </div>
        <div id="admin-tab-dashboard" class="admin-nav-tab ${currentAdminTab === 'dashboard' ? 'active' : ''}">
          Dashboard
        </div>
      </div>

      <!-- main content view -->
      <main id="admin-content-view">
        <div style="font-size: 0.8125rem; color: var(--muted);">Loading view...</div>
      </main>
    </div>
  `;

  // attach tab click listeners
  ['theatres', 'shows', 'bookings', 'dashboard'].forEach(tab => {
    const el = document.getElementById(`admin-tab-${tab}`);
    if (el) {
      el.addEventListener('click', () => {
        document.querySelectorAll('.admin-nav-tab').forEach(i => i.classList.remove('active'));
        el.classList.add('active');
        currentAdminTab = tab;
        renderAdminTabContent(document.getElementById('admin-content-view'));
      });
    }
  });

  renderAdminTabContent(document.getElementById('admin-content-view'));
}

async function renderAdminTabContent(viewEl) {
  if (currentAdminTab === 'theatres') {
    renderAdminTheatresView(viewEl);
  } else if (currentAdminTab === 'shows') {
    renderAdminShowsView(viewEl);
  } else if (currentAdminTab === 'bookings' || currentAdminTab === 'payments') {
    renderAdminBookingsView(viewEl);
  } else if (currentAdminTab === 'dashboard') {
    renderAdminDashboardView(viewEl);
  }
}

/* 1. combined theatre & screen management view */
async function renderAdminTheatresView(viewEl) {
  viewEl.innerHTML = `
    <div>
      <!-- top header row -->
      <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem; flex-wrap: wrap; gap: 1rem;">
        <div>
          <h1 class="page-title" style="font-size: 2rem; margin-bottom: 0.25rem; color: #ffffff;">Theatre Management</h1>
          <p style="color: var(--text-muted); font-size: 0.875rem;">
            Configure multiple branches, geographical locations, and facilities.
          </p>
        </div>
        <button id="btn-toggle-add-theatre" class="btn btn-primary" style="display: flex; align-items: center; gap: 0.5rem;">
          <span>+ Add Theatre</span>
        </button>
      </div>

      <!-- create theatre card -->
      <div id="card-add-theatre" class="card" style="display: none; max-width: 600px; margin-bottom: 2rem; border: 1px solid var(--accent-emerald);">
        <h2 class="section-title" style="font-size: 1.1rem; color: #ffffff; margin-bottom: 1rem;">Create New Theatre</h2>
        <form id="admin-form-add-theatre">
          <div style="display: flex; gap: 1rem; flex-wrap: wrap;">
            <div class="form-group" style="flex: 1; min-width: 200px; margin-bottom: 0;">
              <label class="form-label" for="t-name">Theatre Name</label>
              <input type="text" id="t-name" class="form-input" required placeholder="e.g. PVR Elante" />
            </div>
            <div class="form-group" style="flex: 1; min-width: 200px; margin-bottom: 0;">
              <label class="form-label" for="t-loc">Location</label>
              <input type="text" id="t-loc" class="form-input" required placeholder="e.g. Elante Mall, Industrial Area Phase 1" />
            </div>
          </div>
          <div style="display: flex; gap: 0.75rem; justify-content: flex-end; margin-top: 1rem;">
            <button type="button" id="btn-cancel-add-theatre" class="btn btn-secondary btn-sm">Cancel</button>
            <button type="submit" class="btn btn-primary btn-sm">Save Theatre</button>
          </div>
        </form>
      </div>

      <!-- search & filter bar -->
      <div style="margin-bottom: 1.5rem;">
        <input type="text" id="t-search-input" class="form-input" placeholder="🔍 Search theatres by branch or city..." style="background-color: #121212; max-width: 500px;" />
      </div>

      <!-- existing theatres cards grid -->
      <div id="admin-theatres-grid" class="grid" style="grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 1.5rem;">
        <div style="font-size: 0.8125rem; color: var(--text-muted);">Loading theatres...</div>
      </div>
    </div>
  `;

  // toggle create theatre form
  const addTheatreCard = document.getElementById('card-add-theatre');
  document.getElementById('btn-toggle-add-theatre').addEventListener('click', () => {
    addTheatreCard.style.display = addTheatreCard.style.display === 'none' ? 'block' : 'none';
  });
  document.getElementById('btn-cancel-add-theatre').addEventListener('click', () => {
    addTheatreCard.style.display = 'none';
  });

  // handle submit new theatre
  document.getElementById('admin-form-add-theatre').addEventListener('submit', async (e) => {
    e.preventDefault();
    const name = document.getElementById('t-name').value.trim();
    const location = document.getElementById('t-loc').value.trim();
    try {
      await API.post('/theatres', { name, location });
      alert('Theatre created successfully!');
      renderAdminTheatresView(viewEl);
    } catch (err) { alert(err.message || 'Failed to create theatre'); }
  });

  // fetch theatres & screens data
  const gridEl = document.getElementById('admin-theatres-grid');
  let allTheatres = [];
  let allScreens = [];

  function renderTheatresList(filterQuery = '') {
    const q = filterQuery.toLowerCase();
    const filtered = allTheatres.filter(t => 
      (t.name || '').toLowerCase().includes(q) || (t.location || '').toLowerCase().includes(q)
    );

    if (filtered.length > 0) {
      gridEl.innerHTML = filtered.map(t => {
        const tScreens = allScreens.filter(s => s.theatre_id === t.id);

        return `
          <div class="card" style="display: flex; flex-direction: column; justify-content: space-between; background-color: #121212; border: 1px solid var(--border); border-radius: 12px; padding: 1.25rem;">
            
            <div>
              <!-- card header -->
              <div style="display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.75rem;">
                <div style="background-color: rgba(239,68,68,0.15); color: #ef4444; width: 36px; height: 36px; border-radius: 8px; display: flex; align-items: center; justify-content: center; font-size: 1.1rem; flex-shrink: 0;">
                  🍿
                </div>
                <div>
                  <h3 style="font-size: 1.15rem; font-weight: 800; color: #ffffff; line-height: 1.2;">${t.name}</h3>
                  <div style="font-size: 0.78rem; color: var(--text-muted); margin-top: 0.2rem; display: flex; align-items: center; gap: 0.3rem;">
                    <span>📍 ${t.location}</span>
                  </div>
                </div>
              </div>

              <!-- screens list section -->
              <div style="background-color: #0a0a0a; border-radius: 8px; padding: 0.85rem; border: 1px solid var(--border); margin-top: 1rem; margin-bottom: 1rem;">
                <div style="font-size: 0.78rem; font-weight: 700; color: #e5e5e5; margin-bottom: 0.5rem; display: flex; justify-content: space-between;">
                  <span>Screens (${tScreens.length})</span>
                </div>

                ${tScreens.length > 0 ? `
                  <div style="display: flex; flex-direction: column; gap: 0.4rem;">
                    ${tScreens.map(s => {
                      const sName = s.screen_number ? `Screen #${s.screen_number}` : (s.name || `Screen #${s.id}`);
                      const seats = s.total_seats || s.capacity || 60;
                      return `
                        <div style="display: flex; justify-content: space-between; align-items: center; background-color: #161616; padding: 0.4rem 0.6rem; border-radius: 6px; font-size: 0.78rem;">
                          <span style="color: #ffffff; font-weight: 600;">${sName} <span style="color: var(--text-muted); font-weight: 400;">(${seats} seats)</span></span>
                          <button class="btn btn-danger btn-sm btn-delete-screen" data-id="${s.id}" style="padding: 0.15rem 0.4rem; font-size: 0.7rem;">
                            ✕ Remove
                          </button>
                        </div>
                      `;
                    }).join('')}
                  </div>
                ` : '<div style="font-size: 0.75rem; color: var(--text-muted);">No screens added yet.</div>'}

                <!-- inline add screen toggle -->
                <details style="margin-top: 0.75rem; font-size: 0.75rem;">
                  <summary style="cursor: pointer; color: var(--accent); font-weight: 600;">+ Add Screen</summary>
                  <form class="form-add-screen-inline" data-theatre-id="${t.id}" style="display: flex; flex-direction: column; gap: 0.5rem; margin-top: 0.5rem; background-color: #1a1a1a; padding: 0.6rem; border-radius: 6px;">
                    <input type="number" class="form-input input-screen-number" required min="1" placeholder="Screen Number (e.g. 1)" style="font-size: 0.75rem; padding: 0.35rem 0.5rem;" />
                    <input type="number" class="form-input input-screen-seats" required min="1" placeholder="Total Seats (e.g. 60)" style="font-size: 0.75rem; padding: 0.35rem 0.5rem;" />
                    <button type="submit" class="btn btn-primary btn-sm" style="align-self: flex-end; font-size: 0.75rem; padding: 0.3rem 0.75rem;">Add Screen</button>
                  </form>
                </details>

              </div>
            </div>

            <!-- card bottom action buttons -->
            <div style="display: flex; gap: 0.5rem; margin-top: 0.5rem;">
              <button class="btn btn-secondary btn-sm btn-edit-theatre" data-id="${t.id}" data-name="${t.name}" data-location="${t.location}" style="flex: 1; text-align: center; justify-content: center;">
                ✏ Edit
              </button>
              <button class="btn btn-danger btn-sm btn-delete-theatre" data-id="${t.id}" style="flex: 1; text-align: center; justify-content: center;">
                🗑 Delete
              </button>
            </div>

          </div>
        `;
      }).join('');
    } else {
      gridEl.innerHTML = '<div class="card" style="text-align: center; color: var(--text-muted); padding: 2rem; grid-column: 1/-1;">No matching theatres found.</div>';
    }
  }

  try {
    allTheatres = (await API.get('/theatres')) || [];
    allScreens = (await API.get('/screens')) || [];
    renderTheatresList();
  } catch (err) {
    gridEl.innerHTML = `<div class="card" style="color: var(--accent); grid-column: 1/-1;">${err.message || 'Failed to load theatres.'}</div>`;
  }

  // live search filter listener
  document.getElementById('t-search-input').addEventListener('input', (e) => {
    renderTheatresList(e.target.value.trim());
  });

  // event delegation for delete theatre, edit theatre, and add/delete screen
  gridEl.addEventListener('click', async (e) => {
    const delTheatreBtn = e.target.closest('.btn-delete-theatre');
    const editTheatreBtn = e.target.closest('.btn-edit-theatre');
    const delScreenBtn = e.target.closest('.btn-delete-screen');

    if (delTheatreBtn) {
      const id = delTheatreBtn.dataset.id;
      if (confirm(`Delete theatre #${id} and associated screens?`)) {
        try {
          await API.delete(`/theatres/${id}`);
          renderAdminTheatresView(viewEl);
        } catch (err) { alert(err.message || 'Failed to delete theatre'); }
      }
    }

    if (editTheatreBtn) {
      const id = editTheatreBtn.dataset.id;
      const oldName = editTheatreBtn.dataset.name;
      const oldLoc = editTheatreBtn.dataset.location;

      const newName = prompt('Update Theatre Name:', oldName);
      if (newName === null) return;
      const newLoc = prompt('Update Theatre Location:', oldLoc);
      if (newLoc === null) return;

      try {
        await API.put(`/theatres/${id}`, { name: newName.trim(), location: newLoc.trim() });
        alert('Theatre updated!');
        renderAdminTheatresView(viewEl);
      } catch (err) { alert(err.message || 'Failed to update theatre'); }
    }

    if (delScreenBtn) {
      const id = delScreenBtn.dataset.id;
      if (confirm(`Remove Screen #${id}?`)) {
        try {
          await API.delete(`/screens/${id}`);
          renderAdminTheatresView(viewEl);
        } catch (err) {
          const msg = (err.message || '').toLowerCase();
          if (msg.includes('fk_shows_screen') || msg.includes('foreign key') || msg.includes('shows')) {
            alert(`Cannot remove Screen #${id} because shows are scheduled on it.\n\nPlease go to the Shows tab and cancel the scheduled shows for this screen first.`);
          } else {
            alert(err.message || 'Failed to remove screen');
          }
        }
      }
    }
  });

  // listener for inline add screen form submissions
  gridEl.addEventListener('submit', async (e) => {
    const form = e.target.closest('.form-add-screen-inline');
    if (!form) return;
    e.preventDefault();

    const theatre_id = Number(form.dataset.theatreId);
    const screen_number = Number(form.querySelector('.input-screen-number').value);
    const total_seats = Number(form.querySelector('.input-screen-seats').value);

    try {
      const createdScreen = await API.post('/screens', { theatre_id, screen_number, total_seats });
      
      // auto-generate seat layout matching total_seats in database
      const seatsPerRow = total_seats > 50 ? 15 : 10;
      const rows = Math.ceil(total_seats / seatsPerRow);
      try {
        await API.post(`/screens/${createdScreen.id}/seats/generate`, {
          rows: rows,
          seats_per_row: seatsPerRow
        });
      } catch (seatErr) {}

      alert(`Screen #${screen_number} (${total_seats} seats) added successfully!`);
      renderAdminTheatresView(viewEl);
    } catch (err) { alert(err.message || 'Failed to add screen'); }
  });
}

/* 2. shows tab view: search & auto-import from tmdb api */
async function renderAdminShowsView(viewEl) {
  let selectedMovie = null;

  viewEl.innerHTML = `
    <div>
      <h1 class="page-title" style="font-size: 2rem; margin-bottom: 0.25rem; color: #ffffff;">Manage & Schedule Shows</h1>
      <p style="color: var(--text-muted); font-size: 0.875rem; margin-bottom: 1.5rem;">
        Search any movie from TMDB API, auto-import to local database, and assign showtimes.
      </p>

      <!-- schedule show card -->
      <div class="card" style="max-width: 600px; margin-bottom: 2rem;">
        <h2 class="section-title">Schedule New Show</h2>
        
        <form id="admin-form-add-show" style="display: flex; flex-direction: column; gap: 1.25rem;">
          
          <!-- movie selection with interactive tmdb api search -->
          <div class="form-group" style="margin-bottom: 0; position: relative;">
            <label class="form-label">Search & Select Movie (TMDB Import Enabled)</label>
            
            <!-- selected movie preview display -->
            <div id="sh-selected-movie-card" style="display: none; background-color: #0a0a0a; border: 1px solid var(--accent-emerald); border-radius: 8px; padding: 0.75rem; align-items: center; gap: 1rem;">
              <img id="sh-selected-poster" src="" alt="Poster" style="width: 45px; aspect-ratio: 2/3; object-fit: cover; border-radius: 4px; border: 1px solid var(--border);" />
              <div style="flex: 1;">
                <div id="sh-selected-title" style="font-size: 0.9rem; font-weight: 800; color: #ffffff;">Movie Title</div>
                <div id="sh-selected-sub" style="font-size: 0.75rem; color: var(--accent-emerald);">Published Local ID: #1 • 120 mins</div>
              </div>
              <button type="button" id="btn-remove-selected-movie" class="btn btn-danger btn-sm">✕ Remove</button>
            </div>

            <!-- search bar input -->
            <div id="sh-search-container">
              <input type="text" id="sh-movie-search" class="form-input" placeholder="Type movie title to search TMDB catalog (e.g. Inception, Avatar)..." autocomplete="off" />
              <div id="sh-import-status" style="display: none; font-size: 0.75rem; color: var(--accent-emerald); margin-top: 0.25rem;">
                ⌛ Importing movie metadata from TMDB into local database...
              </div>

              <!-- dropdown search results box -->
              <div id="sh-movie-results" style="display: none; position: absolute; left: 0; right: 0; top: 100%; z-index: 100; background-color: #121212; border: 1px solid var(--border); border-radius: 8px; max-height: 300px; overflow-y: auto; box-shadow: 0 10px 25px rgba(0,0,0,0.9); margin-top: 4px;">
              </div>
            </div>
          </div>

          <!-- screen selection dropdown -->
          <div class="form-group" style="margin-bottom: 0;">
            <label class="form-label" for="sh-s-id">Select Screen</label>
            <select id="sh-s-id" class="form-input" required>
              <option value="">Loading screens...</option>
            </select>
          </div>

          <!-- start time picker -->
          <div class="form-group" style="margin-bottom: 0;">
            <label class="form-label" for="sh-time">Show Start Time</label>
            <input type="datetime-local" id="sh-time" class="form-input" required />
          </div>

          <!-- ticket price input -->
          <div class="form-group" style="margin-bottom: 0;">
            <label class="form-label" for="sh-price">Ticket Price (Rs.)</label>
            <input type="number" step="0.01" id="sh-price" class="form-input" required placeholder="250.00" />
          </div>

          <button type="submit" class="btn btn-primary" style="width: 100%; margin-top: 0.5rem;">Schedule Show</button>
        </form>
      </div>

      <!-- scheduled shows table -->
      <div>
        <h2 class="section-title">All Scheduled Shows</h2>
        <div id="admin-shows-list">
          <div style="font-size: 0.8125rem; color: var(--text-muted);">Loading scheduled shows...</div>
        </div>
      </div>

    </div>
  `;

  // populate screens dropdown with correct field names
  let screens = [];
  try {
    screens = (await API.get('/screens')) || [];
    const theatres = (await API.get('/theatres')) || [];
    const screenSelect = document.getElementById('sh-s-id');

    if (screens.length > 0) {
      screenSelect.innerHTML = '<option value="">-- Choose Screen --</option>' + screens.map(s => {
        const t = theatres.find(th => th.id === s.theatre_id);
        const tName = t ? t.name : `Theatre #${s.theatre_id}`;
        const sName = s.screen_number ? `Screen #${s.screen_number}` : (s.name || `Screen #${s.id}`);
        const seats = s.total_seats || s.capacity || 60;
        return `<option value="${s.id}">${tName} — ${sName} (${seats} Seats)</option>`;
      }).join('');
    } else {
      screenSelect.innerHTML = '<option value="">No screens available (Add theatre/screen first)</option>';
    }
  } catch (e) {}

  // interactive tmdb api search handler
  const searchInput = document.getElementById('sh-movie-search');
  const resultsBox = document.getElementById('sh-movie-results');
  const searchContainer = document.getElementById('sh-search-container');
  const selectedCard = document.getElementById('sh-selected-movie-card');
  const removeMovieBtn = document.getElementById('btn-remove-selected-movie');
  const importStatus = document.getElementById('sh-import-status');

  let debounceTimer = null;

  async function performTMDBSearch(query) {
    if (!query) {
      resultsBox.style.display = 'none';
      return;
    }

    resultsBox.innerHTML = '<div style="padding: 0.8rem; font-size: 0.8125rem; color: var(--text-muted);">Searching TMDB catalog...</div>';
    resultsBox.style.display = 'block';

    try {
      // 1. fetch local published movies
      const localMovies = (await API.get('/movies')) || [];
      const matchingLocal = localMovies.filter(m => (m.title || '').toLowerCase().includes(query.toLowerCase()));

      // 2. fetch tmdb api search results
      let tmdbResults = [];
      try {
        const res = await API.get(`/tmdb/search?query=${encodeURIComponent(query)}`);
        tmdbResults = res.results || res || [];
      } catch (e) {}

      let html = '';

      // published movies section
      if (matchingLocal.length > 0) {
        html += `<div style="font-size: 0.72rem; font-weight: 700; color: var(--accent-emerald); padding: 0.4rem 0.8rem; background-color: #0a0a0a;">PUBLISHED LOCAL MOVIES</div>`;
        html += matchingLocal.map(m => {
          const poster = m.poster_path
            ? (m.poster_path.startsWith('http') ? m.poster_path : `https://image.tmdb.org/t/p/w500${m.poster_path}`)
            : 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?q=80&w=1000';
          return `
            <div class="sh-movie-item" data-type="local" data-id="${m.id}" style="display: flex; gap: 0.75rem; padding: 0.6rem 0.8rem; border-bottom: 1px solid var(--border); cursor: pointer; align-items: center;">
              <img src="${poster}" style="width: 32px; aspect-ratio: 2/3; object-fit: cover; border-radius: 4px;" />
              <div>
                <div style="font-size: 0.84rem; font-weight: 700; color: #ffffff;">${m.title}</div>
                <div style="font-size: 0.72rem; color: var(--text-muted);">Local ID: #${m.id} • Published</div>
              </div>
            </div>
          `;
        }).join('');
      }

      // tmdb api catalog results section
      if (tmdbResults.length > 0) {
        html += `<div style="font-size: 0.72rem; font-weight: 700; color: var(--accent); padding: 0.4rem 0.8rem; background-color: #0a0a0a;">TMDB API CATALOG (CLICK TO IMPORT)</div>`;
        html += tmdbResults.slice(0, 8).map(m => {
          const poster = m.poster_path
            ? (m.poster_path.startsWith('http') ? m.poster_path : `https://image.tmdb.org/t/p/w500${m.poster_path}`)
            : 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?q=80&w=1000';
          const title = m.title || m.name || 'Untitled';
          const year = m.release_date ? m.release_date.split('-')[0] : '';
          return `
            <div class="sh-movie-item" data-type="tmdb" data-tmdb-id="${m.id}" style="display: flex; gap: 0.75rem; padding: 0.6rem 0.8rem; border-bottom: 1px solid var(--border); cursor: pointer; align-items: center;">
              <img src="${poster}" style="width: 32px; aspect-ratio: 2/3; object-fit: cover; border-radius: 4px;" />
              <div>
                <div style="font-size: 0.84rem; font-weight: 700; color: #ffffff;">${title}</div>
                <div style="font-size: 0.72rem; color: var(--text-muted);">${year ? year + ' • ' : ''}TMDB ID: #${m.id} (Auto-imports on click)</div>
              </div>
            </div>
          `;
        }).join('');
      }

      if (!html) {
        html = '<div style="padding: 0.8rem; font-size: 0.8125rem; color: var(--text-muted);">No matching movies found in TMDB or local database.</div>';
      }

      resultsBox.innerHTML = html;
      resultsBox.style.display = 'block';

    } catch (err) {
      resultsBox.innerHTML = `<div style="padding: 0.8rem; font-size: 0.8125rem; color: var(--accent);">${err.message || 'Search failed.'}</div>`;
    }
  }

  if (searchInput) {
    searchInput.addEventListener('input', (e) => {
      clearTimeout(debounceTimer);
      const q = e.target.value.trim();
      debounceTimer = setTimeout(() => performTMDBSearch(q), 300);
    });
  }

  // handle clicking a movie from search dropdown
  if (resultsBox) {
    resultsBox.addEventListener('click', async (e) => {
      const item = e.target.closest('.sh-movie-item');
      if (!item) return;

      const type = item.dataset.type;

      if (type === 'local') {
        const mId = Number(item.dataset.id);
        const localMovies = (await API.get('/movies')) || [];
        selectedMovie = localMovies.find(m => m.id === mId);
      } else if (type === 'tmdb') {
        const tmdbId = Number(item.dataset.tmdbId);
        resultsBox.style.display = 'none';
        importStatus.style.display = 'block';

        try {
          // 1. check if already published
          let published = null;
          try {
            published = await API.get(`/movies/tmdb/${tmdbId}`);
          } catch (e) {}

          // 2. publish from tmdb if not already in local database
          if (!published || !published.id) {
            published = await API.post(`/tmdb/publish/${tmdbId}`);
          }

          selectedMovie = published;
          importStatus.style.display = 'none';
        } catch (err) {
          importStatus.style.display = 'none';
          alert(err.message || 'Failed to import movie from TMDB.');
          return;
        }
      }

      if (selectedMovie && selectedMovie.id) {
        const poster = selectedMovie.poster_path
          ? (selectedMovie.poster_path.startsWith('http') ? selectedMovie.poster_path : `https://image.tmdb.org/t/p/w500${selectedMovie.poster_path}`)
          : 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?q=80&w=1000';

        document.getElementById('sh-selected-poster').src = poster;
        document.getElementById('sh-selected-title').textContent = selectedMovie.title;
        document.getElementById('sh-selected-sub').textContent = `Published Local ID: #${selectedMovie.id} • ${selectedMovie.duration || selectedMovie.runtime || 'N/A'} mins`;

        searchContainer.style.display = 'none';
        resultsBox.style.display = 'none';
        selectedCard.style.display = 'flex';
      }
    });
  }

  // remove selected movie button handler
  if (removeMovieBtn) {
    removeMovieBtn.addEventListener('click', () => {
      selectedMovie = null;
      selectedCard.style.display = 'none';
      searchContainer.style.display = 'block';
      if (searchInput) searchInput.value = '';
    });
  }

  // submit schedule show form
  document.getElementById('admin-form-add-show').addEventListener('submit', async (e) => {
    e.preventDefault();
    if (!selectedMovie || !selectedMovie.id) {
      alert('Please search and select a movie for this show!');
      return;
    }

    const movie_id = selectedMovie.id;
    const screen_id = Number(document.getElementById('sh-s-id').value);
    const start_time = new Date(document.getElementById('sh-time').value).toISOString();
    const price = parseFloat(document.getElementById('sh-price').value);

    try {
      await API.post('/shows', { movie_id, screen_id, start_time, price });
      alert('Show scheduled successfully!');
      renderAdminShowsView(viewEl);
    } catch (err) { alert(err.message || 'Failed to schedule show'); }
  });

  // render existing scheduled shows table
  const showsListEl = document.getElementById('admin-shows-list');
  try {
    const shows = await API.get('/shows');
    const localMovies = (await API.get('/movies')) || [];
    const screens = (await API.get('/screens')) || [];
    const theatres = (await API.get('/theatres')) || [];

    if (shows && shows.length > 0) {
      showsListEl.innerHTML = `
        <div class="table-container card" style="padding: 0;">
          <table>
            <thead>
              <tr>
                <th>Show ID</th>
                <th>Movie Title</th>
                <th>Theatre & Screen</th>
                <th>Start Time</th>
                <th>Ticket Price</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              ${shows.map(s => {
                const startTimeStr = s.start_time ? new Date(s.start_time).toLocaleString([], { dateStyle: 'medium', timeStyle: 'short' }) : 'N/A';
                const m = localMovies.find(pm => pm.id === s.movie_id);
                const title = m ? m.title : (s.movie ? s.movie.title : `Movie #${s.movie_id}`);

                const scr = screens.find(sc => sc.id === s.screen_id);
                const th = scr ? theatres.find(t => t.id === scr.theatre_id) : null;
                const theatreName = th ? th.name : 'Theatre';
                const screenNumStr = scr ? (scr.screen_number ? `Screen ${scr.screen_number}` : `Screen #${scr.id}`) : `Screen #${s.screen_id}`;
                const theatreScreenDisplay = `${theatreName} (${screenNumStr})`;

                return `
                  <tr>
                    <td style="font-weight: 700;">#${s.id}</td>
                    <td>${title}</td>
                    <td>${theatreScreenDisplay}</td>
                    <td>${startTimeStr}</td>
                    <td style="color: var(--accent-emerald); font-weight: 700;">Rs. ${s.price}</td>
                    <td>
                      <button class="btn btn-danger btn-sm btn-delete-show" data-id="${s.id}">
                        Cancel Show
                      </button>
                    </td>
                  </tr>
                `;
              }).join('')}
            </tbody>
          </table>
        </div>
      `;

      showsListEl.addEventListener('click', async (e) => {
        const delBtn = e.target.closest('.btn-delete-show');
        if (delBtn) {
          const id = delBtn.dataset.id;
          if (confirm(`Cancel Show #${id}?`)) {
            try {
              await API.delete(`/shows/${id}`);
              renderAdminShowsView(viewEl);
            } catch (err) { alert(err.message || 'Failed to cancel show'); }
          }
        }
      });
    } else {
      showsListEl.innerHTML = '<div class="card" style="text-align: center; color: var(--text-muted); padding: 2rem;">No scheduled shows found.</div>';
    }
  } catch (err) {
    showsListEl.innerHTML = `<div class="card" style="color: var(--accent);">${err.message || 'Failed to load shows.'}</div>`;
  }
}

/* 3. bookings tab view */
async function renderAdminBookingsView(viewEl) {
  viewEl.innerHTML = `
    <div>
      <h1 class="page-title" style="font-size: 2rem; margin-bottom: 0.25rem; color: #ffffff;">Booking Transactions</h1>
      <p style="color: var(--text-muted); font-size: 0.875rem; margin-bottom: 1.5rem;">
        Review customer ticket transactions and manually confirm or cancel bookings.
      </p>

      <div class="card" style="padding: 1rem 1.25rem; margin-bottom: 1.5rem; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 1rem;">
        <div>
          <div style="font-size: 0.875rem; font-weight: 700; color: #ffffff;">Filter by Date</div>
          <div id="tx-count-label" style="font-size: 0.75rem; color: var(--text-muted);">Showing all bookings</div>
        </div>

        <div style="display: flex; gap: 0.5rem; align-items: center;">
          <button id="btn-filter-today" class="btn btn-secondary btn-sm">Today</button>
          <button id="btn-filter-all" class="btn btn-primary btn-sm">All Time</button>
        </div>
      </div>

      <div id="admin-tx-list">
        <div style="font-size: 0.8125rem; color: var(--text-muted);">Loading bookings...</div>
      </div>
    </div>
  `;

  try {
    const allBookings = (await API.get('/admin/bookings/recent')) || [];
    const txList = document.getElementById('admin-tx-list');
    const countLabel = document.getElementById('tx-count-label');
    const btnToday = document.getElementById('btn-filter-today');
    const btnAll = document.getElementById('btn-filter-all');

    let activeFilter = 'all'; // 'all' or 'today'

    function renderFilteredBookings() {
      const todayStr = new Date().toISOString().split('T')[0];
      const filtered = activeFilter === 'today'
        ? allBookings.filter(b => {
            if (!b.created_at) return false;
            try {
              const dStr = new Date(b.created_at).toISOString().split('T')[0];
              return dStr === todayStr;
            } catch (e) { return false; }
          })
        : allBookings;

      if (activeFilter === 'today') {
        btnToday.className = 'btn btn-primary btn-sm';
        btnAll.className = 'btn btn-secondary btn-sm';
        countLabel.textContent = `Showing today's transactions • ${filtered.length} records`;
      } else {
        btnAll.className = 'btn btn-primary btn-sm';
        btnToday.className = 'btn btn-secondary btn-sm';
        countLabel.textContent = `Showing all bookings • ${filtered.length} records`;
      }

      if (filtered && filtered.length > 0) {
        const now = new Date();
        const todayStr = now.toISOString().split('T')[0];

        txList.innerHTML = filtered.map(b => {
          const isConfirmed = b.status === 'CONFIRMED' || b.status === 'SUCCESS';
          const isPending = b.status === 'PENDING';
          const createdDate = b.created_at ? new Date(b.created_at).toLocaleString([], { dateStyle: 'medium', timeStyle: 'short' }) : 'N/A';

          const createdDateStr = b.created_at ? new Date(b.created_at).toISOString().split('T')[0] : '';
          const isToday = (createdDateStr === todayStr);

          const showEndTime = b.show_end_time ? new Date(b.show_end_time) : (b.show_start_time ? new Date(new Date(b.show_start_time).getTime() + 3 * 3600 * 1000) : null);
          const isShowEnded = showEndTime ? (now > showEndTime) : false;

          const canModify = isToday && !isShowEnded;

          return `
            <div class="booking-tx-card">
              <div class="tx-info">
                <div>
                  <div class="tx-id">Booking #${b.booking_id}</div>
                  <div class="tx-sub">${createdDate}</div>
                </div>
                <div>
                  <div class="tx-meta-label">Movie Title</div>
                  <div class="tx-meta-val">${b.movie_title || 'N/A'}</div>
                </div>
                <div>
                  <div class="tx-meta-label">Customer</div>
                  <div class="tx-meta-val">${b.customer_name || 'Customer'}</div>
                </div>
                <div>
                  <div class="tx-meta-label">Amount</div>
                  <div class="tx-meta-val" style="color: var(--accent-emerald);">Rs. ${b.amount ? b.amount.toFixed(2) : '0.00'}</div>
                </div>
              </div>

              <div style="display: flex; gap: 0.75rem; align-items: center;">
                <span class="badge ${isConfirmed ? 'badge-confirmed' : (isPending ? 'badge-pending' : 'badge-cancelled')}">${b.status}</span>
                
                ${canModify && isPending ? `
                  <button class="btn btn-primary btn-sm btn-admin-confirm" data-id="${b.booking_id}" title="Confirm Booking">✓</button>
                ` : ''}
                
                ${canModify && b.status !== 'CANCELLED' ? `
                  <button class="btn btn-danger btn-sm btn-admin-cancel" data-id="${b.booking_id}" title="Cancel Booking">✕</button>
                ` : ''}
              </div>
            </div>
          `;
        }).join('');
      } else {
        txList.innerHTML = `<div class="card" style="text-align: center; color: var(--text-muted); padding: 2rem;">No ${activeFilter === 'today' ? "today's" : ''} booking records found.</div>`;
      }
    }

    btnToday.addEventListener('click', () => {
      activeFilter = 'today';
      renderFilteredBookings();
    });

    btnAll.addEventListener('click', () => {
      activeFilter = 'all';
      renderFilteredBookings();
    });

    renderFilteredBookings();

    txList.addEventListener('click', async (e) => {
      const confirmBtn = e.target.closest('.btn-admin-confirm');
      const cancelBtn = e.target.closest('.btn-admin-cancel');

      if (confirmBtn) {
        const id = confirmBtn.dataset.id;
        try {
          await API.put(`/bookings/${id}/confirm`);
          renderAdminBookingsView(viewEl);
        } catch (err) { alert(err.message || 'Failed to confirm booking'); }
      }

      if (cancelBtn) {
        const id = cancelBtn.dataset.id;
        if (confirm('Are you sure you want to cancel booking?')) {
          try {
            await API.delete(`/bookings/${id}`);
            renderAdminBookingsView(viewEl);
          } catch (err) { alert(err.message || 'Failed to cancel booking'); }
        }
      }
    });

  } catch (err) {
    document.getElementById('admin-tx-list').innerHTML = `<div class="card" style="color: var(--accent);">${err.message || 'Failed to load bookings.'}</div>`;
  }
}

/* 4. dashboard tab view */
async function renderAdminDashboardView(viewEl) {
  viewEl.innerHTML = '<div style="font-size: 0.8125rem; color: var(--text-muted);">Loading stats...</div>';
  try {
    const stats = await API.get('/admin/dashboard');
    viewEl.innerHTML = `
      <h1 class="page-title">Admin Dashboard</h1>
      <div class="grid" style="grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));">
        <div class="card" style="text-align: center;">
          <div style="font-size: 0.75rem; color: var(--text-muted);">Total Theatres</div>
          <div style="font-size: 2rem; font-weight: 900; color: #ffffff;">${stats.total_theatres || 0}</div>
        </div>
        <div class="card" style="text-align: center;">
          <div style="font-size: 0.75rem; color: var(--text-muted);">Total Screens</div>
          <div style="font-size: 2rem; font-weight: 900; color: #ffffff;">${stats.total_screens || 0}</div>
        </div>
        <div class="card" style="text-align: center;">
          <div style="font-size: 0.75rem; color: var(--text-muted);">Active Shows</div>
          <div style="font-size: 2rem; font-weight: 900; color: #ffffff;">${stats.total_shows || 0}</div>
        </div>
        <div class="card" style="text-align: center;">
          <div style="font-size: 0.75rem; color: var(--text-muted);">Total Bookings</div>
          <div style="font-size: 2rem; font-weight: 900; color: #ffffff;">${stats.total_bookings || 0}</div>
        </div>
        <div class="card" style="text-align: center;">
          <div style="font-size: 0.75rem; color: var(--text-muted);">Total Revenue</div>
          <div style="font-size: 1.75rem; font-weight: 900; color: var(--accent-emerald);">Rs. ${stats.total_revenue || 0}</div>
        </div>
      </div>
    `;
  } catch (err) {
    viewEl.innerHTML = `<div class="card" style="color: var(--accent);">${err.message || 'Failed to load stats.'}</div>`;
  }
}

Router.addRoute('#/admin', renderAdmin);
