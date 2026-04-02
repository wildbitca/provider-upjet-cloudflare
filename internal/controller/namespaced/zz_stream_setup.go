// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	audiotrack "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/audiotrack"
	captionlanguage "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/captionlanguage"
	download "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/download"
	key "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/key"
	liveinput "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/liveinput"
	watermark "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/watermark"
	webhook "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/webhook"
)

// Setup_stream creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_stream(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		audiotrack.Setup,
		captionlanguage.Setup,
		download.Setup,
		key.Setup,
		liveinput.Setup,
		watermark.Setup,
		webhook.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_stream creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_stream(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		audiotrack.SetupGated,
		captionlanguage.SetupGated,
		download.SetupGated,
		key.SetupGated,
		liveinput.SetupGated,
		watermark.SetupGated,
		webhook.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
