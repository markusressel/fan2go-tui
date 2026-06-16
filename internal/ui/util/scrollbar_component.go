package util

import (
	"fan2go-tui/internal/ui/theme"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ScrollBarOrientation int

type ScrollbarRuneType int

const (
	ScrollbarRuneTypeLeft ScrollbarRuneType = iota
	ScrollbarRuneTypeRight
	ScrollbarRuneTypeTop
	ScrollbarRuneTypeBottom
)

const (
	ScrollBarVertical ScrollBarOrientation = iota
	ScrollBarHorizontal

	ScrollIndicatorMiddle = "■"
	ScrollIndicatorTop    = "▲"
	ScrollIndicatorBottom = "▼"
	ScrollIndicatorLeft   = "◀"
	ScrollIndicatorRight  = "▶"
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

	lastKnownWidth  int
	lastKnownHeight int
}

// NewScrollbarComponent creates a new ScrollbarComponent.
// The application is used to redraw the component.
// The orientation is used to set the orientation of the scrollbar.
// The min is the minimum value of the scrollbar.
// The max is the maximum value of the scrollbar.
// The scrollPosition is the current position of the scrollbar.
// The barWidth is the width of the scrollbar.
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
	layout.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		if width != c.lastKnownWidth || height != c.lastKnownHeight {
			c.lastKnownWidth = width
			c.lastKnownHeight = height
			c.updateLayoutInternal()
		}
		return layout.GetInnerRect()
	})

	c.topArrow = tview.NewTextView()
	layout.AddItem(c.topArrow, 1, 0, false)
	c.upperBox = tview.NewBox()
	c.upperBox.SetBackgroundColor(theme.Colors.List.Scrollbar.Background)
	layout.AddItem(c.upperBox, 1, 0, false)
	c.scrollBarBox = tview.NewBox()
	c.scrollBarBox.SetBackgroundColor(theme.Colors.List.Scrollbar.Bar)
	layout.AddItem(c.scrollBarBox, 1, 0, false)
	c.lowerBox = tview.NewBox()
	c.lowerBox.SetBackgroundColor(theme.Colors.List.Scrollbar.Background)
	layout.AddItem(c.lowerBox, 1, 0, false)
	c.bottomArrow = tview.NewTextView()
	layout.AddItem(c.bottomArrow, 1, 0, false)

	c.layout = layout
}

func (c *ScrollbarComponent) UpdateLayout() {
	c.updateLayoutInternal()

	c.application.ForceDraw()
}

func (c *ScrollbarComponent) updateLayoutInternal() {
	c.updateTopEndText()
	c.updateScrollbar()
	c.updateBottomEndText()
}

func (c *ScrollbarComponent) SetOrientation(orientation ScrollBarOrientation) {
	c.orientation = orientation
	switch c.orientation {
	case ScrollBarVertical:
		c.layout.SetDirection(tview.FlexRow)
	case ScrollBarHorizontal:
		c.layout.SetDirection(tview.FlexColumn)
	}
	c.UpdateLayout()
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
	c.UpdateLayout()
}

func (c *ScrollbarComponent) SetMinInternal(min int) {
	c.min = min
	c.updateLayoutInternal()
}

func (c *ScrollbarComponent) SetMax(max int) {
	c.max = max
	c.UpdateLayout()
}

func (c *ScrollbarComponent) SetMaxInternal(max int) {
	c.max = max
	c.updateLayoutInternal()
}

func (c *ScrollbarComponent) SetPosition(position int) {
	if position < 0 {
		position = 0
	}
	c.scrollPosition = position
	c.UpdateLayout()
}

func (c *ScrollbarComponent) SetPositionInternal(position int) {
	if position < 0 {
		position = 0
	}
	c.scrollPosition = position
	c.updateLayoutInternal()
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

	c.UpdateLayout()
}

func (c *ScrollbarComponent) ScrollToTop() {
	c.scroll(-c.GetPosition())
}

func (c *ScrollbarComponent) calculateBarWidth() int {
	// calculate the bar width
	if c.max <= c.min {
		return 1
	}
	barWidth := int(math.Max(1, float64(c.max/(c.max-c.min))))
	return barWidth
}

