fetch
=====

Simple example on how to fetch multiple URLs and combine them in one response to
the client.

Example usage:

	curl -X POST 'http://localhost:8080' \
	-d '{ "urls" : ["http://google.com", "http://facebook.com", "http://twitter.com"] }'

Code was quickly hacked together and can be improved :).
