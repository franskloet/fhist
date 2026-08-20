package shell

import "fmt"

func GetInitScript(shellType string) (string, error) {
	switch shellType {
	case "bash":
		return `# fhist bash integration (Ctrl+R replacement)
fhist_widget() {
    local selected_cmd
    selected_cmd=$(fhist --context 5)
    if [ -n "$selected_cmd" ]; then
        READLINE_LINE="$selected_cmd"
        READLINE_POINT=${#READLINE_LINE}
    fi
}
bind -x '"\C-r": fhist_widget'
`, nil

	case "zsh":
		return `# fhist zsh integration (Ctrl+R replacement)
fhist_widget() {
    local selected_cmd
    selected_cmd=$(fhist --context 5)
    if [ -n "$selected_cmd" ]; then
        BUFFER="$selected_cmd"
        CURSOR=${#BUFFER}
    fi
    zle reset-prompt
}
zle -N fhist_widget
bindkey '^R' fhist_widget
`, nil

	case "fish":
		return `# fhist fish integration (Ctrl+R replacement)
function fhist_widget
    set -l selected_cmd (fhist --context 5)
    if test -n "$selected_cmd"
        commandline -r "$selected_cmd"
    end
    commandline -f repaint
end
bind \cr fhist_widget
`, nil

	default:
		return "", fmt.Errorf("unsupported shell '%s'. Supported shells: bash, zsh, fish", shellType)
	}
}
