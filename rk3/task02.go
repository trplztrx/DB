package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Satellite struct {
	ID              int       `gorm:"column:ID Спутника;primaryKey"`
	Name            string    `gorm:"column:Название спутника"`
	ManufactureDate time.Time `gorm:"column:Дата производства"`
	Country         string    `gorm:"column:Страна"`
}

func (Satellite) TableName() string {
	return "rk3.satellite"
}

type Flight struct {
	ID          int       `gorm:"column:id;primaryKey"`
	SatelliteID int       `gorm:"column:ID Спутника"`
	LaunchDate  time.Time `gorm:"column:Дата запуска"`
	LaunchTime  time.Time `gorm:"column:Время запуска"`
	Weekday     string    `gorm:"column:День недели"`
	Type        int       `gorm:"column:Тип"`
}

func (Flight) TableName() string {
	return "rk3.flight"
}

func main() {
	// Подключение через pgxpool
	connString := "postgres://postgres:postgres@localhost:5432/bmstu_db"
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Printf("Ошибка подключения к БД: %v\n", err)
		return
	}
	defer pool.Close()

	// Подключение через GORM
	dsn := "host=localhost user=postgres password=postgres dbname=bmstu_db port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД через GORM: %v", err)
	}

	// Запрос 1 через pgxpool
	query1 := `
		SELECT DISTINCT "Страна"
		FROM rk3.satellite
		WHERE "Дата производства" = (
			SELECT MAX("Дата производства") FROM rk3.satellite
		);
	`
	rows, err := pool.Query(context.Background(), query1)
	if err != nil {
		fmt.Printf("Ошибка query через pgxpool: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("--- Запрос 1 через pgxpool ---")
	for rows.Next() {
		var country string
		if err := rows.Scan(&country); err != nil {
			fmt.Printf("Ошибка row через pgxpool: %v\n", err)
			return
		}
		fmt.Println(country)
	}

	// Запрос 1 через GORM
	var countries []string
	subQueryMaxDate := db.Model(&Satellite{}).
		Select("max(\"Дата производства\")")

	err = db.Model(&Satellite{}).
		Where("\"Дата производства\" = (?)", subQueryMaxDate).
		Distinct("\"Страна\"").
		Pluck("\"Страна\"", &countries).Error

	if err != nil {
		log.Println("Ошибка query через GORM", err)
	} else {
		fmt.Println("--- Запрос 1 через GORM ---")
		fmt.Println(countries)
	}

	// Запрос 2 через pgxpool
	query2 := `
		SELECT f."ID Спутника", s."Название спутника"
		FROM rk3.flight f
		JOIN rk3.satellite s ON f."ID Спутника" = s."ID Спутника"
		WHERE EXTRACT(YEAR FROM f."Дата запуска") = EXTRACT(YEAR FROM CURRENT_DATE)
		GROUP BY f."ID Спутника", s."Название спутника"
		HAVING COUNT(*) > 2;
	`
	rows, err = pool.Query(context.Background(), query2)
	if err != nil {
		fmt.Printf("Ошибка query через pgxpool: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("--- Запрос 2 через pgxpool ---")
	for rows.Next() {
		var satelliteID int
		var satelliteName string
		if err := rows.Scan(&satelliteID, &satelliteName); err != nil {
			fmt.Printf("Ошибка row через pgxpool: %v\n", err)
			return
		}
		fmt.Printf("ID: %d, Name: %s\n", satelliteID, satelliteName)
	}

	// Запрос 2 через GORM
	var results []struct {
		SatelliteID int
		Name        string
	}

	err = db.Table("rk3.flight f").
		Select("f.\"ID Спутника\" as satellite_id, s.\"Название спутника\" as name").
		Joins("join rk3.satellite s on f.\"ID Спутника\" = s.\"ID Спутника\"").
		Where("extract(year from f.\"Дата запуска\") = extract(year from current_date)").
		Group("f.\"ID Спутника\", s.\"Название спутника\"").
		Having("count(*) > 2").
		Scan(&results).Error

	if err != nil {
		log.Println("Ошибка query через GORM", err)
	} else {
		fmt.Println("--- Запрос 2 через GORM ---")
		fmt.Println(results)
	}

	// Запрос 3 через pgxpool
	query3 := `
		SELECT DISTINCT s."Страна"
		FROM rk3.satellite s
		WHERE s."ID Спутника" IN (
			SELECT f."ID Спутника"
			FROM rk3.flight f
			GROUP BY f."ID Спутника"
			HAVING MIN(f."Дата запуска") > DATE '2024-10-01'
		);
	`
	rows, err = pool.Query(context.Background(), query3)
	if err != nil {
		fmt.Printf("Ошибка query через pgxpool: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("--- Запрос 3 через pgxpool ---")
	for rows.Next() {
		var country string
		if err := rows.Scan(&country); err != nil {
			fmt.Printf("Ошибка row через pgxpool: %v\n", err)
			return
		}
		fmt.Println(country)
	}

	// Запрос 3 через GORM
	var countriesQ3 []string
	subQueryFirstLaunch := db.Table("rk3.flight f").
		Select("f.\"ID Спутника\"").
		Group("f.\"ID Спутника\"").
		Having("min(f.\"Дата запуска\") > ?", "2024-10-01")

	err = db.Table("rk3.satellite s").
		Distinct("s.\"Страна\"").
		Where("s.\"ID Спутника\" in (?)", subQueryFirstLaunch).
		Pluck("s.\"Страна\"", &countriesQ3).Error

	if err != nil {
		log.Println("Ошибка query через GORM", err)
	} else {
		fmt.Println("--- Запрос 3 через GORM ---")
		fmt.Println(countriesQ3)
	}
}
