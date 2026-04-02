// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	networkmonitoringconfiguration "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/networkmonitoringconfiguration"
	networkmonitoringrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/networkmonitoringrule"
	transitconnector "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/transitconnector"
	transitsite "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/transitsite"
	transitsiteacl "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/transitsiteacl"
	transitsitelan "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/transitsitelan"
	transitsitewan "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/transitsitewan"
	wangretunnel "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/wangretunnel"
	wanipsectunnel "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/wanipsectunnel"
	wanstaticroute "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/wanstaticroute"
)

// Setup_magic creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_magic(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		networkmonitoringconfiguration.Setup,
		networkmonitoringrule.Setup,
		transitconnector.Setup,
		transitsite.Setup,
		transitsiteacl.Setup,
		transitsitelan.Setup,
		transitsitewan.Setup,
		wangretunnel.Setup,
		wanipsectunnel.Setup,
		wanstaticroute.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_magic creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_magic(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		networkmonitoringconfiguration.SetupGated,
		networkmonitoringrule.SetupGated,
		transitconnector.SetupGated,
		transitsite.SetupGated,
		transitsiteacl.SetupGated,
		transitsitelan.SetupGated,
		transitsitewan.SetupGated,
		wangretunnel.SetupGated,
		wanipsectunnel.SetupGated,
		wanstaticroute.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
