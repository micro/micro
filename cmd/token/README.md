# Token

Token is a small http server backed by boltdb and OTP for token generation

## Overview

The token server generates tokens on demand but first verifies via email/OTP

- Runs on port :10001
- `/pass` - generates an OTP for an email
- `/generate` - generates a token for an email/pass
- `/verify` - verifies a token is valid
- `/list` - lists the tokens for an email address
- `/revoke` - revokes any tokens specified

## Usage

The following env vars must be specified

```
# Storage encryption key
STORE_KEY=ewfewfwefwefweewf

# Token encryption key
TOKEN_KEY=fewfwefewfewfwe

# Email domain e.g micro.mu
EMAIL_DOMAIN=micro.mu
```

The token server makes use of gmail smtp relay to send emails. Domain must be set to `micro.mu` for internal usage.
