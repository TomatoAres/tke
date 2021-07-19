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

const (
	// LatestVersion is latest version of addon.
	// TODO: bump up to v1.0.3
	LatestVersion = "v1.0.2"

	// Provisioner and ... csi sidecar and controller images name
	Provisioner         = "csi-provisioner"
	Attacher            = "csi-attacher"
	Snapshotter         = "csi-snapshotter"
	LivenessProbe       = "livenessprobe"
	DriverNodeRegistrar = "driver-registrar"
	CSINodeRegistrar    = "csi-node-driver-registrar"
	Resizer             = "csi-resizer"
	RBDDriver           = "rbdplugin"
	CephFSDriver        = "cephfsplugin"
	TencentCBSDriver    = "csi-tencentcloud-cbs"
)

type Components struct {
	CSIOperator containerregistry.Image

	// csi sidecar images
	ProvisionerV042   containerregistry.Image
	AttacherV042      containerregistry.Image
	SnapshotterV041   containerregistry.Image
	LivenessProbeV041 containerregistry.Image
	NodeRegistrarV030 containerregistry.Image
	ProvisionerV101   containerregistry.Image
	AttacherV110      containerregistry.Image
	SnapshotterV110   containerregistry.Image
	LivenessProbeV110 containerregistry.Image
	NodeRegistrarV110 containerregistry.Image
	ProvisionerV120   containerregistry.Image
	SnapshotterV122   containerregistry.Image
	ResizerV050       containerregistry.Image

	// Ceph RBD V0
	CephRBDDriverV030 containerregistry.Image

	// Ceph RBD V1
	CephRBDDriverV100 containerregistry.Image

	// Ceph FS V0, sidecar images are same to Ceph RBD V0
	CephFsDriverV030 containerregistry.Image

	// Ceph FS V1, sidecar images are same to Ceph RBD V1
	CephFsDriverV100 containerregistry.Image

	// Tencent CBS V0
	CbsDriverV021 containerregistry.Image

	// Tencent CBS V1 && V1P1
	CbsDriverV120 containerregistry.Image
	CbsDriverV100 containerregistry.Image // TODO: v1.0.0 is DEPRECATING
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

		ProvisionerV042:   containerregistry.Image{Name: Provisioner, Tag: "v0.4.2"},
		AttacherV042:      containerregistry.Image{Name: Attacher, Tag: "v0.4.2"},
		SnapshotterV041:   containerregistry.Image{Name: Snapshotter, Tag: "v0.4.1"},
		LivenessProbeV041: containerregistry.Image{Name: LivenessProbe, Tag: "v0.4.1"},
		NodeRegistrarV030: containerregistry.Image{Name: DriverNodeRegistrar, Tag: "v0.3.0"},
		CephRBDDriverV030: containerregistry.Image{Name: RBDDriver, Tag: "v0.3.0"},

		ProvisionerV101:   containerregistry.Image{Name: Provisioner, Tag: "v1.0.1"},
		AttacherV110:      containerregistry.Image{Name: Attacher, Tag: "v1.1.0"},
		SnapshotterV110:   containerregistry.Image{Name: Snapshotter, Tag: "v1.1.0"},
		LivenessProbeV110: containerregistry.Image{Name: LivenessProbe, Tag: "v0.4.1"},
		NodeRegistrarV110: containerregistry.Image{Name: CSINodeRegistrar, Tag: "v1.1.0"},
		CephRBDDriverV100: containerregistry.Image{Name: RBDDriver, Tag: "v1.0.0"},

		CephFsDriverV030: containerregistry.Image{Name: CephFSDriver, Tag: "v0.3.0"},
		CephFsDriverV100: containerregistry.Image{Name: CephFSDriver, Tag: "v1.0.0"},

		CbsDriverV021: containerregistry.Image{Name: TencentCBSDriver, Tag: "v0.2.1"},
		CbsDriverV100: containerregistry.Image{Name: TencentCBSDriver, Tag: "v1.0.0"},
		CbsDriverV120: containerregistry.Image{Name: TencentCBSDriver, Tag: "v1.2.0"},
		ResizerV050:   containerregistry.Image{Name: Resizer, Tag: "0.5.0"},
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

	return items
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
