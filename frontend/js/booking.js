/* seat grid & booking summary module */

let selectedSeatsList = []; // objects: { id, row, number, seatNumber }

async function renderSeatLayout(container, params) {
  const showId = params.showId;
  if (!API.getToken()) {
    window.location.hash = '#/login';
    return;
  }

  container.innerHTML = '<div style="text-align: center; padding: 3rem; color: var(--text-muted);">Loading seat map...</div>';

  try {
    const show = await API.get(`/shows/${showId}`);

    // resolve movie details: check preloaded show.movie first, then fall back to local/tmdb lookup
    let movie = (show && show.movie && show.movie.title) ? show.movie : null;
    if (!movie) {
      const movies = (await API.get('/movies')) || [];
      movie = movies.find(m => m.id === show.movie_id || m.tmdb_id === show.movie_id);
    }
    if (!movie || !movie.title) {
      try { movie = await API.get(`/movies/${show.movie_id}`); } catch (e) {}
    }
    if (!movie || !movie.title) {
      try { movie = await API.get(`/movies/tmdb/${show.movie_id}`); } catch (e) {}
    }
    if (!movie || !movie.title) {
      try { movie = await API.get(`/tmdb/movie/${show.movie_id}`); } catch (e) {}
    }
    if (!movie || !movie.title) {
      movie = { title: `Movie #${show.movie_id}`, poster_path: '' };
    }

    const screen = await API.get(`/screens/${show.screen_id}`) || { total_seats: 60, screen_number: 1 };
    const theatre = await API.get(`/theatres/${screen.theatre_id}`) || { name: 'Theatre' };

    let seats = [];
    try {
      const seatRes = await API.get(`/shows/${showId}/seats`);
      if (seatRes && seatRes.rows) {
        // flat array of seats from show seat layout response dto
        seatRes.rows.forEach(r => {
          (r.seats || []).forEach(s => {
            seats.push({
              id: s.seat_id,
              row: r.row,
              number: parseInt(s.seat_number.replace(/^[A-Z]+/i, '')) || 1,
              seatNumber: s.seat_number,
              is_booked: s.status === 'booked'
            });
          });
        });
      }
    } catch (e) {}

    if (seats.length === 0) {
      container.innerHTML = `
        <div style="max-width: 600px; margin: 3rem auto; text-align: center; background-color: #121212; border: 1px solid var(--border); border-radius: 12px; padding: 3rem 2rem;">
          <div style="font-size: 2.5rem; margin-bottom: 1rem;">⚠️</div>
          <h2 style="color: #ffffff; font-size: 1.3rem; font-weight: 800; margin-bottom: 0.5rem;">No Seats Configured</h2>
          <p style="color: var(--text-muted); font-size: 0.875rem; margin-bottom: 1.5rem;">
            There are currently no seats configured for this show screen in the database. Please contact the theatre administrator or select another showtime.
          </p>
          <a href="#/shows/${theatre.id || 1}" class="btn btn-secondary" style="display: inline-block;">&larr; Back to Showtimes</a>
        </div>
      `;
      return;
    }

    selectedSeatsList = [];

    // group seats by row
    const rowMap = {};
    seats.forEach(s => {
      if (!rowMap[s.row]) rowMap[s.row] = [];
      rowMap[s.row].push(s);
    });

    const screenName = screen.screen_number ? `Screen ${screen.screen_number}` : (screen.name || `Screen #${screen.id}`);
    const showTimeStr = show.start_time ? new Date(show.start_time).toLocaleString([], { dateStyle: 'medium', timeStyle: 'short' }) : 'Showtime';

    container.innerHTML = `
      <div style="max-width: 1200px; margin: 0 auto; padding-top: 1rem;">
        
        <a href="#/shows/${theatre.id || 1}" style="font-size: 0.8125rem; color: var(--text-muted); text-decoration: none; display: inline-block; margin-bottom: 1.5rem;">
          &larr; Back to Showtimes
        </a>

        <!-- main layout: seat map on left (70%), booking summary on right (30%) -->
        <div style="display: flex; gap: 2rem; flex-wrap: wrap; align-items: flex-start;">
          
          <!-- left side: seat map area -->
          <div style="flex: 2; min-width: 320px; background-color: var(--surface); border: 1px solid var(--border); border-radius: 6px; padding: 1.5rem; text-align: center;">
            
            <div style="font-size: 0.7rem; font-weight: 800; color: var(--accent); letter-spacing: 1px; text-transform: uppercase;">SCREEN SELECTION</div>
            <h2 style="font-size: 1.25rem; font-weight: 800; color: #ffffff; margin-bottom: 1.25rem;">${screenName}</h2>

            <!-- screen bar -->
            <div style="margin: 0 auto 1.5rem auto; max-width: 400px; padding: 0.35rem; background-color: var(--surface-light); border: 1px solid var(--border); border-radius: 4px; font-size: 0.7rem; font-weight: 700; color: var(--muted); letter-spacing: 2px; text-transform: uppercase;">
              SCREEN THIS WAY
            </div>

            <!-- legend -->
            <div style="display: flex; justify-content: center; gap: 1.5rem; margin-bottom: 1.5rem; font-size: 0.75rem; color: var(--muted);">
              <div style="display: flex; align-items: center; gap: 0.4rem;">
                <div class="seat-grid-cell available" style="width: 16px; height: 16px;"></div> Available
              </div>
              <div style="display: flex; align-items: center; gap: 0.4rem;">
                <div class="seat-grid-cell selected" style="width: 16px; height: 16px;"></div> Selected
              </div>
              <div style="display: flex; align-items: center; gap: 0.4rem;">
                <div class="seat-grid-cell booked" style="width: 16px; height: 16px;"></div> Booked
              </div>
            </div>

            <!-- interactive seat grid container -->
            <div id="seat-layout-grid" style="display: flex; flex-direction: column; gap: 0.5rem; align-items: center; overflow-x: auto; padding-bottom: 1rem;">
              ${Object.keys(rowMap).sort().map(row => `
                <div style="display: flex; gap: 0.4rem; align-items: center;">
                  <span style="font-size: 0.75rem; font-weight: 800; width: 22px; text-align: right; color: #ef4444; margin-right: 0.4rem;">${row}</span>
                  ${rowMap[row].sort((a,b) => a.number - b.number).map(seat => `
                    <div class="seat-grid-cell ${seat.is_booked ? 'booked' : 'available'}" 
                         data-id="${seat.id}" 
                         data-row="${seat.row}" 
                         data-num="${seat.number}"
                         data-code="${seat.seatNumber}">
                      ${seat.number}
                    </div>
                  `).join('')}
                </div>
              `).join('')}
            </div>

          </div>

          <!-- right sidebar: booking summary -->
          <div style="flex: 1; min-width: 280px; background-color: var(--surface); border: 1px solid var(--border); border-radius: 6px; padding: 1.25rem;">
            
            <div style="margin-bottom: 1rem;">
              <h2 style="font-size: 1rem; font-weight: 800; color: #ffffff; margin: 0;">BOOKING SUMMARY</h2>
            </div>

            <div style="display: flex; flex-direction: column; gap: 0.875rem; font-size: 0.8125rem; margin-bottom: 1.25rem;">
              
              <div>
                <div style="font-size: 0.7rem; color: var(--muted); font-weight: 700; text-transform: uppercase; margin-bottom: 0.1rem;">MOVIE</div>
                <div style="font-size: 0.95rem; font-weight: 800; color: #ffffff;">${movie.title}</div>
              </div>

              <div>
                <div style="font-size: 0.7rem; color: var(--text-muted); font-weight: 700; text-transform: uppercase; margin-bottom: 0.2rem;">THEATRE & SCREEN</div>
                <div style="font-size: 0.85rem; font-weight: 700; color: #e5e5e5;">${theatre.name} (${screenName})</div>
              </div>

              <div>
                <div style="font-size: 0.7rem; color: var(--text-muted); font-weight: 700; text-transform: uppercase; margin-bottom: 0.2rem;">SHOW TIME</div>
                <div style="font-size: 0.85rem; font-weight: 600; color: #e5e5e5;">${showTimeStr}</div>
              </div>

              <div>
                <div style="font-size: 0.7rem; color: var(--text-muted); font-weight: 700; text-transform: uppercase; margin-bottom: 0.4rem;">SELECTED SEATS</div>
                <div id="summary-seats-pills" style="display: flex; gap: 0.4rem; flex-wrap: wrap;">
                  <span style="color: var(--text-muted); font-style: italic;">No seats selected</span>
                </div>
              </div>

            </div>

            <hr style="border-color: var(--border); margin-bottom: 1.25rem;" />

            <!-- price breakdown -->
            <div style="display: flex; flex-direction: column; gap: 0.5rem; font-size: 0.8125rem; margin-bottom: 1.5rem;">
              <div style="display: flex; justify-content: space-between; color: var(--text-muted);">
                <span>Ticket Price</span>
                <span id="summary-ticket-price">Rs. ${show.price} x 0</span>
              </div>
              <div style="display: flex; justify-content: space-between; color: var(--text-muted);">
                <span>Convenience Fee</span>
                <span>Rs. 0.00</span>
              </div>
              <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 0.5rem;">
                <span style="font-size: 0.9rem; font-weight: 800; color: #ffffff;">Total Payable</span>
                <span id="summary-total-price" style="font-size: 1.4rem; font-weight: 900; color: #ef4444;">Rs. 0.00</span>
              </div>
            </div>

            <!-- submit booking button -->
            <button id="btn-create-booking" class="btn btn-primary" style="width: 100%; padding: 0.85rem; font-size: 0.95rem; font-weight: 800;" disabled>
              Book & Pay Rs. 0.00
            </button>

          </div>

        </div>

      </div>
    `;

    // interactive seat click event handler
    const grid = document.getElementById('seat-layout-grid');
    const seatsPillsEl = document.getElementById('summary-seats-pills');
    const ticketPriceEl = document.getElementById('summary-ticket-price');
    const totalPriceEl = document.getElementById('summary-total-price');
    const bookBtn = document.getElementById('btn-create-booking');

    grid.addEventListener('click', (e) => {
      const cell = e.target.closest('.seat-grid-cell.available');
      if (!cell) return;

      const id = Number(cell.dataset.id);
      const code = cell.dataset.code;

      const idx = selectedSeatsList.findIndex(s => s.id === id);
      if (idx > -1) {
        selectedSeatsList.splice(idx, 1);
        cell.style.backgroundColor = '#262626';
        cell.style.borderColor = '#3d3d3d';
      } else {
        selectedSeatsList.push({ id, code });
        cell.style.backgroundColor = '#ef4444';
        cell.style.borderColor = '#ef4444';
      }

      // update right sidebar summary
      const count = selectedSeatsList.length;
      const total = count * show.price;

      if (count > 0) {
        seatsPillsEl.innerHTML = selectedSeatsList.map(s => `
          <span style="background-color: rgba(239, 68, 68, 0.2); color: #ef4444; border: 1px solid #ef4444; padding: 0.2rem 0.5rem; border-radius: 4px; font-weight: 800; font-size: 0.75rem;">
            ${s.code}
          </span>
        `).join('');
        ticketPriceEl.textContent = `Rs. ${show.price} x ${count}`;
        totalPriceEl.textContent = `Rs. ${total.toFixed(2)}`;
        bookBtn.textContent = `Book & Pay Rs. ${total.toFixed(2)}`;
        bookBtn.disabled = false;
      } else {
        seatsPillsEl.innerHTML = '<span style="color: var(--text-muted); font-style: italic;">No seats selected</span>';
        ticketPriceEl.textContent = `Rs. ${show.price} x 0`;
        totalPriceEl.textContent = `Rs. 0.00`;
        bookBtn.textContent = `Book & Pay Rs. 0.00`;
        bookBtn.disabled = true;
      }
    });

    // handle create booking & payment
    bookBtn.addEventListener('click', async () => {
      const seatIds = selectedSeatsList.map(s => s.id);
      try {
        const booking = await API.post('/bookings', {
          show_id: Number(showId),
          seat_ids: seatIds,
          payment_method: 'RAZORPAY'
        });

        if (booking.razorpay_order_id && booking.razorpay_key_id) {
          const options = {
            key: booking.razorpay_key_id,
            amount: booking.total_amount * 100,
            currency: 'INR',
            name: 'Movie Booking',
            description: `Payment for Booking #${booking.id}`,
            order_id: booking.razorpay_order_id,
            handler: async function (rzpResponse) {
              try {
                await API.post('/payments/verify', {
                  razorpay_order_id: rzpResponse.razorpay_order_id,
                  razorpay_payment_id: rzpResponse.razorpay_payment_id,
                  razorpay_signature: rzpResponse.razorpay_signature,
                });
                alert('Payment verified successfully!');
                window.location.hash = '#/bookings';
              } catch (e) {
                alert('Payment verification error: ' + e.message);
              }
            },
            theme: { color: '#ef4444' }
          };

          const rzp = new window.Razorpay(options);
          rzp.open();
        } else {
          window.location.hash = `#/bookings/confirm/${booking.id}`;
        }
      } catch (err) {
        alert(err.message || 'Failed to create booking.');
      }
    });

  } catch (err) {
    container.innerHTML = `<div class="card" style="color: var(--accent);">${err.message || 'Failed to load seat layout.'}</div>`;
  }
}

