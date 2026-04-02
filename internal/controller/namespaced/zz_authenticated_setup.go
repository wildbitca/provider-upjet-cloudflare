// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	originpulls "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/authenticated/originpulls"
	originpullscertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/authenticated/originpullscertificate"
	originpullshostnamecertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/authenticated/originpullshostnamecertificate"
	originpullssettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/authenticated/originpullssettings"
)

// Setup_authenticated creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_authenticated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		originpulls.Setup,
		originpullscertificate.Setup,
		originpullshostnamecertificate.Setup,
		originpullssettings.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_authenticated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_authenticated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		originpulls.SetupGated,
		originpullscertificate.SetupGated,
		originpullshostnamecertificate.SetupGated,
		originpullssettings.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
