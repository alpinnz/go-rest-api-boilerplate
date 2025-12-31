package localization

// Dictionary represents language translation map.
// Uses nested map structure for flexibility and ease of use.
// Keys follow dot notation path (e.g., "validation.email.required").
type Dictionary map[string]any
