package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	accountCluster "github.com/wildbitca/provider-upjet-cloudflare/config/cluster/account"
	addressCluster "github.com/wildbitca/provider-upjet-cloudflare/config/cluster/address"
	dnsCluster "github.com/wildbitca/provider-upjet-cloudflare/config/cluster/dns"
	zoneCluster "github.com/wildbitca/provider-upjet-cloudflare/config/cluster/zone"
	accountNamespaced "github.com/wildbitca/provider-upjet-cloudflare/config/namespaced/account"
	addressNamespaced "github.com/wildbitca/provider-upjet-cloudflare/config/namespaced/address"
	dnsNamespaced "github.com/wildbitca/provider-upjet-cloudflare/config/namespaced/dns"
	zoneNamespaced "github.com/wildbitca/provider-upjet-cloudflare/config/namespaced/zone"
)

const (
	resourcePrefix = "upjet-cloudflare"
	modulePath     = "github.com/wildbitca/provider-upjet-cloudflare"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("upjet-cloudflare.upbound.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		addressCluster.Configure,
		zoneCluster.Configure,
		dnsCluster.Configure,
		accountCluster.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns the namespaced provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("upjet-cloudflare.m.upbound.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}))

	for _, configure := range []func(provider *ujconfig.Provider){
		addressNamespaced.Configure,
		zoneNamespaced.Configure,
		dnsNamespaced.Configure,
		accountNamespaced.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
