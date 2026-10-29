// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	version "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/worker/version"
)

// Setup_worker creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_worker(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		version.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_worker creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_worker(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		version.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager_worker registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager_worker(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		version.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
