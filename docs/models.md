# Modeling

## Domain modeling

### Schedule - Entity, Aggregate Root

- Identity
- Title
- Workers
- Tasks
- Locations
- Constraints

---
Even though workers, tasks and locations could be logically one and the same, leaving them separate could be a good idea for future development, where a worker could have an associated user account, who would flag their own constraints.

### Worker - Entity, Aggregate

- Identity
- First name
- Last name
- ScheduleId

### Task - Entity, Aggregate

- Identity
- Title
- Description
- ScheduleId

### Location - Entity, Aggregate

- Identity
- Title
- Description
- ScheduleId

### Constraint - Value Object

- ScheduleId
- WorkerId
- TaskId
- LocationId
- StartTime
- EndTime
- Type

## Database modeling

### Schedule

- Id
- Title

### Worker

- Id
- Title
- ScheduleId

### Task

- Id
- Title
- Story
- ScheduleId

### Location

- Id
- Title
- Story
- ScheduleId

### Constraint

- ScheduleId
- WorkerId
- TaskId
- LocationId
- StartTime
- EndTime
- Type
