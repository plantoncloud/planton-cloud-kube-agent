package exporter

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/opencost/opencost/pkg/kubecost"
	"github.com/pkg/errors"
	hostingcostrpc "github.com/plantoncloud-inc/company-protos/zzgo/planton/company/proto/v1/hosting/cost/rpc"
	commonsresource "github.com/plantoncloud-inc/go-commons/domain/common/resource"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/apiclient"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/auth/token"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/config"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/opencost/client"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/opencost/labels"
	"github.com/plantoncloud-inc/proto-commons/zzgo/planton/commons/proto/v1/resource/enums"
	log "github.com/sirupsen/logrus"
	"time"
)

func Export(ctx context.Context, c *config.Config) error {
	openCostClient := client.Client{BaseUrl: c.OpenCostApiEndpoint}
	openCostAllocations, err := openCostClient.GetAllocationAggregatedByPodByWindow(getHourlyWindow())
	if err != nil {
		return errors.Wrap(err, "failed to get cost-allocations from open-cost")
	}
	conn, err := apiclient.NewConn(c.PlantonCloudServiceApiEndpoint, token.Token)
	if err != nil {
		return errors.Wrap(err, "failed to get new planton-cloud api client")
	}
	costAllocationCmdControllerClient := hostingcostrpc.NewCostAllocationCmdControllerClient(conn)
	plantonCostAllocations := make([]*hostingcostrpc.CostAllocation, 0)
	for _, ca := range openCostAllocations {
		workloadLabels := labels.GetLabels(ca)

		if commonsresource.ResourceTypeStringToEnum(workloadLabels.ResourceType) == enums.ResourceType_RESOURCE_TYPE_UNSPECIFIED {
			//skip export for unknown workloads
			log.Debugf("skipping export for %s as the resource-type is unknown", ca.Name)
			continue
		}

		var plantonCloudCostAllocation hostingcostrpc.CostAllocation

		if err := copier.Copy(&plantonCloudCostAllocation, &ca); err != nil {
			return errors.Wrap(err, "failed to copy cost-allocation object from open-cost to planton-cloud")
		}

		copy(&plantonCloudCostAllocation, ca)

		plantonCloudCostAllocation.CompanyId = workloadLabels.Company
		plantonCloudCostAllocation.ProductId = workloadLabels.Product
		plantonCloudCostAllocation.HostingClusterId = c.PlantonCloudKubeAgentHostingClusterId
		plantonCloudCostAllocation.EnvironmentId = workloadLabels.Environment
		plantonCloudCostAllocation.ResourceType = workloadLabels.ResourceType
		plantonCloudCostAllocation.ResourceId = workloadLabels.ResourceId
		plantonCloudCostAllocation.StartTs = ca.Start
		plantonCloudCostAllocation.EndTs = ca.End
		plantonCostAllocations = append(plantonCostAllocations, &plantonCloudCostAllocation)
	}
	log.Debugf("exporting %d entriesto planton-cloud", len(plantonCostAllocations))
	if _, err := costAllocationCmdControllerClient.Create(ctx,
		&hostingcostrpc.CostAllocations{
			Entries: plantonCostAllocations,
		}); err != nil {
		return errors.Wrap(err, "failed to export cost-allocation to planton-cloud")
	}
	log.Infof("successfully exported %d cost-allocations to planton-cloud", len(plantonCostAllocations))
	return nil
}

func getHourlyWindow() string {
	now := time.Now()
	endWindow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	startWindow := endWindow.Add(-time.Hour * 1)
	return fmt.Sprintf("%s,%s", startWindow.Format("2006-01-02T15:04:05Z"), endWindow.Format("2006-01-02T15:04:05Z"))
}

func copy(costAllocation *hostingcostrpc.CostAllocation, allocationJson *kubecost.AllocationJSON) {
	costAllocation.CpuCores = *allocationJson.CPUCores
	costAllocation.CpuCoreRequestAverage = *allocationJson.CPUCoreRequestAverage
	costAllocation.CpuCost = *allocationJson.CPUCost
	costAllocation.CpuCostAdjustment = *allocationJson.CPUCostAdjustment
	costAllocation.CpuCoreHours = *allocationJson.CPUCoreHours
	costAllocation.CpuEfficiency = *allocationJson.CPUEfficiency
	costAllocation.CpuCoreUsageAverage = *allocationJson.CPUCoreUsageAverage

	costAllocation.GpuCost = *allocationJson.GPUCost
	costAllocation.GpuCostAdjustment = *allocationJson.GPUCostAdjustment
	costAllocation.GpuCount = *allocationJson.GPUCount
	costAllocation.GpuHours = *allocationJson.GPUHours

	costAllocation.PvCost = *allocationJson.PVCost
	costAllocation.PvCostAdjustment = *allocationJson.PVCostAdjustment
	costAllocation.PvBytes = *allocationJson.PVBytes
	costAllocation.PvByteHours = *allocationJson.PVByteHours

	costAllocation.RamCost = *allocationJson.RAMCost
	costAllocation.RamCostAdjustment = *allocationJson.RAMCostAdjustment
	costAllocation.RamByteHours = *allocationJson.RAMByteHours
	costAllocation.RamBytes = *allocationJson.RAMBytes
	costAllocation.RamByteRequestAverage = *allocationJson.RAMByteRequestAverage
	costAllocation.RamByteUsageAverage = *allocationJson.RAMByteUsageAverage
	costAllocation.RamEfficiency = *allocationJson.RAMEfficiency
}
