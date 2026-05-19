package graph

import (
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

// OverlayPlot extends tvxwidgets.Plot with lightweight overlay rendering.
type OverlayPlot[T any] struct {
	*tvxwidgets.Plot
	overlays   []GraphComponentOverlay[T]
	overlayCtx OverlayRenderContext[T]
}

func NewOverlayPlot[T any]() *OverlayPlot[T] {
	return &OverlayPlot[T]{
		Plot: tvxwidgets.NewPlot(),
	}
}

func (p *OverlayPlot[T]) SetOverlays(overlays []GraphComponentOverlay[T]) {
	p.overlays = append([]GraphComponentOverlay[T]{}, overlays...)
}

func (p *OverlayPlot[T]) SetOverlayContext(ctx OverlayRenderContext[T]) {
	p.overlayCtx = ctx
}

func (p *OverlayPlot[T]) Draw(screen tcell.Screen) {
	p.Plot.Draw(screen)
	ctx := p.overlayCtx
	ctx.Plot = p.Plot
	ctx.Background = p.GetBackgroundColor()
	for _, overlay := range p.overlays {
		overlay.draw(screen, ctx)
	}
}
