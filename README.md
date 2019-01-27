# DDNS-Proxy
Dynamic DNS Proxy

The goal of this project is to provide an interface for creating/updating DNS records to act as a Dynamic DNS (DDNS)
Currently the goal is to support OpenWRT Dynamic DNS features to use this service.

Supported DNS providers:
- Cloudflare

Attention: A lot of service providers use TTL to enforce caching DNS results to increase speed, this might cause interruptions or loss of conectivity on one of the nodes if the ip address gets updated during a short period.