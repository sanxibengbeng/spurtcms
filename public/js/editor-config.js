/**
 * SpurtCMS Editor Configuration
 * This file configures the image upload URL for the built-in editor
 */

// Get the base URL from the environment or use the current domain
function getBaseUrl() {
  // If we're in a browser environment
  if (typeof window !== 'undefined') {
    // Get the current domain and protocol
    const protocol = window.location.protocol;
    const host = window.location.host;
    return `${protocol}//${host}/`;
  }
  
  // Default fallback (should not be used in production)
  return 'http://localhost:8082/';
}

// Editor configuration object
const editorConfig = {
  imageUpload: {
    path: getBaseUrl() + "uploadb64image",
    payload: "imagedata",
    success: {
      imagepath: "imagepath",
      imagename: "imagename"
    }
  }
};

// Make the configuration available globally
if (typeof window !== 'undefined') {
  window.spurtEditorConfig = editorConfig;
}

// Export for module environments
if (typeof module !== 'undefined' && module.exports) {
  module.exports = editorConfig;
}
