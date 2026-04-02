// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	cachereserve "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/cachereserve"
	cachevariants "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/cachevariants"
	dnssec "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/dnssec"
	dnssettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/dnssettings"
	hold "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/hold"
	lockdown "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/lockdown"
	setting "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/setting"
	subscription "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/subscription"
	zone "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/zone"
)

// Setup_zone creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_zone(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		cachereserve.Setup,
		cachevariants.Setup,
		dnssec.Setup,
		dnssettings.Setup,
		hold.Setup,
		lockdown.Setup,
		setting.Setup,
		subscription.Setup,
		zone.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_zone creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_zone(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		cachereserve.SetupGated,
		cachevariants.SetupGated,
		dnssec.SetupGated,
		dnssettings.SetupGated,
		hold.SetupGated,
		lockdown.SetupGated,
		setting.SetupGated,
		subscription.SetupGated,
		zone.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
