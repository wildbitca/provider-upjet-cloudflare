// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	account "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/account"
	dnssettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/dnssettings"
	dnssettingsinternalview "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/dnssettingsinternalview"
	member "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/member"
	subscription "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/subscription"
	token "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/token"
)

// Setup_account creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_account(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		account.Setup,
		dnssettings.Setup,
		dnssettingsinternalview.Setup,
		member.Setup,
		subscription.Setup,
		token.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_account creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_account(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		account.SetupGated,
		dnssettings.SetupGated,
		dnssettingsinternalview.SetupGated,
		member.SetupGated,
		subscription.SetupGated,
		token.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
