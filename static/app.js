// ========================================
// Menu — uses event delegation so it
// survives HTMX swaps of #page-wrapper
// ========================================
(function () {
  var menu = document.getElementById('nav-menu');
  var overlay = document.getElementById('nav-overlay');

  function getToggle() {
    return document.querySelector('.menu-toggle');
  }

  function openMenu() {
    var toggle = getToggle();
    if (!toggle || !menu || !overlay) return;
    toggle.classList.add('active');
    menu.classList.add('active');
    overlay.classList.add('active');
    toggle.setAttribute('aria-expanded', 'true');
    toggle.setAttribute('aria-label', 'Close menu');
    var firstLink = menu.querySelector('a');
    if (firstLink) firstLink.focus();
  }

  function closeMenu() {
    var toggle = getToggle();
    if (toggle) {
      toggle.classList.remove('active');
      toggle.setAttribute('aria-expanded', 'false');
      toggle.setAttribute('aria-label', 'Open menu');
    }
    if (menu) menu.classList.remove('active');
    if (overlay) overlay.classList.remove('active');
  }

  // Delegate click on toggle — works even after DOM replacement
  document.addEventListener('click', function (e) {
    if (e.target.closest('.menu-toggle')) {
      e.preventDefault();
      var toggle = getToggle();
      if (toggle && toggle.classList.contains('active')) {
        closeMenu();
      } else {
        openMenu();
      }
    }
  });

  // Overlay click closes menu
  if (overlay) {
    overlay.addEventListener('click', closeMenu);
  }

  // Escape key closes menu
  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') {
      closeMenu();
    }
  });

  // Close menu before HTMX navigations
  document.body.addEventListener('htmx:beforeRequest', closeMenu);

  // Scroll to top after HTMX swaps
  document.body.addEventListener('htmx:afterSettle', function () {
    window.scrollTo(0, 0);
  });

  // Update document title from server response header
  document.body.addEventListener('htmx:afterRequest', function (e) {
    var title = e.detail.xhr.getResponseHeader('HX-Title');
    if (title) {
      document.title = title;
    }
  });

  // Browser back/forward — re-fetch the page content
  window.addEventListener('popstate', function () {
    var wrapper = document.getElementById('page-wrapper');
    if (wrapper) {
      htmx.ajax('GET', window.location.pathname, { target: '#page-wrapper', swap: 'innerHTML' });
    }
  });
})();