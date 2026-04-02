// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	bucket "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucket"
	bucketcors "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucketcors"
	bucketeventnotification "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucketeventnotification"
	bucketlifecycle "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucketlifecycle"
	bucketlock "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucketlock"
	bucketsippy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucketsippy"
	customdomain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/customdomain"
	manageddomain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/manageddomain"
)

// Setup_r2 creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_r2(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		bucket.Setup,
		bucketcors.Setup,
		bucketeventnotification.Setup,
		bucketlifecycle.Setup,
		bucketlock.Setup,
		bucketsippy.Setup,
		customdomain.Setup,
		manageddomain.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_r2 creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_r2(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		bucket.SetupGated,
		bucketcors.SetupGated,
		bucketeventnotification.SetupGated,
		bucketlifecycle.SetupGated,
		bucketlock.SetupGated,
		bucketsippy.SetupGated,
		customdomain.SetupGated,
		manageddomain.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
