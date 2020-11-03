package plans

import (
	"database/sql"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Plan struct {
	Id    int
	Name  string
	Price float32
}

func GetPlans(db *sql.DB) []Plan {
	query := `SELECT id, name, price FROM plan`

	rows, err := db.Query(query)
	HandleErr(err)
	defer rows.Close()

	var plans []Plan

	for rows.Next() {
		plan := Plan{}
		err = rows.Scan(&plan.Id, &plan.Name, &plan.Price)
		HandleErr(err)

		plans = append(plans, plan)
	}

	return plans
}

func GetPlanById(db *sql.DB, id int) (Plan, error) {
	query := `SELECT id, name, price FROM plan WHERE id  = ?`
	plan := Plan{}
	err := db.QueryRow(query, id).Scan(&plan.Id, &plan.Name, &plan.Price)

	return plan, err
}
