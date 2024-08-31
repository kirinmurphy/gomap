(function () {
  const validateUrl = (url) => {
    const regex = /^https:\/\/docs\.google\.com\/spreadsheets\/d\/e\/([a-zA-Z0-9-_]+)\/pubhtml$/;
    return regex.test(url) ? url.match(regex)[1] : null;
  };

  const displayError = (message) => {
    const errorElement = document.getElementById("urlError");
    errorElement.textContent = message;
    errorElement.classList.remove("hidden");
  };

  const hideError = () => {
    document.getElementById("urlError").classList.add("hidden");
  };

  document.getElementById("sheetUrl").addEventListener("paste", (event) => {
    setTimeout(() => {
      const url = event.target.value;
      const sheetId = validateUrl(url);
      if (!sheetId) {
        displayError("I think you got a bad URL there bud.");
        console.log('SHEETID', sheetId);

        return;
      } else {
        hideError();
      }

      document.getElementById("loadingIndicator").classList.remove("hidden");

      htmx.ajax('POST', '/updateMapUI', { 
        values: { sheetId },
        target: '#mapContainer',
        swap: 'innerHtml'
      });
      
    }, 0);
  });

  document.getElementById("sheetForm").addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
      event.preventDefault(); 
    }
  });

  document.body.addEventListener("htmx:afterSwap", (event) => {
    setTimeout(() => {
      const mapContainer = document.getElementById("mapContainer");
      document.getElementById("loadingIndicator").classList.add("hidden");
      if (mapContainer) {
        const topPosition = mapContainer.getBoundingClientRect().top + window.scrollY - 10;
        window.scrollTo({
          top: topPosition,
          behavior: "smooth"
        });
      }
    }, 50);
  });

  const demoPrompt = document.getElementById("demoPrompt");
  if (demoPrompt) {
    const pasteDemoUrlButton = document.getElementById("pasteDemoUrl");
    const sheetUrlInput = document.getElementById("sheetUrl");
    pasteDemoUrlButton.addEventListener("click", () => {
      const urlToPaste = pasteDemoUrlButton.getAttribute("text-to-copy");
      sheetUrlInput.value = urlToPaste;
      sheetUrlInput.dispatchEvent(new Event("paste"));
    });
  }
})()