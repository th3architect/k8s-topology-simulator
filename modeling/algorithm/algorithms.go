/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package algorithm

import "k8s.io/klog/v2"

// NewAlgorithm serves as an algorithm constructor based on the algroithm name
func NewAlgorithm(name string) RoutingAlgorithm {
	switch name {
	case "SharedGlobal", "SharedGlobalAlgorithm":
		klog.Info("SharedGlobalAlgorithm created")
		return SharedGlobalAlgorithm{sharedCoreAlgorithm: sharedGlobalAlgorithmCore{globalWeight: 0.4, globalThreshold: 100}}
	case "SharedMultiZone", "SharedMultiZoneAlgorithm":
		klog.Info("SharedMultiZoneAlgorithm created")
		return SharedMultiZoneAlgorithm{sharedCoreAlgorithm: sharedGlobalAlgorithmCore{globalWeight: 1, globalThreshold: 100}}
	case "Local", "LocalAlgorithm":
		klog.Info("LocalSliceAlgorithm created")
		return LocalSliceAlgorithm{threshold: 0.5, startingThreshold: 3}
	case "LocalWeighted", "LocalWeightedAlgorithm":
		klog.Info("LocalWeightedSliceAlgorithm created")
		return LocalWeightedSliceAlgorithm{}
	case "LocalOpt", "LocalOptAlgorithm":
		klog.Info("LocalSliceAlgorithmOpt created")
		return LocalSliceAlgorithmOpt{}
	case "LocalShared", "LocalSharedAlgorithm":
		klog.Info("LocalSharedSliceAlgorithm created")
		return LocalSharedSliceAlgorithm{threshold: 0.5}
	case "Original", "OriginalAlgorithm":
		klog.Info("OriginalAlgorithm created")
		return OriginalAlgorithm{}
	}
	klog.Warningf("[WARNINIG] unknown algorithm %v, return LocalSliceAlgorithm as default\n", name)
	return LocalSliceAlgorithm{}
}
