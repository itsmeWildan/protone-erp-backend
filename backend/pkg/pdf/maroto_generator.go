package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type PayrollSlipData struct {
	CompanyName    string
	Period         string
	EmployeeName   string
	EmployeeNIK    string
	Department     string
	Position       string
	BasicSalary    float64
	Allowances     []ComponentData
	Deductions     []ComponentData
	OvertimeAmount float64
	NetSalary      float64
}

type ComponentData struct {
	Name   string
	Amount float64
}

type marotoGenerator struct{}

func NewMarotoGenerator() Generator {
	return &marotoGenerator{}
}

func (m *marotoGenerator) GeneratePayrollSlip(ctx context.Context, data interface{}) (io.Reader, error) {
	d := data.(PayrollSlipData)

	cfg := config.NewBuilder().
		WithPageNumber().
		Build()

	mrt := maroto.New(cfg)

	// Header
	mrt.AddRows(
		row.New(20).Add(
			col.New(8).Add(
				text.New(d.CompanyName, props.Text{
					Size:  16,
					Style: fontstyle.Bold,
				}),
				text.New("SLIP GAJI KARYAWAN", props.Text{
					Top:  7,
					Size: 12,
				}),
			),
			col.New(4).Add(
				text.New(fmt.Sprintf("Periode: %s", d.Period), props.Text{
					Align: align.Right,
					Size:  10,
				}),
			),
		),
	)

	mrt.AddRows(row.New(5)) // Spacer

	// Employee Info
	mrt.AddRows(
		row.New(15).Add(
			col.New(6).Add(
				text.New(fmt.Sprintf("Nama       : %s", d.EmployeeName), props.Text{Size: 10}),
				text.New(fmt.Sprintf("NIK        : %s", d.EmployeeNIK), props.Text{Size: 10, Top: 5}),
			),
			col.New(6).Add(
				text.New(fmt.Sprintf("Departemen : %s", d.Department), props.Text{Size: 10}),
				text.New(fmt.Sprintf("Jabatan    : %s", d.Position), props.Text{Size: 10, Top: 5}),
			),
		),
	)

	mrt.AddRows(row.New(10)) // Spacer

	// Income Header
	mrt.AddRows(
		row.New(10).Add(
			col.New(12).Add(
				text.New("PENGHASILAN", props.Text{Style: fontstyle.Bold, Size: 10}),
			),
		),
		row.New(2).Add(col.New(12).Add(line.New())),
	)

	// Basic Salary
	mrt.AddRows(
		row.New(8).Add(
			col.New(8).Add(text.New("- Gaji Pokok", props.Text{Size: 10})),
			col.New(4).Add(text.New(formatCurrency(d.BasicSalary), props.Text{Size: 10, Align: align.Right})),
		),
	)

	// Allowances
	for _, a := range d.Allowances {
		mrt.AddRows(
			row.New(8).Add(
				col.New(8).Add(text.New(fmt.Sprintf("- %s", a.Name), props.Text{Size: 10})),
				col.New(4).Add(text.New(formatCurrency(a.Amount), props.Text{Size: 10, Align: align.Right})),
			),
		)
	}

	// Overtime
	if d.OvertimeAmount > 0 {
		mrt.AddRows(
			row.New(8).Add(
				col.New(8).Add(text.New("- Uang Lembur", props.Text{Size: 10})),
				col.New(4).Add(text.New(formatCurrency(d.OvertimeAmount), props.Text{Size: 10, Align: align.Right})),
			),
		)
	}

	mrt.AddRows(row.New(5))

	// Deduction Header
	mrt.AddRows(
		row.New(10).Add(
			col.New(12).Add(
				text.New("POTONGAN", props.Text{Style: fontstyle.Bold, Size: 10}),
			),
		),
		row.New(2).Add(col.New(12).Add(line.New())),
	)

	// Deductions
	for _, dec := range d.Deductions {
		mrt.AddRows(
			row.New(8).Add(
				col.New(8).Add(text.New(fmt.Sprintf("- %s", dec.Name), props.Text{Size: 10})),
				col.New(4).Add(text.New(formatCurrency(dec.Amount), props.Text{Size: 10, Align: align.Right})),
			),
		)
	}

	mrt.AddRows(row.New(10).Add(col.New(12).Add(line.New()))) // Separator before total

	// Total / Take Home Pay
	mrt.AddRows(
		row.New(15).Add(
			col.New(8).Add(
				text.New("TAKE HOME PAY (TOTAL BERSIH)", props.Text{Style: fontstyle.Bold, Size: 11}),
			),
			col.New(4).Add(
				text.New(formatCurrency(d.NetSalary), props.Text{
					Style: fontstyle.Bold,
					Size:  11,
					Align: align.Right,
				}),
			),
		),
	)

	mrt.AddRows(row.New(20)) // Spacer

	// Footer
	mrt.AddRows(
		row.New(20).Add(
			col.New(8).Add(
				text.New("Catatan:", props.Text{Size: 8, Style: fontstyle.Italic}),
				text.New("Slip gaji ini digenerate secara otomatis oleh sistem ProtoERP.", props.Text{Size: 8, Top: 4}),
			),
			col.New(4).Add(
				text.New("Hormat Kami,", props.Text{Size: 10, Align: align.Center}),
				text.New("HR Manager", props.Text{Size: 10, Top: 15, Align: align.Center, Style: fontstyle.Bold}),
			),
		),
	)

	doc, err := mrt.Generate()
	if err != nil {
		return nil, err
	}

	// Maroto v2 uses GetBytes() instead of GetReader()
	return bytes.NewReader(doc.GetBytes()), nil
}

func formatCurrency(amount float64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("Rp %d", int64(amount))
}
