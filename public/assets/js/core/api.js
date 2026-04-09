window.SensorPanel = window.SensorPanel || {};
window.SensorPanel.core = window.SensorPanel.core || {};

(function registerApi(core) {
  async function fetchJSON(url) {
    const resp = await fetch(url, { cache: "no-store" });
    if (!resp.ok) {
      throw new Error(`HTTP ${resp.status}`);
    }
    return await resp.json();
  }

  async function fetchDeviceModel(deviceName) {
    return await fetchJSON(`/api/device/${encodeURIComponent(deviceName)}`);
  }

  core.api = {
    fetchJSON,
    fetchDeviceModel
  };
})(window.SensorPanel.core);