async function renderTicketConfirm(container, params) {
  const bookingId = params.id;
  container.innerHTML = '<div style="text-align: center; padding: 3rem; color: var(--text-muted);">Loading booking details...</div>';

  try {
    const booking = await API.get(`/bookings/${bookingId}`);

    container.innerHTML = `
      <h1 class="page-title">Booking Confirmation</h1>
      <div class="card" style="max-width: 450px; margin: 0 auto;">
        <h2 style="font-size: 1.125rem; font-weight: 800; margin-bottom: 1rem;">Booking #${booking.id}</h2>
        <div style="display: flex; flex-direction: column; gap: 0.5rem; font-size: 0.875rem; margin-bottom: 1.5rem;">
          <div><strong>Status:</strong> <span class="badge ${booking.status === 'CONFIRMED' ? 'badge-confirmed' : 'badge-pending'}">${booking.status}</span></div>
          <div><strong>Total Amount:</strong> <span style="color: var(--accent-emerald); font-weight: 700;">Rs. ${booking.total_amount}</span></div>
        </div>

        <a href="#/bookings" class="btn btn-secondary" style="width: 100%; text-align: center;">View All My Bookings</a>
      </div>
    `;
  } catch (err) {
    container.innerHTML = `<div class="card" style="color: var(--accent);">${err.message || 'Failed to load booking.'}</div>`;
  }
}

