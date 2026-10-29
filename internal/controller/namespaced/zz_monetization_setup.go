// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	rules "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/monetization/rules"
)

// Setup_monetization creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_monetization(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		rules.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_monetization creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_monetization(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		rules.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager_monetization registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager_monetization(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		rules.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
