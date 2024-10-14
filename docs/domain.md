# Domain modeling

Basic components of the system include:

- constraints
  - adding
  - removing
  - validation
  - connection with notifications for change notify use case?

- scheduling (the computation)
  - calculation
  - setting parameters
  - presentation of different variants
  - connection to schedule cache to save resources
- schedule computation caching
  - caching
  - handling of periodic removal of cached schedules
  - handling of how many variants should be cached

- profiles
  - user profiles
  - organizations
  - profile images (through some bucket cloud service?)
  - connecting users into organizations

- notifications
  - notifying schedule owners when the user's constraints change
  - notifying the user when the schedule is recalculated
  - notifying the user when the schedule owner changes the constraints
  - can be extended with other notifications

- authentication and authorization
  - providing identity verification for users
  - 2FA and webauthn auth
  -

- frontend (not a microservice candidate)

Scheduling and schedule computation caching can be in one bounded context.
