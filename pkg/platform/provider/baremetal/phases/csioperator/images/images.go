/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package images

import (
	"fmt"
	"reflect"
	"sort"

	"tkestack.io/tke/pkg/util/containerregistry"
)

// CSIVersion indicates the version of CSI external components.
type CSIVersion string

const (
	// LatestVersion is latest version of addon.
	// TODO: bump up to v1.0.3
	LatestVersion = "v1.0.2"

	// CSIVersionV0 indicates the 0.3.0 version of CSI.
	CSIVersionV0 = "v0.0"
	// CSIVersionV1 indicates the 1.x version of CSI.
	CSIVersionV1 = "v1.0"
	// CSIVersionV1p1 indicates the 1.1+ version of CSI in tencent cloud cvm,
	// which does not need to use secret id and key.
	CSIVersionV1p1 = "v1.1"

	// CSIDriverCephRBD indicates the CephRBD storage type.
	CSIDriverCephRBD = "csi-rbd"
	// CSIDriverCephFS indicates the CephFS storage type.
	CSIDriverCephFS = "csi-cephfs"
	// CSIDriverTencentCBS indicates the Tencent Cloud CBS storage type.
	CSIDriverTencentCBS = "com.tencent.cloud.csi.cbs"
)

// csiVersion is the set of versions of all CSI components.
type csiVersion struct {
	Provisioner      string
	Attacher         string
	Resizer          string
	Snapshotter      string
	LivenessProbe    string
	NodeRegistrar    string
	ClusterRegistrar string
	Driver           string
}

// csiVersionMap stores all images of CSI need. Refer from
// <https://github.com/tkestack/csi-operator/blob/74188bd0f7462446109ee82f7488d8bd3646f525/pkg/controller/csi/enhancer/enhancer.go#L64>
// Need to keep same with the csi-operator version.
var csiVersionMap = map[string]map[CSIVersion]*csiVersion{
	CSIDriverCephRBD: {
		CSIVersionV0: {
			Provisioner:   "csi-provisioner:v0.4.2",
			Attacher:      "csi-attacher:v0.4.2",
			Snapshotter:   "csi-snapshotter:v0.4.1",
			LivenessProbe: "livenessprobe:v0.4.1",
			NodeRegistrar: "driver-registrar:v0.3.0",
			Driver:        "rbdplugin:v0.3.0",
		},
		CSIVersionV1: {
			Provisioner:   "csi-provisioner:v1.0.1",
			Attacher:      "csi-attacher:v1.1.0",
			Snapshotter:   "csi-snapshotter:v1.1.0",
			LivenessProbe: "livenessprobe:v1.1.0",
			NodeRegistrar: "csi-node-driver-registrar:v1.1.0",
			Driver:        "rbdplugin:v1.0.0",
			// TODO: Add resizer.
			// Resizer:          "v0.1.0",
		},
	},
	CSIDriverCephFS: {
		CSIVersionV0: {
			Provisioner:   "csi-provisioner:v0.4.2",
			Attacher:      "csi-attacher:v0.4.2",
			LivenessProbe: "livenessprobe:v0.4.1",
			NodeRegistrar: "driver-registrar:v0.3.0",
			Driver:        "cephfsplugin:v0.3.0",
		},
		CSIVersionV1: {
			Provisioner:   "csi-provisioner:v1.0.1",
			Attacher:      "csi-attacher:v1.1.0",
			LivenessProbe: "livenessprobe:v1.1.0",
			NodeRegistrar: "csi-node-driver-registrar:v1.1.0",
			Driver:        "cephfsplugin:v1.0.0",
			// TODO: Add resizer.
			// Resizer:          "v0.1.0",
		},
	},
	CSIDriverTencentCBS: {
		CSIVersionV0: {
			Provisioner:   "csi-provisioner:v0.4.2",
			Attacher:      "csi-attacher:v0.4.2",
			NodeRegistrar: "driver-registrar:v0.3.0",
			Driver:        "csi-tencentcloud-cbs:v0.2.1",
		},
		CSIVersionV1: {
			Provisioner:   "csi-provisioner:v1.2.0",
			Attacher:      "csi-attacher:v1.1.0",
			Snapshotter:   "csi-snapshotter:v1.2.2",
			NodeRegistrar: "csi-node-driver-registrar:v1.1.0",
			// TODO:NOTE--TKE Stack now use a old version csi-operator image (ID sha256:b77952b83730),
			// which only looks like v1.0.2. Version of driver in this image is v1.0.0.
			// TODO: FIX--After csi-operator bump up to v1.0.3, use the right version v1.2.0
			//Driver:        "csi-tencentcloud-cbs:v1.2.0",
			Driver:  "csi-tencentcloud-cbs:v1.0.0",
			Resizer: "csi-resizer:v0.5.0",
		},
		CSIVersionV1p1: {
			Provisioner:   "csi-provisioner:v1.2.0",
			Attacher:      "csi-attacher:v1.1.0",
			Snapshotter:   "csi-snapshotter:v1.2.2",
			NodeRegistrar: "csi-node-driver-registrar:v1.1.0",
			Driver:        "csi-tencentcloud-cbs:v1.2.0",
			Resizer:       "csi-resizer:v0.5.0",
		},
	},
}

type Components struct {
	CSIOperator containerregistry.Image
}

func (c Components) Get(name string) *containerregistry.Image {
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumField(); i++ {
		v, _ := v.Field(i).Interface().(containerregistry.Image)
		if v.Name == name {
			return &v
		}
	}
	return nil
}

var versionMap = map[string]Components{
	LatestVersion: {
		// TODO: bump up to v1.0.3
		CSIOperator: containerregistry.Image{Name: "csi-operator", Tag: "v1.0.2"},
	},
}

func List() []string {
	items := make([]string, 0, len(versionMap))
	versions := Versions()
	for _, version := range versions {
		v := reflect.ValueOf(versionMap[version])
		for i := 0; i < v.NumField(); i++ {
			v, _ := v.Field(i).Interface().(containerregistry.Image)
			items = append(items, v.BaseName())
		}
	}

	for _, storages := range csiVersionMap {
		for _, csiV := range storages {
			items = append(items, getImages(csiV)...)
		}
	}

	return items
}

// getImages return images needed by the csi
func getImages(csi *csiVersion) []string {
	images := []string{
		csi.Attacher,
		csi.Provisioner,
		csi.Snapshotter,
		csi.Resizer,
		csi.LivenessProbe,
		csi.NodeRegistrar,
		csi.ClusterRegistrar,
		csi.Driver,
	}

	imagesNeed := make([]string, 0)
	for _, image := range images {
		if image != "" {
			imagesNeed = append(imagesNeed, image)
		}
	}

	return imagesNeed
}

func Versions() []string {
	keys := make([]string, 0, len(versionMap))
	for key := range versionMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

func Get(version string) Components {
	cv, ok := versionMap[version]
	if !ok {
		panic(fmt.Sprintf("the component version definition corresponding to version %s could not be found", version))
	}
	return cv
}
