# api-streaming-availability

Stan's transcoder pipeline emits events whenever a video segment becomes available for playback. Downstream CDN origin servers poll a manifest service to discover which segments are ready to serve to players.

This is the Go HTTP service that sits between the transcoder and the CDN.
