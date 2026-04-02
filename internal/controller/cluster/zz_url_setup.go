// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	normalizationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/url/normalizationsettings"
)

// Setup_url creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_url(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		normalizationsettings.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_url creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_url(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		normalizationsettings.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
