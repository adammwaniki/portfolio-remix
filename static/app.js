(function () {
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

  document.addEventListener('click', function (e) {
    // Hamburger toggle
    var toggle = e.target.closest('.menu-toggle');
    if (toggle) {
      var menu = document.getElementById('nav-menu');
      var overlay = document.getElementById('nav-overlay');
      if (!menu || !overlay) return;

      // Menu is the source of truth (it persists, toggle gets replaced by HTMX)
      if (menu.classList.contains('active')) {
        closeMenu();
      } else {
        toggle.classList.add('active');
        toggle.setAttribute('aria-expanded', 'true');
        toggle.setAttribute('aria-label', 'Close menu');
        menu.classList.add('active');
        overlay.classList.add('active');
      }
      return;
    }

    // Overlay click closes menu
    if (e.target.closest('#nav-overlay')) {
      closeMenu();
    }
  });

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') closeMenu();
  });

  function updateNavHighlight() {
    var path = window.location.pathname;
    var links = document.querySelectorAll('.nav-menu a');
    for (var i = 0; i < links.length; i++) {
      var href = links[i].getAttribute('href');
      if (href === path || (href !== '/' && path.startsWith(href + '/'))) {
        links[i].classList.add('current');
      } else {
        links[i].classList.remove('current');
      }
    }
  }

  document.body.addEventListener('htmx:beforeRequest', closeMenu);
  document.body.addEventListener('htmx:afterSettle', function () { window.scrollTo(0, 0); });
  document.body.addEventListener('htmx:pushedIntoHistory', updateNavHighlight);
  document.body.addEventListener('htmx:historyRestore', function () { closeMenu(); updateNavHighlight(); });
  document.body.addEventListener('htmx:afterRequest', function (e) {
    var title = e.detail.xhr.getResponseHeader('HX-Title');
    if (title) document.title = title;
  });
})();
