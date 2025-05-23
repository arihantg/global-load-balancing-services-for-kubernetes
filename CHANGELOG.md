# Change log:

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
 

## AMKO-1.4.1-beta

### Changed:
  - `GlobalDeploymentPolicy` object structure changes. GDP `v1alpha1` is now deprecated and `v1alpha2` supported.
  - Log fixes for error reporting.

### Added:
  - AMKO support for selecting custom FQDNs for GSs via `GSLBConfig`.
  - AMKO support for `GSLBHostRule`.
  - AMKO support for GslbService properties: TTL, Site Persistence and Custom Health Monitors via `GlobalDeploymentPolicy` and `GSLBHostRule`.
  - AMKO support for adding third party site members via `GSLBHostRule`.

## AMKO-1.4.2

### Added:
  - Support for [federation](docs/AMKO/federation.md). It has to be enabled during installation.
  - A new custom resource `AMKOCluster` for federation configuration.
  - AMKO can now boot up, even if one of the member clusters is unreachable. If the cluster is available later on, AMKO will start it's informers.

### Bugs fixed:
  - Parsing error for TTL and hash mask fields after creation of GslbServices
  - Path based health monitors gets unnecessarily created and then deleted sometimes, even if custom health monitor refs are provided
  - GslbServices unnecessarily updated due to incorrect parsing of site persistence field
  - AMKO doesn't panic if the GSLB leader details couldn't be fetched

## AMKO-1.5.1

### Added:
  - Support for rebooting AMKO if the GSLB leader IP address is updated in the `GSLBConfig` object.

### Bugs fixed:
  - Fixed a status update race between re-sync interval goroutine and modification to `GSLBConfig` object.

## AMKO-1.5.2

### Updated:
  - Base image `photon:4.0` updated with fixes for latest vulnerabilities

## AMKO-1.6.1

### Added:
  - AMKO supports Kubernetes 1.22.
  - Support for multiple AMKO installations.
  - Support for GSLB pool property `priority` via `GlobalDeploymentPolicy` and `GSLBHostRule`.
  - Introduced support for broadcasting AMKO pod `Events` in order to enhance the observability and monitoring aspects.

### Changed:
 - Encode names of all HM objects except HM created for passthrough ingress/routes.


## AMKO-1.7.1

### Added:
  - Support for multiple FQDNs to a single GS using `HostRule` CRD.
  - Support for multi-cluster load balancing for applications deployed in the public cloud.
  - GSLB Monitors settings created by AMKO can now be customized via `GlobalDeploymentPolicy` and `GSLBHostRule`.

### Bugs fixed:
  - Fixed the issue of the AMKO pod does not respond periodically.
  - AMKO now takes into account the priority values given in the `GSLBHostRule` objects.

## AMKO-1.8.1

### Added:
  - Support for Kubernetes 1.24.

### Bugs fixed:
  - Fixed an issue of AMKO updating the health monitors with the wrong ports.

## AMKO-1.8.2

### Bugs fixed:
  - Fixed an issue of AMKO deleting the health monitors references specified as `healthMonitorRefs` via `GlobalDeploymentPolicy` and `GSLBHostRule` from the controller.
  - Fixed few security vulnerabilities in the dependant Golang packages and the base image.

## AMKO-1.9.1

### Added:
  - AMKO now claims support for Kubernetes 1.25 and OCP 4.11.
  - AMKO support for GslbService property Down Response via `GlobalDeploymentPolicy` and `GSLBHostRule`.

### Bugs fixed:
  - Fixed an issue of UUID not being federated across member clusters.

### Known Issues:
  - AMKO fails to create GSLB services with site persistence enabled for controller versions above 22.1.2.

## AMKO-1.10.1

### Added:
  - AMKO now claims support for Kubernetes 1.26 and OCP 4.12.
  - Support for GslbService Pool property `Public IP Address` via `GSLBHostRule`.

### Bugs fixed:
  - Fixed an issue where AMKO was not considering the updates to priority in `GSLBHostRule` objects.
  - Fixed an issue of AMKO not taking into account the TLS annotation in an ingress which was causing
    site persistence configurations to fail.

## AMKO-1.11.1

### Added:
  - AMKO now claims support for Kubernetes 1.27 and OCP 4.13.
  - Support for [Multitenancy](docs/amko_tenancy.md) 

### Changed
  - AKO library has been upgraded. User have to update CRD definitions and AKO to 1.11.1 before upgrading AMKO to 1.11.1.

## AMKO-1.12.1

### Added:
  - AMKO now claims support for Kubernetes 1.28, 1.29 and OCP 4.14

### Changed:
  - PKI Profile used in GSLB Site Persistence can be added via `GlobalDeploymentPolicy` and `GSLBHostRule`

## AMKO-1.13.1

### Added:
  - AMKO now claims support for Kubernetes 1.31 and OCP 4.17
  - AMKO now support multitenancy at namespace level. This feature allows AMKO to map each namespace in Kubernetes/OpenShift cluster uniquely to a tenant in Avi.
  - AMKO now supports OpenShift Routes with `spec.subdomain` field using the [`defaultDomain`](https://github.com/vmware/global-load-balancing-services-for-kubernetes/blob/release-1.13.1/docs/crds/gdp.md?plain=1#L128) field in GDP object.
  - AMKO now has ability to support only control plane health monitoring.

### Bugs fixed:
  - Fixed an issue where AMKO was not adding a new GS poolmember if GSLBHostrule was defined with publicIP
  - Improved logic to remove stale health monitors created by AMKO


