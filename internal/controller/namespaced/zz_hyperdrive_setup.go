// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	config "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/hyperdrive/config"
)

// Setup_hyperdrive creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_hyperdrive(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		config.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_hyperdrive creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_hyperdrive(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		config.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
