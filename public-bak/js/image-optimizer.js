/**
 * Image optimization utilities for SpurtCMS
 * Helps reduce CPU usage during image cropping operations
 */

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

// Show loading indicator
function showLoadingIndicator() {
  // Check if loading indicator exists, if not create it
  if ($('#image-processing-indicator').length === 0) {
    $('body').append(
      '<div id="image-processing-indicator" style="display:none; position:fixed; top:0; left:0; width:100%; height:100%; background:rgba(0,0,0,0.5); z-index:9999;">' +
      '<div style="position:absolute; top:50%; left:50%; transform:translate(-50%,-50%); background:white; padding:20px; border-radius:5px; text-align:center;">' +
      '<div class="spinner-border text-primary" role="status"></div>' +
      '<p style="margin-top:10px;">Processing image...</p>' +
      '</div></div>'
    );
  }
  $('#image-processing-indicator').fadeIn(200);
}

// Hide loading indicator
function hideLoadingIndicator() {
  $('#image-processing-indicator').fadeOut(200);
}
