/* ENBUILD white-label overlay for Headlamp — runs in the PAGE context
 * (nginx-injected <script>, NOT plugin code; per OPERATIONS_GOTCHAS §8 only
 * page-context JS can touch document/localStorage reliably).
 *
 * Responsibilities:
 *   1. Seed `ENBUILD Dark` as the default theme on first visit.
 *   2. Pin document.title (Headlamp's RouteSwitcher resets it per cluster).
 *   3. Build the fixed ENBUILD top nav from the (CSS-hidden) native sidebar
 *      anchors, proxy-clicking them so Headlamp's own router does the nav.
 *   4. Keep all of the above enforced across SPA re-renders via a debounced
 *      MutationObserver, with a self-diagnostic if selectors miss.
 *
 * Defensive throughout: every selector has fallbacks; nothing throws. DOM is
 * built with createElement/textContent only (no innerHTML — no XSS surface).
 */
(function () {
  'use strict';
  var BRAND_TITLE = 'ENBUILD Cluster Console';
  var THEME_KEY = 'headlampThemePreference';
  var THEME_VAL = 'ENBUILD Dark';

  // 1. Default-theme seed (only if the operator hasn't chosen one).
  try {
    if (!localStorage.getItem(THEME_KEY)) localStorage.setItem(THEME_KEY, THEME_VAL);
  } catch (e) { /* localStorage disabled — theme plugin's name='dark' override still applies */ }

  function pinTitle() {
    if (document.title !== BRAND_TITLE) document.title = BRAND_TITLE;
  }

  function el(tag, cls, text) {
    var n = document.createElement(tag);
    if (cls) n.className = cls;
    if (text != null) n.textContent = text;
    return n;
  }

  // Find the native sidebar nav anchors. Headlamp's Drawer is an MUI Drawer;
  // entries are react-router <Link> => <a href="/c/<cluster>/...">. We keep
  // it in the DOM (CSS display:none) so these still route on .click().
  function sidebarAnchors() {
    var roots = document.querySelectorAll(
      '.MuiDrawer-root, .MuiDrawer-docked, [class*="Sidebar"]'
    );
    var seen = {}, out = [];
    for (var i = 0; i < roots.length; i++) {
      var as = roots[i].querySelectorAll('a[href]');
      for (var j = 0; j < as.length; j++) {
        var a = as[j];
        var href = a.getAttribute('href') || '';
        var label = (a.textContent || '').trim();
        if (!label || !href || href === '#' || seen[href]) continue;
        if (href.indexOf('/c/') === -1 && href.indexOf('/cluster') === -1) continue;
        seen[href] = 1;
        out.push({ a: a, href: href, label: label });
      }
    }
    return out;
  }

  function ensureBar() {
    var bar = document.getElementById('enbuild-topnav');
    if (bar) return bar;
    bar = el('header');
    bar.id = 'enbuild-topnav';

    var brand = el('div', 'eb-brand');
    var logo = el('img', 'eb-logo');
    logo.src = '/headlamp-brand/favicon.ico';
    logo.alt = 'ENBUILD';
    brand.appendChild(logo);
    brand.appendChild(el('span', 'eb-word', 'ENBUILD'));
    brand.appendChild(el('span', 'eb-sub', 'Cluster Console'));

    var nav = el('nav', 'eb-links');
    nav.setAttribute('aria-label', 'ENBUILD navigation');

    bar.appendChild(brand);
    bar.appendChild(nav);
    document.body.appendChild(bar);
    return bar;
  }

  function clear(node) {
    while (node.firstChild) node.removeChild(node.firstChild);
  }

  function buildNav() {
    var anchors = sidebarAnchors();
    if (!anchors.length) return false;
    var bar = ensureBar();
    var navEl = bar.querySelector('.eb-links');
    if (!navEl) return false;
    var sig = anchors.map(function (x) { return x.label; }).join('|');
    if (navEl.getAttribute('data-sig') === sig) {
      markActive(navEl);
      return true;
    }
    clear(navEl);
    anchors.forEach(function (x) {
      var item = el('a', 'eb-link', x.label);
      item.href = x.href;
      item.setAttribute('data-href', x.href);
      item.addEventListener('click', function (ev) {
        ev.preventDefault();
        // Proxy-click the real sidebar anchor → Headlamp's router handles it.
        x.a.click();
        setTimeout(function () { markActive(navEl); pinTitle(); }, 50);
      });
      navEl.appendChild(item);
    });
    navEl.setAttribute('data-sig', sig);
    markActive(navEl);
    return true;
  }

  function markActive(navEl) {
    var path = location.pathname;
    var items = navEl.querySelectorAll('.eb-link');
    for (var i = 0; i < items.length; i++) {
      var h = items[i].getAttribute('data-href') || '';
      items[i].classList.toggle('eb-active', !!h && path.indexOf(h) === 0);
    }
  }

  var missTicks = 0, warned = false;
  function tick() {
    pinTitle();
    var ok = buildNav();
    if (!ok) {
      if (++missTicks > 40 && !warned) {
        warned = true;
        console.warn(
          '[enbuild-whitelabel] native sidebar anchors not found after 40 ' +
          'ticks — Headlamp DOM may have changed; update selectors in ' +
          'whitelabel.js sidebarAnchors() for this Headlamp version.'
        );
      }
    } else {
      missTicks = 0;
    }
  }

  var raf = 0;
  function schedule() {
    if (raf) return;
    raf = window.requestAnimationFrame(function () { raf = 0; tick(); });
  }

  function start() {
    document.documentElement.classList.add('enbuild-whitelabel');
    pinTitle();
    try {
      new MutationObserver(schedule).observe(document.body, {
        childList: true, subtree: true,
      });
    } catch (e) { /* no MO — interval below still drives it */ }
    tick();
    // RouteSwitcher rewrites document.title per cluster nav; cheap re-pin.
    setInterval(pinTitle, 1500);
    // Safety net if the MutationObserver misses the first sidebar mount.
    var n = 0, iv = setInterval(function () {
      tick();
      if (++n > 30 || missTicks === 0) clearInterval(iv);
    }, 500);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', start);
  } else {
    start();
  }
})();
