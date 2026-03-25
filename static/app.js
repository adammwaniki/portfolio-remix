// ========================================
// Menu — all DOM refs re-queried every call
// (survives HTMX swaps + history restore).
// The nav-menu is the source of truth for
// open/closed state because it persists
// outside #page-wrapper, while the toggle
// gets replaced on every HTMX swap.
// ========================================
(function () {
  function isMenuOpen() {
    var menu = document.getElementById('nav-menu');
    return menu && menu.classList.contains('active');
  }

  function openMenu(toggle) {
    if (!toggle) toggle = document.querySelector('.menu-toggle');
    var menu = document.getElementById('nav-menu');
    var overlay = document.getElementById('nav-overlay');
    if (!toggle || !menu || !overlay) return;
    toggle.offsetHeight;
    toggle.classList.add('active');
    menu.classList.add('active');
    overlay.classList.add('active');
    toggle.setAttribute('aria-expanded', 'true');
    toggle.setAttribute('aria-label', 'Close menu');
    var firstLink = menu.querySelector('a');
    if (firstLink) firstLink.focus();
  }

  function closeMenu() {
    var toggle = document.querySelector('.menu-toggle');
    var menu = document.getElementById('nav-menu');
    var overlay = document.getElementById('nav-overlay');
    if (toggle) {
      toggle.classList.remove('active');
      toggle.setAttribute('aria-expanded', 'false');
      toggle.setAttribute('aria-label', 'Open menu');
    }
    if (menu) menu.classList.remove('active');
    if (overlay) overlay.classList.remove('active');
  }

  // After HTMX replaces #page-wrapper the new toggle has no classes.
  // Sync it with the menu so the hamburger/X visual matches reality.
  function syncToggleState() {
    var toggle = document.querySelector('.menu-toggle');
    var menu = document.getElementById('nav-menu');
    if (!toggle || !menu) return;
    if (menu.classList.contains('active')) {
      toggle.classList.add('active');
      toggle.setAttribute('aria-expanded', 'true');
      toggle.setAttribute('aria-label', 'Close menu');
    } else {
      toggle.classList.remove('active');
      toggle.setAttribute('aria-expanded', 'false');
      toggle.setAttribute('aria-label', 'Open menu');
    }
  }

  // All clicks delegated to document — survives any DOM replacement.
  // Uses menu state (persistent) as source of truth, not toggle state.
  document.addEventListener('click', function (e) {
    var toggle = e.target.closest('.menu-toggle');
    if (toggle) {
      e.preventDefault();
      if (isMenuOpen()) {
        closeMenu();
      } else {
        openMenu(toggle);
      }
      return;
    }
    if (e.target.closest('#nav-overlay')) {
      closeMenu();
    }
  });

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') closeMenu();
  });

  function updateNavHighlight() {
    var currentPath = window.location.pathname;
    var links = document.querySelectorAll('.nav-menu a');
    for (var i = 0; i < links.length; i++) {
      var href = links[i].getAttribute('href');
      var isCurrent = href === currentPath || (href !== '/' && currentPath.startsWith(href + '/'));
      if (isCurrent) {
        links[i].classList.add('current');
        links[i].setAttribute('aria-current', 'page');
      } else {
        links[i].classList.remove('current');
        links[i].removeAttribute('aria-current');
      }
    }
  }

  // Close menu before HTMX navigations
  document.body.addEventListener('htmx:beforeRequest', closeMenu);

  // After HTMX swaps #page-wrapper: sync the new toggle with menu state
  document.body.addEventListener('htmx:afterSettle', function () {
    syncToggleState();
    window.scrollTo(0, 0);
  });

  // Update nav highlight after URL is pushed
  document.body.addEventListener('htmx:pushedIntoHistory', updateNavHighlight);

  // After history restore: ensure clean state
  document.body.addEventListener('htmx:historyRestore', function () {
    closeMenu();
    updateNavHighlight();
  });

  // Update document title from server response header
  document.body.addEventListener('htmx:afterRequest', function (e) {
    var title = e.detail.xhr.getResponseHeader('HX-Title');
    if (title) {
      document.title = title;
    }
  });
})();
