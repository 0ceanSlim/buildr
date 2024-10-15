package utils

import (
	"html/template"
	"regexp"
)

// Function to convert image links in note content into <img> tags
func renderNoteContent(content string) template.HTML {
    // Regular expression to detect image links (e.g., ending with .png, .jpg, etc.)
    imageRegex := regexp.MustCompile(`(https?://[^\s]+(?:png|jpg|jpeg|gif))`)
    
    // Replace image links with <img> tags
    contentWithImages := imageRegex.ReplaceAllString(content, `<img src="$1" alt="Image" class="rounded-md note-image" />`)
    
    // Use template.HTML to mark content as safe HTML (be careful with user-generated content)
    return template.HTML(contentWithImages)
}