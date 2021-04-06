package client

import "github.com/urfave/cli/v2"

// Flags common to all clients
var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:  "local",
		Usage: "Enable local only development: Defaults to true.",
	},
	&cli.BoolFlag{
		Name:    "enable_acme",
		Usage:   "Enables ACME support via Let's Encrypt. ACME hosts should also be specified.",
		EnvVars: []string{"MICRO_ENABLE_ACME"},
	},
	&cli.StringFlag{
		Name:    "acme_hosts",
		Usage:   "Comma separated list of hostnames to manage ACME certs for",
		EnvVars: []string{"MICRO_ACME_HOSTS"},
	},
	&cli.StringFlag{
		Name:    "acme_provider",
		Usage:   "The provider that will be used to communicate with Let's Encrypt. Valid options: autocert, certmagic",
		EnvVars: []string{"MICRO_ACME_PROVIDER"},
	},
	&cli.BoolFlag{
		Name:    "enable_tls",
		Usage:   "Enable TLS support. Expects cert and key file to be specified",
		EnvVars: []string{"MICRO_ENABLE_TLS"},
	},
	&cli.StringFlag{
		Name:    "tls_cert_file",
		Usage:   "Path to the TLS Certificate file",
		EnvVars: []string{"MICRO_TLS_CERT_FILE"},
	},
	&cli.StringFlag{
		Name:    "tls_key_file",
		Usage:   "Path to the TLS Key file",
		EnvVars: []string{"MICRO_TLS_KEY_FILE"},
	},
	&cli.StringFlag{
		Name:    "tls_client_ca_file",
		Usage:   "Path to the TLS CA file to verify clients against",
		EnvVars: []string{"MICRO_TLS_CLIENT_CA_FILE"},
	},
}
