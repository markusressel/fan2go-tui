package util

import (
	"github.com/rivo/tview"
	"reflect"
)

func IsTxViewVisible(view tview.Primitive) bool {
	var viewToCheck interface{}
	switch view.(type) {
	case *tview.Box:
		box := view.(*tview.Box)
		_, _, width, height := box.GetInnerRect()
		return width <= 0 && height <= 0
	case *tview.Flex:
		viewToCheck = view.(*tview.Flex).Box
	case *tview.Grid:
		viewToCheck = view.(*tview.Grid).Box
	case *tview.Frame:
		viewToCheck = view.(*tview.Frame).Box
	case *tview.List:
		viewToCheck = view.(*tview.List).Box
	default:
		viewToCheck = view
	}

	innerWidth := reflect.ValueOf(viewToCheck).Elem().FieldByName("innerWidth").Int()
	innerHeight := reflect.ValueOf(viewToCheck).Elem().FieldByName("innerHeight").Int()
	return innerWidth <= 0 && innerHeight <= 0
}
