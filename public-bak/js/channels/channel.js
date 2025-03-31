/** Gobal Variables */
var languagedata
let sections = []
var fiedlvalue = []
let optionvalues = []
let deleteoption = []
let deletesecion = []
let deletefields = []
let fielid = 0             //db field id
let newFieldid = 1        //temporary field id
let orderindex = 1        //order index field
let optionid = 1
var sectionid = 1
let RelationalMember = 14
let RelationalMedia = 15
let RelationalVideo = 16
let SelectedCategoryValue = []  // Initialize SelectedCategoryValue array
let isEventProcessing = false   // Flag to prevent recursive event handling

// Add custom validation method for space checking
if ($.validator && !$.validator.methods.space) {
    $.validator.addMethod("space", function (value) {
        return /^[^-\s]/.test(value);
    }, "Cannot start with space or hyphen");
}

$(document).ready(async function () {
  var languagepath = $('.language-group>button').attr('data-path')

  await $.getJSON(languagepath, function (data) {
    languagedata = data
  })

  $('.drag-btn').attr({
    'data-bs-toggle': 'tooltip',
    'data-bs-placement': 'bottom',
    'data-bs-custom-class': 'custom-tooltip',
    'data-bs-title': 'Drag & Drop'
  });

  $('.edit-field').attr({
    'data-bs-toggle': 'tooltip',
    'data-bs-placement': 'bottom',
    'data-bs-custom-class': 'custom-tooltip',
    'data-bs-title': 'Edit'
  });

  $('.duplicate-field').attr({
    'data-bs-toggle': 'tooltip',
    'data-bs-placement': 'bottom',
    'data-bs-custom-class': 'custom-tooltip',
    'data-bs-title': 'Duplicate'
  });

  $('.delete-field').attr({
    'data-bs-toggle': 'tooltip',
    'data-bs-placement': 'bottom',
    'data-bs-custom-class': 'custom-tooltip',
    'data-bs-title': 'Delete'
  });

  const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]')
  const tooltipList = [...tooltipTriggerList].map(tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl))

  // Add event handlers for category selection buttons
  $(document).on("click", ".category-select-btn", function(e) {
    // Prevent recursive calls
    if (isEventProcessing) {
      e.stopPropagation();
      return;
    }
    
    isEventProcessing = true;
    
    // Get category ID
    var categoryid = $(this).attr("data-categoryid");
    
    // Add to selected categories if not already there
    if (!SelectedCategoryValue.includes(categoryid)) {
      SelectedCategoryValue.push(categoryid);
    }
    
    // Update UI or perform other actions
    console.log("Selected categories:", SelectedCategoryValue);
    
    isEventProcessing = false;
  });

  // Add event handlers for category unselection buttons
  $(document).on("click", ".category-unselect-btn", function(e) {
    // Prevent recursive calls
    if (isEventProcessing) {
      e.stopPropagation();
      return;
    }
    
    isEventProcessing = true;
    
    // Get category ID
    var categoryid = $(this).attr("data-categoryid");
    
    // Remove from selected categories
    SelectedCategoryValue = SelectedCategoryValue.filter(function(value) {
      return value !== categoryid;
    });
    
    // Update UI or perform other actions
    console.log("Selected categories after removal:", SelectedCategoryValue);
    
    isEventProcessing = false;
  });

  // Modified event handlers for category div clicks
  $(document).on("click", ".categorypdiv", function(e) {
    // Prevent recursive calls
    if (isEventProcessing) {
      e.stopPropagation();
      return;
    }
    
    isEventProcessing = true;
    var btn = $(this).find('.category-select-btn');
    
    // Get category ID
    var categoryid = btn.attr("data-categoryid");
    
    // Add to selected categories if not already there
    if (!SelectedCategoryValue.includes(categoryid)) {
      SelectedCategoryValue.push(categoryid);
    }
    
    // Update UI or perform other actions
    console.log("Selected categories from div click:", SelectedCategoryValue);
    
    isEventProcessing = false;
  });

  $(document).on("click", ".selectedcategorydiv", function(e) {
    // Prevent recursive calls
    if (isEventProcessing) {
      e.stopPropagation();
      return;
    }
    
    isEventProcessing = true;
    var btn = $(this).find('.category-unselect-btn');
    
    // Get category ID
    var categoryid = btn.attr("data-categoryid");
    
    // Remove from selected categories
    SelectedCategoryValue = SelectedCategoryValue.filter(function(value) {
      return value !== categoryid;
    });
    
    // Update UI or perform other actions
    console.log("Selected categories after removal from div click:", SelectedCategoryValue);
    
    isEventProcessing = false;
  });
});
