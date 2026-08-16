# Reverse Proxy in Go

In this project, I'm continuing to learn Go.

I decided to build my own reverse proxy service. Even though there are many existing solutions, my main goal is to deepen my understanding of Go and learn how the underlying components work.

I'll primarily use the standard library packages that come with Go, **except** for `httputil/reverseproxy`, since it already implements most of the functionality needed for a reverse proxy.

### Project Goals:
- Forward requests to different backend services
- Thread-safe load balancing
- Active health checking
- Connection pooling with strict timeouts
- Asynchronous logging pipeline
