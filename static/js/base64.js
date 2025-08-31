(function () {
  function setButtonText(btn, text, timeout) {
    if (!btn) return;
    var original = btn.dataset.origText || btn.textContent || 'COPY';
    btn.textContent = text;
    clearTimeout(btn._copyTimer);
    btn._copyTimer = setTimeout(function () {
      btn.textContent = original;
    }, timeout || 1200);
  }

  function copyToClipboard(btn) {
    try {
      var el = document.getElementById('b64-result');
      var value = '';
      if (el) {
        value = (typeof el.value === 'string' && el.value.length ? el.value : (el.textContent || '')) || '';
      }
      navigator.clipboard.writeText(value)
        .then(function () { setButtonText(btn, 'COPIED', 1200); })
        .catch(function () { setButtonText(btn, 'COPY FAILED', 1500); });
    } catch (e) {
      setButtonText(btn, 'COPY FAILED', 1500);
    }
  }

  function updateVisibility() {
    var el = document.getElementById('b64-result');
    var text = '';
    if (el) {
      text = (typeof el.value === 'string' && el.value.length ? el.value : (el.textContent || '')) || '';
    }
    var show = !!el && text.length > 0;
    ['copy-top', 'copy-top-decode'].forEach(function (id) {
      var btn = document.getElementById(id);
      if (!btn) return;
      if (show) { btn.classList.remove('hidden'); } else { btn.classList.add('hidden'); }
    });
  }

  function attachHandlers() {
    ['copy-top', 'copy-top-decode'].forEach(function (id) {
      var btn = document.getElementById(id);
      if (!btn) return;
      if (!btn.dataset.origText) btn.dataset.origText = btn.textContent || 'COPY';
      if (!btn._copyBound) {
        btn.addEventListener('click', function () { copyToClipboard(btn); });
        btn._copyBound = true;
      }
    });
    updateVisibility();
  }

  document.addEventListener('DOMContentLoaded', attachHandlers);
  document.addEventListener('htmx:afterSwap', function (evt) {
    if (evt && evt.target && evt.target.id === 'response') {
      attachHandlers();
    }
  });
})();

