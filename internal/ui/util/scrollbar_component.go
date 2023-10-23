package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"math"
)

type ScrollBarOrientation int

const (
	ScrollBarVertical ScrollBarOrientation = iota
	ScrollBarHorizontal
)

type ScrollbarComponent struct {
	application *tview.Application

	layout       *tview.Flex
	topArrow     *tview.TextView
	upperBox     *tview.Box
	scrollBarBox *tview.Box
	lowerBox     *tview.Box
	bottomArrow  *tview.TextView

	inputCapture func(event *tcell.EventKey) *tcell.EventKey

	orientation    ScrollBarOrientation
	scrollPosition int
	barWidth       int
	min            int
	max            int
}

func NewScrollbarComponent(
	application *tview.Application,
	orientation ScrollBarOrientation,
	min int,
	max int,
	scrollPosition int,
	barWidth int,
) *ScrollbarComponent {
	scrollbarComponent := &ScrollbarComponent{
		application: application,
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		min:            min,
		max:            max,
		scrollPosition: scrollPosition,
		barWidth:       barWidth,
	}
	scrollbarComponent.createLayout()
	scrollbarComponent.SetOrientation(orientation)
	return scrollbarComponent
}

func (c *ScrollbarComponent) createLayout() {
	layout := tview.NewFlex()

	c.topArrow = tview.NewTextView()
	layout.AddItem(c.topArrow, 1, 0, false)
	c.upperBox = tview.NewBox()
	c.upperBox.SetBackgroundColor(tcell.ColorGray)
	layout.AddItem(c.upperBox, 1, 0, false)
	c.scrollBarBox = tview.NewBox()
	c.scrollBarBox.SetBackgroundColor(tcell.ColorWhite)
	layout.AddItem(c.scrollBarBox, 1, 0, false)
	c.lowerBox = tview.NewBox()
	c.lowerBox.SetBackgroundColor(tcell.ColorGray)
	layout.AddItem(c.lowerBox, 1, 0, false)
	c.bottomArrow = tview.NewTextView()
	layout.AddItem(c.bottomArrow, 1, 0, false)

	c.layout = layout
}

func (c *ScrollbarComponent) updateLayout() {
	// update the top arrow
	c.layout.ResizeItem(c.topArrow, 1, 0)

	c.updateTopEndText()
	c.updateScrollbar()
	c.updateBottomEndText()

	c.application.ForceDraw()
}

func (c *ScrollbarComponent) SetOrientation(orientation ScrollBarOrientation) {
	c.orientation = orientation
	switch c.orientation {
	case ScrollBarVertical:
		c.layout.SetDirection(tview.FlexRow)
	case ScrollBarHorizontal:
		c.layout.SetDirection(tview.FlexColumn)
	}
	c.updateLayout()
}

func (c *ScrollbarComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *ScrollbarComponent) SetTitle(title string) {
	SetupWindow(c.layout, title)
}

func (c *ScrollbarComponent) GetMin() int {
	return c.min
}

func (c *ScrollbarComponent) GetMax() int {
	return c.max
}

func (c *ScrollbarComponent) GetPosition() int {
	return c.scrollPosition
}

func (c *ScrollbarComponent) SetMin(min int) {
	c.min = min
	c.updateLayout()
}

func (c *ScrollbarComponent) SetMax(max int) {
	c.max = max
	c.updateLayout()
}

func (c *ScrollbarComponent) SetPosition(position int) {
	if position < 0 {
		position = 0
	}
	c.scrollPosition = position
	c.updateLayout()
}

func (c *ScrollbarComponent) HasFocus() bool {
	return c.layout.HasFocus()
}

func (c *ScrollbarComponent) SetInputCapture(inputCapture func(event *tcell.EventKey) *tcell.EventKey) {
	c.inputCapture = inputCapture
}

func (c *ScrollbarComponent) scrollUp() {
	c.scroll(-1)
}

func (c *ScrollbarComponent) scrollDown() {
	c.scroll(+1)
}

// scroll moves the scrollbar to the specified position
func (c *ScrollbarComponent) scroll(amount int) {
	oldPosition := c.GetPosition()
	newPosition := oldPosition + amount
	if newPosition < c.GetMin() {
		newPosition = c.GetMin()
	}
	if newPosition > c.GetMax() {
		newPosition = c.GetMax()
	}
	c.SetPosition(newPosition)

	newBarWidth := c.calculateBarWidth()
	c.barWidth = newBarWidth

	c.updateLayout()
}

func (c *ScrollbarComponent) ScrollToTop() {
	c.scroll(-c.GetPosition())
}

func (c *ScrollbarComponent) calculateBarWidth() int {
	// calculate the bar width
	barWidth := int(math.Max(1, float64(MaxVisibleItems/(c.max-c.min))))
	return barWidth
}

func (c *ScrollbarComponent) updateTopEndText() {
	c.layout.ResizeItem(c.topArrow, 1, 0)
	text := ""
	//isAtLimit := c.scrollPosition <= c.min
	isAtLimit := false
	switch c.orientation {
	case ScrollBarVertical:
		if isAtLimit {
			text = "="
		} else {
			text = "^"
		}
	case ScrollBarHorizontal:
		if isAtLimit {
			text = "="
		} else {
			text = "<"
		}
	}
	c.topArrow.SetText(text)
}

func (c *ScrollbarComponent) updateBottomEndText() {
	c.layout.ResizeItem(c.bottomArrow, 1, 0)
	text := ""
	//isAtLimit := c.scrollPosition+c.barWidth >= c.max
	isAtLimit := false
	switch c.orientation {
	case ScrollBarVertical:
		if isAtLimit {
			text = "="
		} else {
			text = "v"
		}
	case ScrollBarHorizontal:
		if isAtLimit {
			text = "="
		} else {
			text = ">"
		}
	}
	c.bottomArrow.SetText(text)
}

func (c *ScrollbarComponent) updateScrollbar() {
	// calculate the box sizes
	upperBoxStart := float64(c.min)
	upperBoxEnd := float64(c.min + c.scrollPosition)
	scrollBarBoxStart := upperBoxEnd
	scrollBarBoxEnd := scrollBarBoxStart + float64(c.barWidth)
	lowerBoxStart := scrollBarBoxEnd
	lowerBoxEnd := float64(c.max)

	//x, y, width, height := c.layout.GetInnerRect()
	_, _, width, height := c.layout.GetInnerRect()
	var total = height - 2
	if c.orientation == ScrollBarHorizontal {
		total = width - 2
	}
	scale := float64(total) / (lowerBoxEnd - upperBoxStart)

	upperBoxStart = upperBoxStart
	upperBoxEnd = upperBoxEnd * scale
	scrollBarBoxStart = scrollBarBoxStart * scale
	scrollBarBoxEnd = scrollBarBoxEnd * scale
	lowerBoxStart = lowerBoxStart * scale
	lowerBoxEnd = lowerBoxEnd * scale

	// update scrollbar and "padding" boxes
	upperBoxSize := int(math.Max(0, upperBoxEnd-upperBoxStart))
	c.layout.ResizeItem(c.upperBox, upperBoxSize, 0)
	scrollBarBoxSize := int(math.Max(0, scrollBarBoxEnd-scrollBarBoxStart))
	c.layout.ResizeItem(c.scrollBarBox, scrollBarBoxSize, 0)
	lowerBoxSize := int(math.Max(0, lowerBoxEnd-lowerBoxStart))
	c.layout.ResizeItem(c.lowerBox, lowerBoxSize, 0)
}

func (c *ScrollbarComponent) SetWidth(width int) {
	c.barWidth = width
	c.updateLayout()
}
