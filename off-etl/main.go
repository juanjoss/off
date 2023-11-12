package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/juanjoss/off-etl/jobs"
	"github.com/juanjoss/off-etl/repository"
)

func main() {
	repo := repository.NewRepository()

	jobs.RunBrandsETL(repo)

	s := gocron.NewScheduler(time.UTC)

	maxRuns, err := strconv.Atoi(os.Getenv("ETL_MAX_RUNS"))
	if err != nil {
		log.Fatal("unable to set the limit of executions, check the MAX_RUNS env variable.")
	}
	s.LimitRunsTo(maxRuns)

	period, err := strconv.Atoi(os.Getenv("ETL_RUN_PERIOD"))
	if err != nil {
		log.Fatal("unable to set the execution period, check the ETL_RUN_PERIOD env variable.")
	}
	s.Every(period).Millisecond().Do(jobs.RunProductsETL, repo)

	s.StartBlocking()
}
