/* Authentication Controller with 5-Minute OTP Countdown Timer */

let otpTimerInterval = null;

function renderLogin(container) {
  if (otpTimerInterval) {
    clearInterval(otpTimerInterval);
    otpTimerInterval = null;
  }

  container.innerHTML = `
    <div class="card" style="max-width: 400px; margin: 2rem auto;">
      <h1 class="page-title" style="text-align: center; margin-bottom: 0.5rem;">Login</h1>
      <p style="text-align: center; color: var(--text-muted); font-size: 0.875rem; margin-bottom: 1.5rem;">
        Access your account to book tickets
      </p>

      <form id="form-login">
        <div class="form-group">
          <label class="form-label" for="login-email">Email Address</label>
          <input type="email" id="login-email" class="form-input" required placeholder="user@example.com" />
        </div>

        <div class="form-group">
          <label class="form-label" for="login-password">Password</label>
          <input type="password" id="login-password" class="form-input" required placeholder="••••••••" />
        </div>

        <div id="auth-error" style="color: var(--brand-primary); font-size: 0.8125rem; margin-bottom: 1rem;" class="hidden"></div>

        <button type="submit" class="btn btn-primary" style="width: 100%;">Login</button>
      </form>

      <div style="margin-top: 1.25rem; text-align: center; font-size: 0.8125rem; color: var(--text-muted); display: flex; flex-direction: column; gap: 0.5rem;">
        <div>Don't have an account? <a href="#/signup" style="color: var(--brand-primary); font-weight: 600;">Sign up</a></div>
        <div><a href="#/forgot-password" style="color: var(--text-muted);">Forgot Password?</a></div>
      </div>
    </div>
  `;

  document.getElementById('form-login').addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;
    const errorDiv = document.getElementById('auth-error');
    errorDiv.classList.add('hidden');

    try {
      const res = await API.post('/auth/login', { email, password });
      API.setToken(res.token);
      window.location.hash = '#/';
    } catch (err) {
      const errMsg = err.message || 'Invalid email or password.';
      if (errMsg.toLowerCase().includes('verify your email')) {
        errorDiv.innerHTML = `
          <div>${errMsg}</div>
          <a href="#/verify-otp?email=${encodeURIComponent(email)}" class="btn btn-secondary btn-sm" style="margin-top: 0.5rem; display: inline-block;">Verify Email Now</a>
        `;
      } else {
        errorDiv.textContent = errMsg;
      }
      errorDiv.classList.remove('hidden');
    }
  });
}

function renderSignup(container) {
  if (otpTimerInterval) {
    clearInterval(otpTimerInterval);
    otpTimerInterval = null;
  }

  container.innerHTML = `
    <div class="card" style="max-width: 400px; margin: 2rem auto;">
      <h1 class="page-title" style="text-align: center; margin-bottom: 0.5rem;">Sign Up</h1>
      <p style="text-align: center; color: var(--text-muted); font-size: 0.875rem; margin-bottom: 1.5rem;">
        Create a new account
      </p>

      <form id="form-signup">
        <div class="form-group">
          <label class="form-label" for="signup-name">Full Name</label>
          <input type="text" id="signup-name" class="form-input" required placeholder="John Doe" />
        </div>

        <div class="form-group">
          <label class="form-label" for="signup-email">Email Address</label>
          <input type="email" id="signup-email" class="form-input" required placeholder="user@example.com" />
        </div>

        <div class="form-group">
          <label class="form-label" for="signup-password">Password</label>
          <input type="password" id="signup-password" class="form-input" required minlength="6" placeholder="At least 6 characters" />
        </div>

        <div id="auth-error" style="color: var(--brand-primary); font-size: 0.8125rem; margin-bottom: 1rem;" class="hidden"></div>

        <button type="submit" class="btn btn-primary" style="width: 100%;">Create Account</button>
      </form>

      <div style="margin-top: 1.25rem; text-align: center; font-size: 0.8125rem; color: var(--text-muted);">
        Already registered? <a href="#/login" style="color: var(--brand-primary); font-weight: 600;">Login</a>
      </div>
    </div>
  `;

  document.getElementById('form-signup').addEventListener('submit', async (e) => {
    e.preventDefault();
    const name = document.getElementById('signup-name').value;
    const email = document.getElementById('signup-email').value;
    const password = document.getElementById('signup-password').value;
    const errorDiv = document.getElementById('auth-error');
    errorDiv.classList.add('hidden');

    try {
      await API.post('/auth/register', { name, email, password });
      window.location.hash = `#/verify-otp?email=${encodeURIComponent(email)}`;
    } catch (err) {
      errorDiv.textContent = err.message || 'Registration failed.';
      errorDiv.classList.remove('hidden');
    }
  });
}

