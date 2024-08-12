function bindCopyListener ({ elementId, contentToCopyId }) {
  const copyButton = document.getElementById(elementId);

  copyButton.addEventListener('click', function() {
    const contentToCopy = contentToCopyId 
      ? document.getElementById(contentToCopyId).textContent
      : copyButton.getAttribute('copy-target');

    navigator.clipboard.writeText(contentToCopy)
      .then(() => {
        const copyButton = document.getElementById(elementId);
        const originalText = copyButton.textContent;
        copyButton.textContent = 'Copied!';
        setTimeout(() => copyButton.textContent = originalText, 2000); // Revert after 2 seconds
      })
      .catch(err => {
        console.error('Failed to copy text: ', err);
      });
  });  
}
