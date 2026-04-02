// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	credentialcheck "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/leaked/credentialcheck"
	credentialcheckrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/leaked/credentialcheckrule"
)

// Setup_leaked creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_leaked(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		credentialcheck.Setup,
		credentialcheckrule.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_leaked creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_leaked(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		credentialcheck.SetupGated,
		credentialcheckrule.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
