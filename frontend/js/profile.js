/* user profile & favorite genres module using exact tmdb genre ids */

async function renderProfile(container) {
  if (!API.getToken()) {
    window.location.hash = '#/login';
    return;
  }

  container.innerHTML = '<div style="text-align: center; padding: 3rem; color: var(--text-muted);">Loading profile...</div>';

  try {
    const profile = await API.get('/users/me');
    const allGenres = await API.get('/genres');
    let userFavGenres = [];

    try {
      userFavGenres = await API.get('/users/me/genres');
    } catch (e) { console.error('Failed to load favorite genres:', e); }

    // backend add/remove favorite genres requires tmdb genre id (tmdb_id)
    const favTmdbIds = (userFavGenres || []).map(g => g.tmdb_id || g.TMDBID || g.id);

    container.innerHTML = `
      <h1 class="page-title">User Profile</h1>

      <div class="card" style="max-width: 600px;">
        <h2 class="section-title">Account Details</h2>
        <div style="display: flex; flex-direction: column; gap: 0.5rem; font-size: 0.875rem; margin-bottom: 1.5rem;">
          <div><strong>Name:</strong> ${profile.name || 'N/A'}</div>
          <div><strong>Email:</strong> ${profile.email || 'N/A'}</div>
          <div><strong>Role:</strong> <span class="badge badge-confirmed">${profile.role || 'Customer'}</span></div>
        </div>

        <hr style="border-color: var(--border); margin: 1.5rem 0;" />

        <h2 class="section-title">Favorite Genres</h2>
        <p style="font-size: 0.8125rem; color: var(--text-muted); margin-bottom: 1rem;">
          Select your favorite genres to personalize your recommendations.
        </p>

        <div id="profile-genre-chips" class="genre-chips" style="margin-bottom: 1.5rem;">
          ${allGenres.map(g => {
            const tmdbId = g.tmdb_id || g.id;
            const isFav = favTmdbIds.includes(tmdbId);
            return `
              <div class="genre-chip ${isFav ? 'active' : ''}" data-tmdb-id="${tmdbId}">
                ${g.name}
              </div>
            `;
          }).join('')}
        </div>
      </div>
    `;

    // chip click listener to toggle favorites
    const chipsContainer = document.getElementById('profile-genre-chips');
    chipsContainer.addEventListener('click', async (e) => {
      const chip = e.target.closest('.genre-chip');
      if (!chip) return;
      const tmdbId = Number(chip.dataset.tmdbId);
      const isFav = chip.classList.contains('active');

      try {
        if (isFav) {
          await API.delete(`/users/me/genres/${tmdbId}`);
          chip.classList.remove('active');
        } else {
          await API.post(`/users/me/genres/${tmdbId}`);
          chip.classList.add('active');
        }
      } catch (err) {
        alert(err.message || 'Failed to update favorite genre');
      }
    });

  } catch (err) {
    container.innerHTML = `<div class="card" style="color: var(--accent);">${err.message || 'Failed to load profile.'}</div>`;
  }
}

Router.addRoute('#/profile', renderProfile);
