// Copyright 2019 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"time"

	apiMetrics "github.com/cilium/cilium/pkg/api/helpers/metrics"
	ec2shim "github.com/cilium/cilium/pkg/aws/ec2"
	"github.com/cilium/cilium/pkg/aws/eni"
	"github.com/cilium/cilium/pkg/controller"
	"github.com/cilium/cilium/pkg/ipam"
	ipamMetrics "github.com/cilium/cilium/pkg/ipam/metrics"
	"github.com/cilium/cilium/pkg/option"
	"github.com/cilium/cilium/pkg/trigger"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/sirupsen/logrus"
)

// startENIAllocator kicks of ENI allocation, the initial connection to AWS
// APIs is done in a blocking manner, given that is successful, a controller is
// started to manage allocation based on CiliumNode custom resources
func startENIAllocator(awsClientQPSLimit float64, awsClientBurst int, eniTags map[string]string) error {
	log.Info("Starting ENI allocator...")

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return fmt.Errorf("unable to load AWS configuration: %s", err)
	}

	log.Info("Retrieving own metadata from EC2 metadata server...")
	metadataClient := ec2metadata.New(cfg)
	instance, err := metadataClient.GetInstanceIdentityDocument()
	if err != nil {
		return fmt.Errorf("unable to retrieve instance identity document: %s", err)
	}

	log.WithFields(logrus.Fields{
		"instance": instance.InstanceID,
		"region":   instance.Region,
	}).Info("Connected to EC2 metadata server")

	cfg.Region = instance.Region

	var (
		ec2Client *ec2shim.Client
		instances *eni.InstancesManager
	)

	if enableMetrics {
		eniMetrics := apiMetrics.NewLegacyPrometheusMetrics(metricNamespace, "eni", "aws_api_duration_seconds", "ec2_rate_limit_duration_seconds", registry)
		ec2Client = ec2shim.NewClient(ec2.New(cfg), eniMetrics, awsClientQPSLimit, awsClientBurst)
		log.Info("Connected to EC2 service API")
		instances = eni.NewInstancesManager(ec2Client, eniTags)
		iMetrics := ipamMetrics.NewPrometheusMetrics(metricNamespace, registry)
		nodeManager, err = ipam.NewNodeManager(instances, &k8sAPI{}, iMetrics, eniParallelWorkers,
			option.Config.AwsReleaseExcessIps)
		if err != nil {
			return fmt.Errorf("unable to initialize ENI node manager: %s", err)
		}
	} else {
		// Inject dummy metrics operations that do nothing so we don't panic if
		// metrics aren't enabled
		noOpMetric := &noOpMetrics{}
		ec2Client = ec2shim.NewClient(ec2.New(cfg), &apiMetrics.NoOpMetrics{}, awsClientQPSLimit, awsClientBurst)
		log.Info("Connected to EC2 service API")
		instances = eni.NewInstancesManager(ec2Client, eniTags)
		nodeManager, err = ipam.NewNodeManager(instances, &k8sAPI{}, noOpMetric, eniParallelWorkers,
			option.Config.AwsReleaseExcessIps)
		if err != nil {
			return fmt.Errorf("unable to initialize ENI node manager: %s", err)
		}
	}

	// Initial blocking synchronization of all ENIs and subnets
	instances.Resync(context.TODO())

	// Start an interval based  background resync for safety, it will
	// synchronize the state regularly and resolve eventual deficit if the
	// event driven trigger fails, and also release excess IP addresses
	// if release-excess-ips is enabled
	go func() {
		time.Sleep(time.Minute)
		mngr := controller.NewManager()
		mngr.UpdateController("eni-refresh",
			controller.ControllerParams{
				RunInterval: time.Minute,
				DoFunc: func(ctx context.Context) error {
					syncTime := instances.Resync(ctx)
					nodeManager.Resync(ctx, syncTime)
					return nil
				},
			})
	}()

	return nil
}

// The below are types which fulfill various interfaces which are needed by the
// eni / ec2 functions we have which do nothing if metrics are disabled.

type noOpMetricsObserver struct{}

// MetricsObserver implementation
func (m *noOpMetricsObserver) PostRun(callDuration, latency time.Duration, folds int) {}
func (m *noOpMetricsObserver) QueueEvent(reason string)                               {}

type noOpMetrics struct{}

// eni metricsAPI interface implementation
func (m *noOpMetrics) IncAllocationAttempt(status, subnetID string)                              {}
func (m *noOpMetrics) AddIPAllocation(subnetID string, allocated int64)                          {}
func (m *noOpMetrics) AddIPRelease(subnetID string, released int64)                              {}
func (m *noOpMetrics) SetAllocatedIPs(typ string, allocated int)                                 {}
func (m *noOpMetrics) SetAvailableENIs(available int)                                            {}
func (m *noOpMetrics) SetAvailableIPsPerSubnet(subnetID, availabilityZone string, available int) {}
func (m *noOpMetrics) SetNodes(category string, nodes int)                                       {}
func (m *noOpMetrics) IncResyncCount()                                                           {}
func (m *noOpMetrics) PoolMaintainerTrigger() trigger.MetricsObserver {
	return &noOpMetricsObserver{}
}
func (m *noOpMetrics) K8sSyncTrigger() trigger.MetricsObserver {
	return &noOpMetricsObserver{}
}
func (m *noOpMetrics) ResyncTrigger() trigger.MetricsObserver {
	return &noOpMetricsObserver{}
}

// ec2 metricsAPI interface implementation
func (m *noOpMetrics) ObserveEC2APICall(call, status string, duration float64)      {}
func (m *noOpMetrics) ObserveEC2RateLimit(operation string, duration time.Duration) {}
