package scheduler

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gitlab.com/plantoncode/planton/pcs/lib/mod/planton-cloud-kube-agent.git/internal/config"
	"gitlab.com/plantoncode/planton/pcs/lib/mod/planton-cloud-kube-agent.git/internal/opencost/exporter"
	"time"
)

// Start the periodic scheduler to scrape data from open-cost and export it to planton-cloud
func Start(ctx context.Context, c *config.Config) {
	if err := exporter.Export(ctx, c); err != nil {
		log.Fatalf("failed to export cost-allocation data to planton-cloud  with err: %v", err)
	}
	ticker := time.NewTicker(time.Duration(c.OpenCostPollingIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := exporter.Export(ctx, c); err != nil {
				log.Fatalf("failed to export cost-allocation data to planton-cloud  with err: %v", err)
			}
		}
	}
}
