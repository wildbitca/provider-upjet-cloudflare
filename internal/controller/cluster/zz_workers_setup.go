// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	crontrigger "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/crontrigger"
	customdomain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/customdomain"
	deployment "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/deployment"
	forplatformsdispatchnamespace "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/forplatformsdispatchnamespace"
	kv "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/kv"
	kvnamespace "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/kvnamespace"
	route "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/route"
	script "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/script"
	scriptsubdomain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/scriptsubdomain"
)

// Setup_workers creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_workers(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		crontrigger.Setup,
		customdomain.Setup,
		deployment.Setup,
		forplatformsdispatchnamespace.Setup,
		kv.Setup,
		kvnamespace.Setup,
		route.Setup,
		script.Setup,
		scriptsubdomain.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_workers creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_workers(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		crontrigger.SetupGated,
		customdomain.SetupGated,
		deployment.SetupGated,
		forplatformsdispatchnamespace.SetupGated,
		kv.SetupGated,
		kvnamespace.SetupGated,
		route.SetupGated,
		script.SetupGated,
		scriptsubdomain.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
