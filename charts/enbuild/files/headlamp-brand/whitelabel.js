/* ENBUILD white-label overlay for Headlamp — runs in the PAGE context
 * (nginx-injected <script>, NOT plugin code; per OPERATIONS_GOTCHAS §8
 * only page-context JS can touch document/localStorage reliably).
 *
 * Slimmed 2026-05-21 — the previous version (172 lines) injected an
 * ENBUILD top nav that visually duplicated Headlamp's own AppBar. The
 * native AppBar is now the only chrome (rebranded via the enbuild-theme
 * plugin's registerAppLogo + theme tokens), so this file is reduced to
 * the two things that genuinely need page context:
 *
 *   1. Pin document.title (Headlamp's RouteSwitcher resets it per cluster).
 *   2. Seed `ENBUILD Dark` as the default theme on first visit so the
 *      AppBar renders with brand colors before the user touches Settings.
 *
 * No DOM injection, no MutationObserver, no sidebar selectors to brittle
 * against Headlamp upgrades.
 */
(function () {
  'use strict';
  var BRAND_TITLE = 'ENBUILD Cluster Console';
  var THEME_KEY = 'headlampThemePreference';
  var THEME_VAL = 'ENBUILD Dark';

  // Default-theme seed (only if the operator hasn't chosen one). Headlamp's
  // lib/themes.ts:getThemeName reads this; null/empty falls back to its own
  // 'light' default. Setting our 'ENBUILD Dark' upfront avoids a flash of
  // the upstream theme on first paint.
  try {
    if (!localStorage.getItem(THEME_KEY)) localStorage.setItem(THEME_KEY, THEME_VAL);
  } catch (e) { /* localStorage disabled — registerAppTheme name='light'/'dark' overrides cover this */ }

  function pinTitle() {
    if (document.title !== BRAND_TITLE) document.title = BRAND_TITLE;
  }

  // Headlamp's RouteSwitcher rewrites document.title per cluster navigation
  // (e.g. "Headlamp / Clusters / ccm-vendor13-eks-final"). Re-pin on a cheap
  // interval — no MutationObserver needed.
  pinTitle();
  setInterval(pinTitle, 1500);
})();
