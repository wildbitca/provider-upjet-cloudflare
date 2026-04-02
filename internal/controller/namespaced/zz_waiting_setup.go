// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	room "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/waiting/room"
	roomevent "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/waiting/roomevent"
	roomrules "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/waiting/roomrules"
	roomsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/waiting/roomsettings"
)

// Setup_waiting creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_waiting(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		room.Setup,
		roomevent.Setup,
		roomrules.Setup,
		roomsettings.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_waiting creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_waiting(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		room.SetupGated,
		roomevent.SetupGated,
		roomrules.SetupGated,
		roomsettings.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
