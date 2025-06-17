package ui

import "github.com/manifoldco/promptui"

// UI Templates for consistent styling across the application
var (
	// StandardPromptTemplates is the default template for prompt inputs
	StandardPromptTemplates = &promptui.PromptTemplates{
		Prompt:  "{{ . | cyan | bold }}: ",
		Valid:   "{{ . | green | bold }}: ",
		Invalid: "{{ . | red | bold }}: ",
		Success: "{{ . | bold | green }}: ",
	}

	// StandardSelectTemplates is the default template for select menus
	StandardSelectTemplates = &promptui.SelectTemplates{
		Label:    "{{ . | cyan | bold }}",
		Active:   "▸ {{ . | cyan }}",
		Inactive: "  {{ . | faint }}",
		Selected: "{{ . | green | bold }}",
	}

	// ItemSelectTemplates is for select menus with structured items (with .Display field)
	ItemSelectTemplates = &promptui.SelectTemplates{
		Label:    "{{ . | cyan | bold }}",
		Active:   "{{ \"▶\" | cyan | bold }} {{ .Display | cyan }}",
		Inactive: "  {{ .Display | faint }}",
		Selected: "{{ \"✔\" | green | bold }} {{ .Display | green }}",
	}

	// DetailedSelectTemplates is for select menus with additional details
	DetailedSelectTemplates = &promptui.SelectTemplates{
		Label:    "{{ . | cyan | bold }}",
		Active:   "{{ \"▶\" | cyan | bold }} {{ .Display | cyan }}",
		Inactive: "  {{ .Display | faint }}",
		Selected: "{{ \"✔\" | green | bold }} {{ .Display | green }}",
		Details: `
--------- Details ----------
{{ "Title:" | faint }}	{{ .Title }}
{{ "URL:" | faint }}	{{ .URL | cyan }}
{{ "Category:" | faint }}	{{ .Category }}
{{ "Created:" | faint }}	{{ .CreatedAt.Format "2006-01-02 15:04:05" }}`,
	}

	// SimpleSelectTemplates is for basic select menus without styling
	SimpleSelectTemplates = &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ . }}",
		Inactive: "  {{ . }}",
		Selected: "{{ . }}",
	}
)

// GetSelectSize returns appropriate size for select prompt based on item count
func GetSelectSize(itemCount int) int {
	const maxSize = 10
	if itemCount < maxSize {
		return itemCount
	}
	return maxSize
}