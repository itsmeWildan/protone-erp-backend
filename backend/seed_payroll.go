package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr := "postgres://postgres:postgres@localhost:5432/protone_erp?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	tenantID := uuid.MustParse("333ec068-21cf-4ec1-926f-db9a82dcecc6")

	components := []struct {
		Code   string
		Name   string
		Type   string
		Amount float64
	}{
		{"BASIC", "Gaji Pokok", "allowance", 5000000},
		{"TRANSPORT", "Tunjangan Transport", "allowance", 500000},
		{"MEAL", "Tunjangan Makan", "allowance", 300000},
		{"BPJS_TK", "BPJS Ketenagakerjaan", "deduction", 100000},
		{"BPJS_KES", "BPJS Kesehatan", "deduction", 50000},
	}

	for _, c := range components {
		query := `INSERT INTO salary_components (id, tenant_id, code, name, type, default_amount, is_taxable, is_fixed) 
		          VALUES ($1, $2, $3, $4, $5, $6, true, true)
		          ON CONFLICT (tenant_id, code) DO NOTHING`
		_, err := pool.Exec(context.Background(), query, uuid.New(), tenantID, c.Code, c.Name, c.Type, c.Amount)
		if err != nil {
			log.Printf("Failed to seed component %s: %v\n", c.Code, err)
		} else {
			fmt.Printf("Seeded salary component: %s\n", c.Name)
		}
	}

	fmt.Println("Salary components seeding completed!")
}
