# okta-gbac
A thin group ownership wrapper around the Okta Developer API

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/okta-gbac)](https://goreportcard.com/report/github.com/adrianosela/okta-gbac)
[![Documentation](https://godoc.org/github.com/adrianosela/okta-gbac?status.svg)](https://godoc.org/github.com/adrianosela/okta-gbac)
[![GitHub issues](https://img.shields.io/github/issues/adrianosela/okta-gbac)](https://github.com/adrianosela/okta-gbac/issues)

### Pre-requisites

- An Okta API Token with `Super Admin` privileges (required for creating and managing all groups)
- The Okta Groups profile must be modified to contain two custom attributes:
	- Data Type: `boolean`, Display Name: `API Managed`, Variable Name: `api_managed`
	- Data Type: `string array`, Display Name: `Owners`, Variable Name: `owners`
