// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	filter "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/filter"
	healthcheck "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/healthcheck"
	image "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/image"
	list "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/list"
	organization "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/organization"
	queue "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/queue"
	ruleset "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/ruleset"
	snippet "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/snippet"
	snippets "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/snippets"
	stream "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/stream"
	user "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/user"
	worker "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/worker"
	workflow "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/workflow"
)

// Setup_cloudflare creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_cloudflare(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		filter.Setup,
		healthcheck.Setup,
		image.Setup,
		list.Setup,
		organization.Setup,
		queue.Setup,
		ruleset.Setup,
		snippet.Setup,
		snippets.Setup,
		stream.Setup,
		user.Setup,
		worker.Setup,
		workflow.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_cloudflare creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_cloudflare(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		filter.SetupGated,
		healthcheck.SetupGated,
		image.SetupGated,
		list.SetupGated,
		organization.SetupGated,
		queue.SetupGated,
		ruleset.SetupGated,
		snippet.SetupGated,
		snippets.SetupGated,
		stream.SetupGated,
		user.SetupGated,
		worker.SetupGated,
		workflow.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
