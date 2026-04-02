// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	hostname "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/custom/hostname"
	hostnamefallbackorigin "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/custom/hostnamefallbackorigin"
	pages "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/custom/pages"
	ssl "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/custom/ssl"
)

// Setup_custom creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_custom(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		hostname.Setup,
		hostnamefallbackorigin.Setup,
		pages.Setup,
		ssl.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_custom creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_custom(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		hostname.SetupGated,
		hostnamefallbackorigin.SetupGated,
		pages.SetupGated,
		ssl.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
