// ========================================
// Menu — all DOM refs are re-queried on
// every call so they survive HTMX swaps
// AND history restoration (which replaces
// the entire body's children).
// ========================================
(function () {
  function openMenu(toggle) {
    if (!toggle) toggle = document.querySelector('.menu-toggle');
    var menu = document.getElementById('nav-menu');
    var overlay = document.getElementById('nav-overlay');
    if (!toggle || !menu || !overlay) return;
    // Force reflow so the browser commits initial styles before the
    // transition starts — fixes first-open after HTMX DOM insertion.
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

  // All clicks delegated to document — survives any DOM replacement
  document.addEventListener('click', function (e) {
    var toggle = e.target.closest('.menu-toggle');
    if (toggle) {
      e.preventDefault();
      if (toggle.classList.contains('active')) {
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

  // Scroll to top after HTMX swaps
  document.body.addEventListener('htmx:afterSettle', function () {
    window.scrollTo(0, 0);
  });

  // Update nav highlight after URL is pushed (fires after afterSettle)
  document.body.addEventListener('htmx:pushedIntoHistory', updateNavHighlight);

  // Update nav highlight after HTMX restores from history cache
  document.body.addEventListener('htmx:historyRestore', updateNavHighlight);

  // Update document title from server response header
  document.body.addEventListener('htmx:afterRequest', function (e) {
    var title = e.detail.xhr.getResponseHeader('HX-Title');
    if (title) {
      document.title = title;
    }
  });
})();