func (c *ScrollbarComponent) updateTopEndText() {
	c.layout.ResizeItem(c.topArrow, 1, 0)
	isAtLimit := c.scrollPosition <= c.min
	text, textColor := c.determineRuneAndColor(ScrollbarRuneTypeTop, isAtLimit)
	c.topArrow.SetText(text)
	c.topArrow.SetTextColor(textColor)
}

func (c *ScrollbarComponent) updateBottomEndText() {
	c.layout.ResizeItem(c.bottomArrow, 1, 0)
	isAtLimit := c.scrollPosition+c.barWidth >= c.max
	text, textColor := c.determineRuneAndColor(ScrollbarRuneTypeBottom, isAtLimit)
	c.bottomArrow.SetText(text)
	c.bottomArrow.SetTextColor(textColor)
}

func (c *ScrollbarComponent) determineRuneAndColor(
	scrollbarRuneType ScrollbarRuneType,
	isAtLimit bool,
) (text string, textColor tcell.Color) {
	switch c.orientation {
	case ScrollBarVertical:
		if isAtLimit {
			text = ScrollIndicatorMiddle
			textColor = theme.Colors.List.Scrollbar.IndicatorInactive
		} else {
			switch scrollbarRuneType {
			case ScrollbarRuneTypeBottom:
				text = ScrollIndicatorBottom
			case ScrollbarRuneTypeTop:
				fallthrough
			default:
				text = ScrollIndicatorTop
			}
			textColor = theme.Colors.List.Scrollbar.IndicatorActive
		}
	case ScrollBarHorizontal:
		if isAtLimit {
			text = ScrollIndicatorMiddle
			textColor = theme.Colors.List.Scrollbar.IndicatorInactive
		} else {
			switch scrollbarRuneType {
			case ScrollbarRuneTypeRight:
				text = ScrollIndicatorRight
			case ScrollbarRuneTypeLeft:
				fallthrough
			default:
				text = ScrollIndicatorLeft
			}
			textColor = theme.Colors.List.Scrollbar.IndicatorActive
		}
	}
	return text, textColor
}

func (c *ScrollbarComponent) updateScrollbar() {
	// calculate the box sizes
	upperBoxStart := float64(c.min)
	upperBoxEnd := float64(c.min + c.scrollPosition)
	scrollBarBoxEnd := upperBoxEnd + float64(c.barWidth)
	lowerBoxEnd := float64(c.max)

	//x, y, width, height := c.layout.GetInnerRect()
	_, _, width, height := c.layout.GetInnerRect()
	var total = height - 2
	if c.orientation == ScrollBarHorizontal {
		total = width - 2
	}

	if total <= 0 {
		return
	}

	scale := float64(total) / math.Max(1, lowerBoxEnd-upperBoxStart)

	upperBoxEndScaled := (upperBoxEnd - upperBoxStart) * scale
	scrollBarBoxEndScaled := (scrollBarBoxEnd - upperBoxStart) * scale

	upperBoxSize := int(math.Round(upperBoxEndScaled))
	scrollBarBoxEndInt := int(math.Round(scrollBarBoxEndScaled))

	scrollBarBoxSize := scrollBarBoxEndInt - upperBoxSize
	if scrollBarBoxSize < 1 {
		scrollBarBoxSize = 1
		if upperBoxSize+scrollBarBoxSize > total {
			upperBoxSize = total - scrollBarBoxSize
			if upperBoxSize < 0 {
				upperBoxSize = 0
			}
		}
	}

	lowerBoxSize := total - (upperBoxSize + scrollBarBoxSize)
	if lowerBoxSize < 0 {
		lowerBoxSize = 0
	}

	// update scrollbar and "padding" boxes
	c.layout.ResizeItem(c.upperBox, upperBoxSize, 0)
	c.layout.ResizeItem(c.scrollBarBox, scrollBarBoxSize, 0)
	c.layout.ResizeItem(c.lowerBox, lowerBoxSize, 0)
}

func (c *ScrollbarComponent) SetWidth(width int) {
	c.barWidth = width
	c.UpdateLayout()
}

func (c *ScrollbarComponent) SetWidthInternal(width int) {
	c.barWidth = width
	c.updateLayoutInternal()
}