function renderVerifyOTP(container) {
  // Clear any existing countdown timer
  if (otpTimerInterval) {
    clearInterval(otpTimerInterval);
    otpTimerInterval = null;
  }

  const urlParams = new URLSearchParams(window.location.hash.split('?')[1] || '');
  const emailFromUrl = urlParams.get('email') || '';

  container.innerHTML = `
    <div class="card" style="max-width: 400px; margin: 2rem auto;">
      <h1 class="page-title" style="text-align: center; margin-bottom: 0.5rem;">Verify Email</h1>
      <p style="text-align: center; color: var(--text-muted); font-size: 0.875rem; margin-bottom: 1rem;">
        Enter the 6-digit OTP sent to your registered email address
      </p>

      <!-- 10-Minute Timer Display -->
      <div id="otp-timer-wrapper" style="text-align: center; margin-bottom: 1.25rem; background-color: #0a0a0a; padding: 0.6rem; border-radius: 6px; border: 1px solid var(--border-dark);">
        <div style="font-size: 0.75rem; color: var(--text-muted);">OTP Expires In</div>
        <div id="otp-timer-display" style="font-size: 1.5rem; font-weight: 800; color: var(--brand-primary); font-family: monospace; margin-top: 0.1rem;">10:00</div>
      </div>

      <form id="form-otp">
        <div class="form-group">
          <label class="form-label" for="otp-email">Email Address</label>
          <input type="email" id="otp-email" class="form-input" required value="${emailFromUrl}" placeholder="user@example.com" />
        </div>

        <div class="form-group">
          <label class="form-label" for="otp-code">OTP Code</label>
          <input type="text" id="otp-code" class="form-input" required placeholder="123456" maxlength="6" style="text-align: center; font-size: 1.25rem; letter-spacing: 0.25em; font-family: monospace;" />
        </div>

        <div id="auth-error" style="color: var(--brand-primary); font-size: 0.8125rem; margin-bottom: 1rem;" class="hidden"></div>
        <div id="auth-success" style="color: var(--accent-emerald); font-size: 0.8125rem; margin-bottom: 1rem;" class="hidden"></div>

        <button type="submit" class="btn btn-primary" style="width: 100%;">Verify OTP</button>
      </form>

      <div style="margin-top: 1rem; text-align: center; display: flex; flex-direction: column; gap: 0.75rem; align-items: center;">
        <button id="btn-resend-otp" class="btn btn-secondary btn-sm">Resend OTP</button>
        <div style="font-size: 0.8125rem; color: var(--text-muted);">
          Need to change email? <a href="#/signup" style="color: var(--brand-primary); font-weight: 600;">Sign up again</a>
        </div>
      </div>
    </div>
  `;

  // Timer Implementation (600 seconds = 10 minutes)
  function startOTPTimer(seconds) {
    if (otpTimerInterval) clearInterval(otpTimerInterval);

    let remaining = seconds;
    const timerDisplay = document.getElementById('otp-timer-display');
    const resendBtn = document.getElementById('btn-resend-otp');
    const errorDiv = document.getElementById('auth-error');

    if (timerDisplay) timerDisplay.style.color = 'var(--brand-primary)';

    function updateDisplay() {
      if (!timerDisplay) return;
      const mins = Math.floor(remaining / 60);
      const secs = remaining % 60;
      timerDisplay.textContent = `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;

      if (remaining <= 0) {
        clearInterval(otpTimerInterval);
        otpTimerInterval = null;
        timerDisplay.textContent = '00:00';
        timerDisplay.style.color = 'var(--text-muted)';
        if (errorDiv) {
          errorDiv.textContent = 'OTP code expired. Please click "Resend OTP" for a new code.';
          errorDiv.classList.remove('hidden');
        }
      } else {
        remaining--;
      }
    }

    updateDisplay();
    otpTimerInterval = setInterval(updateDisplay, 1000);
  }

  // Start 10-minute timer (600s) on render
  startOTPTimer(600);

  document.getElementById('form-otp').addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('otp-email').value.trim();
    const otp = document.getElementById('otp-code').value.trim();
    const errorDiv = document.getElementById('auth-error');
    const successDiv = document.getElementById('auth-success');
    errorDiv.classList.add('hidden');
    successDiv.classList.add('hidden');

    try {
      await API.post('/auth/verify-otp', { email, otp });
      if (otpTimerInterval) clearInterval(otpTimerInterval);
      successDiv.textContent = 'Email verified successfully! Redirecting to login...';
      successDiv.classList.remove('hidden');
      setTimeout(() => { window.location.hash = '#/login'; }, 1500);
    } catch (err) {
      errorDiv.textContent = err.message || 'OTP verification failed.';
      errorDiv.classList.remove('hidden');
    }
  });

  document.getElementById('btn-resend-otp').addEventListener('click', async () => {
    const email = document.getElementById('otp-email').value.trim();
    const errorDiv = document.getElementById('auth-error');
    const successDiv = document.getElementById('auth-success');
    errorDiv.classList.add('hidden');
    successDiv.classList.add('hidden');

    if (!email) {
      errorDiv.textContent = 'Please enter your email address to resend OTP.';
      errorDiv.classList.remove('hidden');
      return;
    }

    try {
      await API.post('/auth/resend-otp', { email });
      successDiv.textContent = 'A new 6-digit OTP has been sent to your email.';
      successDiv.classList.remove('hidden');
      // Restart 10-minute timer from 10:00 (600s)
      startOTPTimer(600);
    } catch (err) {
      errorDiv.textContent = err.message || 'Failed to resend OTP.';
      errorDiv.classList.remove('hidden');
    }
  });
}

function renderForgotPassword(container) {
  if (otpTimerInterval) {
    clearInterval(otpTimerInterval);
    otpTimerInterval = null;
  }

  container.innerHTML = `
    <div class="card" style="max-width: 400px; margin: 2rem auto;">
      <h1 class="page-title" style="text-align: center; margin-bottom: 0.5rem;">Forgot Password</h1>
      <p style="text-align: center; color: var(--text-muted); font-size: 0.875rem; margin-bottom: 1.5rem;">
        Enter your email to receive a password reset OTP
      </p>

      <form id="form-forgot">
        <div class="form-group">
          <label class="form-label" for="forgot-email">Email Address</label>
          <input type="email" id="forgot-email" class="form-input" required placeholder="user@example.com" />
        </div>

        <div id="auth-error" style="color: var(--brand-primary); font-size: 0.8125rem; margin-bottom: 1rem;" class="hidden"></div>
        <div id="auth-success" style="color: var(--accent-emerald); font-size: 0.8125rem; margin-bottom: 1rem;" class="hidden"></div>

        <button type="submit" class="btn btn-primary" style="width: 100%;">Send Reset Code</button>
      </form>
    </div>
  `;

  document.getElementById('form-forgot').addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('forgot-email').value;
    const errorDiv = document.getElementById('auth-error');
    const successDiv = document.getElementById('auth-success');
    errorDiv.classList.add('hidden');
    successDiv.classList.add('hidden');

    try {
      await API.post('/auth/forgot-password', { email });
      window.location.hash = `#/reset-password?email=${encodeURIComponent(email)}`;
    } catch (err) {
      errorDiv.textContent = err.message || 'Failed to process forgot password request.';
      errorDiv.classList.remove('hidden');
    }
  });
}

