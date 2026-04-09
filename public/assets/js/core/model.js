window.SensorPanel = window.SensorPanel || {};
window.SensorPanel.core = window.SensorPanel.core || {};

(function registerModel(core) {
  function pathParts() {
    return window.location.pathname.split("/").filter(Boolean);
  }

  function getDeviceNameFromPath() {
    const parts = pathParts();
    const idx = parts.indexOf("device");
    if (idx >= 0 && parts[idx + 1]) {
      return decodeURIComponent(parts[idx + 1]);
    }
    return null;
  }

  function componentTitle(component) {
    return component.options?.title || component.component || component.type || "component";
  }

  core.model = {
    pathParts,
    getDeviceNameFromPath,
    componentTitle
  };
})(window.SensorPanel.core);
