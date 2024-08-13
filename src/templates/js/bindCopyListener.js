function bindCopyListener (props) {
  const { 
    elementId, 
    contentToCopyId, 
    copyMessage = 'Copied!' 
  } = props;

  const copyButton = document.getElementById(elementId);

  copyButton.addEventListener('click', function() {
    const contentToCopy = contentToCopyId 
      ? document.getElementById(contentToCopyId).textContent
      : copyButton.getAttribute('text-to-copy');

    navigator.clipboard.writeText(contentToCopy)
      .then(() => {
        const copyButton = document.getElementById(elementId);
        const originalText = copyButton.innerHTML;
        copyButton.innerHTML = copyMessage;
        setTimeout(() => copyButton.innerHTML = originalText, 2000); 
      })
      .catch(err => {
        console.error('Failed to copy text: ', err);
      });
  });  
}
