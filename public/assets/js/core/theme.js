window.SensorPanel = window.SensorPanel || {};
window.SensorPanel.core = window.SensorPanel.core || {};

(function registerTheme(core) {
  function applyThemeTokens(tokens) {
    if (!tokens) return;
    const root = document.documentElement;
    for (const [key, value] of Object.entries(tokens)) {
      root.style.setProperty(`--${key.replace(/_/g, "-")}`, value);
    }
  }

  core.theme = {
    applyThemeTokens
  };
})(window.SensorPanel.core);