async function renderMyBookings(container) {
  if (!API.getToken()) {
    window.location.hash = '#/login';
    return;
  }

  container.innerHTML = '<div style="text-align: center; padding: 3rem; color: var(--text-muted);">Loading bookings...</div>';

  try {
    const user = API.getUser();
    if (!user || !user.user_id) {
      throw new Error('Session expired. Please login again.');
    }
    const bookings = await API.get(`/bookings/user/${user.user_id}`);

    container.innerHTML = `
      <h1 class="page-title">My Bookings</h1>
      ${bookings && bookings.length > 0 ? `
        <div class="table-container card" style="padding: 0;">
          <table>
            <thead>
              <tr>
                <th>Booking ID</th>
                <th>Show ID</th>
                <th>Total Amount</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              ${bookings.map(b => `
                <tr>
                  <td style="font-weight: 700;">#${b.id}</td>
                  <td>Show #${b.show_id}</td>
                  <td style="color: var(--accent-emerald); font-weight: 700;">Rs. ${b.total_amount}</td>
                  <td><span class="badge ${b.status === 'CONFIRMED' ? 'badge-confirmed' : 'badge-pending'}">${b.status}</span></td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      ` : `
        <div class="card" style="text-align: center; padding: 3rem; color: var(--text-muted);">
          You have no active bookings yet.
        </div>
      `}
    `;
  } catch (err) {
    container.innerHTML = `<div class="card" style="color: var(--accent);">${err.message || 'Failed to load bookings.'}</div>`;
  }
}

Router.addRoute('#/book/:showId', renderSeatLayout);
Router.addRoute('#/bookings/confirm/:id', renderTicketConfirm);
Router.addRoute('#/bookings', renderMyBookings);
