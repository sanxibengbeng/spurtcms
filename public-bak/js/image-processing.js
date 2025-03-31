// Image processing functions for SpurtCMS

// Throttle function to limit how often a function can be called
function throttle(func, limit) {
  let inThrottle;
  return function() {
    const args = arguments;
    const context = this;
    if (!inThrottle) {
      func.apply(context, args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, limit);
    }
  };
}

// Resize large images before passing them to the cropper
function resizeImageIfNeeded(img, maxWidth, maxHeight) {
  if (img.width > maxWidth || img.height > maxHeight) {
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    
    // Calculate new dimensions while maintaining aspect ratio
    let newWidth = img.width;
    let newHeight = img.height;
    
    if (newWidth > maxWidth) {
      newHeight = (maxWidth / newWidth) * newHeight;
      newWidth = maxWidth;
    }
    
    if (newHeight > maxHeight) {
      newWidth = (maxHeight / newHeight) * newWidth;
      newHeight = maxHeight;
    }
    
    canvas.width = newWidth;
    canvas.height = newHeight;
    
    // Draw resized image to canvas
    ctx.drawImage(img, 0, 0, newWidth, newHeight);
    
    // Return new data URL
    return canvas.toDataURL('image/jpeg', 0.7);
  }
  return null; // No resize needed
}

// Function to show loading indicator
function showLoadingIndicator() {
  // Add loading indicator to the page
  if (!$('#loading-indicator').length) {
    $('body').append('<div id="loading-indicator" style="position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.5); z-index: 9999; display: flex; justify-content: center; align-items: center;"><div style="background-color: white; padding: 20px; border-radius: 5px;"><p>Processing image...</p></div></div>');
  } else {
    $('#loading-indicator').show();
  }
}

// Function to hide loading indicator
function hideLoadingIndicator() {
  $('#loading-indicator').hide();
}
