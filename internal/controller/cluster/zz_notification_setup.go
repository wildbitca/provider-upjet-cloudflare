// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	policy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/notification/policy"
	policywebhooks "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/notification/policywebhooks"
)

// Setup_notification creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_notification(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		policy.Setup,
		policywebhooks.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_notification creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_notification(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		policy.SetupGated,
		policywebhooks.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
