package pdf

import (
	"context"
	"io"
)

// Generator adalah interface untuk membuat file PDF.
// Dengan interface ini, kita bisa ganti dari Maroto ke Chromedp dengan mudah nanti.
type Generator interface {
	GeneratePayrollSlip(ctx context.Context, data interface{}) (io.Reader, error)
}
