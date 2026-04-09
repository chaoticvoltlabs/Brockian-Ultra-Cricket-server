window.SensorPanel = window.SensorPanel || {};
window.SensorPanel.core = window.SensorPanel.core || {};

(function registerDom(core) {
  function el(tag, className, text) {
    const node = document.createElement(tag);
    if (className) node.className = className;
    if (text !== undefined && text !== null) node.textContent = text;
    return node;
  }

  core.dom = {
    el
  };
})(window.SensorPanel.core);
