// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	validationoperationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/schema/validationoperationsettings"
	validationschemas "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/schema/validationschemas"
	validationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/schema/validationsettings"
)

// Setup_schema creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_schema(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		validationoperationsettings.Setup,
		validationschemas.Setup,
		validationsettings.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_schema creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_schema(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		validationoperationsettings.SetupGated,
		validationschemas.SetupGated,
		validationsettings.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
