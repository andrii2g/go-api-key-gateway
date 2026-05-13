# API Key Format

`ak_{app}_{public_key}_{secret}`

- `app`: lowercase `a-z0-9`, length 1 to 3
- `public_key`: 16 characters from `0123456789ABCDEFGHJKMNPQRSTVWXYZ`
- `secret`: Base64URL without padding, decoded length at least 32 bytes
