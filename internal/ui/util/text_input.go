package util

import "github.com/rivo/tview"

var activeTextInputs = map[tview.Primitive]bool{}

// SetTextInputActive marks the given primitive as capturing text input (or not).
func SetTextInputActive(primitive tview.Primitive, active bool) {
	if active {
		activeTextInputs[primitive] = true
	} else {
		delete(activeTextInputs, primitive)
	}
}

// IsTextInputActive returns whether the given primitive captures text input.
func IsTextInputActive(primitive tview.Primitive) bool {
	if primitive == nil {
		return false
	}
	if _, ok := primitive.(*tview.InputField); ok {
		return true
	}
	return activeTextInputs[primitive]
}
