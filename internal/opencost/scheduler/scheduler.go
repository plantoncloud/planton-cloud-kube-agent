package scheduler

import (
	"context"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/config"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/opencost/exporter"
	log "github.com/sirupsen/logrus"
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
