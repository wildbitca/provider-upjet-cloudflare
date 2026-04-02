// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	job "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/logpush/job"
	ownershipchallenge "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/logpush/ownershipchallenge"
)

// Setup_logpush creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_logpush(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		job.Setup,
		ownershipchallenge.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_logpush creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_logpush(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		job.SetupGated,
		ownershipchallenge.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
