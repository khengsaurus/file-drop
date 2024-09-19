A simple site offering file-sharing and URL-shortening functionalities

Links

- Referenced by a base62 key 6-chars long
- TTL 1 day

Servers

- 1 server for APIs
- 1 server for streaming files
- 2 servers for URL redirection (LB based on ip_hash)

Rate-limiting

- Token retrieved from server, stored in local-storage for persistence
- Uploading of files and URLs: rate-limited based on a client-side token
- Streaming of files: rate-limited based on client IP address
