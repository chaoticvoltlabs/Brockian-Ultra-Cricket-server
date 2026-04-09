window.SensorPanel = window.SensorPanel || {};
window.SensorPanel.core = window.SensorPanel.core || {};

(function registerFormat(core) {
  function formatFixed(value, decimals, fallback = "—") {
    const num = Number(value);
    if (!Number.isFinite(num)) return fallback;
    return num.toFixed(decimals);
  }

  function formatInt(value, fallback = "—") {
    const num = Number(value);
    if (!Number.isFinite(num)) return fallback;
    return String(Math.round(num));
  }

  core.format = {
    formatFixed,
    formatInt
  };
})(window.SensorPanel.core);
