// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	sfuapp "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/calls/sfuapp"
	turnapp "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/calls/turnapp"
)

// Setup_calls creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_calls(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		sfuapp.Setup,
		turnapp.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_calls creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_calls(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		sfuapp.SetupGated,
		turnapp.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
