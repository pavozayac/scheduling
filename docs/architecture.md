# Architecture overview

This is the preliminary list of services which will exist in the domain:

These two the most crucial:

* Constraints
* Computation

And others, supporting the principal services:

* Profile
* Auth
* Notifications (email will likely be here as well)
* Discussion? (for websocket-oriented chats on deliberating schedules)
* Gateway

All synchronous inter-service communication will be handled using gRPC. For example, computation could
invoke constraints for a list of items necessary to calculate a schedule. It could also invoke notification
to send a message once computation is complete.

From the user-facing perspective, all access to resources will be handled by a gateway
service, implementing a GraphQL API over the necessary items. It will also probably expose
a server-sent events endpoint for real-time notifications, as well as a proxy to connect
to the websocket on a discussion service. User authentication will also be handled here, but
for authorization the particular requirements will need to be revised on the basis
of individual services.

The gRPC connections should be operating under a service mesh, probably Linkerd.

There is room for using a message broker for asynchronous communication. Sending notifications, emails
and notifying on the progress of a computation could be the use cases.