function renderResetPassword(container) {
  if (otpTimerInterval) {
    clearInterval(otpTimerInterval);
    otpTimerInterval = null;
  }

  const urlParams = new URLSearchParams(window.location.hash.split('?')[1] || '');
  const emailFromUrl = urlParams.get('email') || '';

  container.innerHTML = `
    <div class="card" style="max-width: 400px; margin: 2rem auto;">
      <h1 class="page-title" style="text-align: center; margin-bottom: 0.5rem;">Reset Password</h1>
      <p style="text-align: center; color: var(--text-muted); font-size: 0.875rem; margin-bottom: 1rem;">
        Enter the 6-digit OTP sent to your registered email address
      </p>

      <!-- 10-Minute Timer Display -->
      <div id="otp-timer-wrapper" style="text-align: center; margin-bottom: 1.25rem; background-color: #0a0a0a; padding: 0.6rem; border-radius: 6px; border: 1px solid var(--border-dark);">
        <div style="font-size: 0.75rem; color: var(--text-muted);">Reset OTP Expires In</div>
        <div id="reset-otp-timer-display" style="font-size: 1.5rem; font-weight: 800; color: var(--brand-primary); font-family: monospace; margin-top: 0.1rem;">10:00</div>
      </div>

      <form id="form-reset">
        <div class="form-group">
          <label class="form-label" for="reset-email">Email Address</label>
          <input type="email" id="reset-email" class="form-input" required value="${emailFromUrl}" placeholder="user@example.com" />
        </div>

        <div class="form-group">
          <label class="form-label" for="reset-otp">OTP Code</label>
          <input type="text" id="reset-otp" class="form-input" required placeholder="123456" maxlength="6" style="text-align: center; font-size: 1.25rem; letter-spacing: 0.25em; font-family: monospace;" />
        </div>

        <div class="form-group">
          <label class="form-label" for="reset-pass">New Password</label>
          <input type="password" id="reset-pass" class="form-input" required minlength="6" placeholder="At least 6 characters" />
        </div>

        <div id="auth-error" style="color: var(--brand-primary); font-size: 0.8125rem; margin-bottom: 1rem;" class="hidden"></div>
        <div id="auth-success" style="color: var(--accent-emerald); font-size: 0.8125rem; margin-bottom: 1rem;" class="hidden"></div>

        <button type="submit" class="btn btn-primary" style="width: 100%;">Reset Password</button>
      </form>

      <div style="margin-top: 1rem; text-align: center; display: flex; flex-direction: column; gap: 0.75rem; align-items: center;">
        <button id="btn-resend-reset-otp" class="btn btn-secondary btn-sm">Resend Reset OTP</button>
        <div style="font-size: 0.8125rem; color: var(--text-muted);">
          <a href="#/login" style="color: var(--brand-primary); font-weight: 600;">Back to Login</a>
        </div>
      </div>
    </div>
  `;

  // Timer Implementation (600 seconds = 10 minutes)
  function startResetOTPTimer(seconds) {
    if (otpTimerInterval) clearInterval(otpTimerInterval);

    let remaining = seconds;
    const timerDisplay = document.getElementById('reset-otp-timer-display');
    const resendBtn = document.getElementById('btn-resend-reset-otp');
    const errorDiv = document.getElementById('auth-error');

    if (timerDisplay) timerDisplay.style.color = 'var(--brand-primary)';

    function updateDisplay() {
      if (!timerDisplay) return;
      const mins = Math.floor(remaining / 60);
      const secs = remaining % 60;
      timerDisplay.textContent = `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;

      if (remaining <= 0) {
        clearInterval(otpTimerInterval);
        otpTimerInterval = null;
        timerDisplay.textContent = '00:00';
        timerDisplay.style.color = 'var(--text-muted)';
        if (errorDiv) {
          errorDiv.textContent = 'OTP code expired. Please click "Resend Reset OTP" for a new code.';
          errorDiv.classList.remove('hidden');
        }
      } else {
        remaining--;
      }
    }

    updateDisplay();
    otpTimerInterval = setInterval(updateDisplay, 1000);
  }

  // Start 10-minute countdown timer (600s) on render
  startResetOTPTimer(600);

  document.getElementById('form-reset').addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('reset-email').value.trim();
    const otp = document.getElementById('reset-otp').value.trim();
    const newPassword = document.getElementById('reset-pass').value;
    const errorDiv = document.getElementById('auth-error');
    const successDiv = document.getElementById('auth-success');
    errorDiv.classList.add('hidden');
    successDiv.classList.add('hidden');

    try {
      await API.post('/auth/reset-password', { email, otp, new_password: newPassword });
      if (otpTimerInterval) clearInterval(otpTimerInterval);
      successDiv.textContent = 'Password reset successfully! Redirecting to login...';
      successDiv.classList.remove('hidden');
      setTimeout(() => { window.location.hash = '#/login'; }, 1500);
    } catch (err) {
      errorDiv.textContent = err.message || 'Failed to reset password.';
      errorDiv.classList.remove('hidden');
    }
  });

  document.getElementById('btn-resend-reset-otp').addEventListener('click', async () => {
    const email = document.getElementById('reset-email').value.trim();
    const errorDiv = document.getElementById('auth-error');
    const successDiv = document.getElementById('auth-success');
    errorDiv.classList.add('hidden');
    successDiv.classList.add('hidden');

    if (!email) {
      errorDiv.textContent = 'Please enter your email address to resend reset OTP.';
      errorDiv.classList.remove('hidden');
      return;
    }

    try {
      await API.post('/auth/forgot-password', { email });
      successDiv.textContent = 'A new password reset OTP has been sent to your email.';
      successDiv.classList.remove('hidden');
      // Restart 10-minute timer from 10:00 (600s)
      startResetOTPTimer(600);
    } catch (err) {
      errorDiv.textContent = err.message || 'Failed to resend reset OTP.';
      errorDiv.classList.remove('hidden');
    }
  });
}

Router.addRoute('#/login', renderLogin);
Router.addRoute('#/signup', renderSignup);
Router.addRoute('#/verify-otp', renderVerifyOTP);
Router.addRoute('#/forgot-password', renderForgotPassword);
Router.addRoute('#/reset-password', renderResetPassword);
