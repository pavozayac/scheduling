# Db models outline

## Schedule

- Id
- Title

## Worker

- Id
- Title
- ScheduleId

## Task

- Id
- Title
- Description
- ScheduleId

## Location

- Id
- Title
- Description
- ScheduleId

## WorkerTimeConstraint

- Id
- WorkerId
- StartTime
- EndTime
- Type (must or cant)

## LocationTimeConstraint

- Id
- WorkerId
- StartTime
- EndTime
- Type (must or cant)

## TaskTimeConstraint

- Id
- WorkerId
- StartTime
- EndTime
- Type (must or cant)

## WorkerTaskConstraint

- Id
- WorkerId
- TaskId
- Type (must or cant)

## LocationTaskConstraint

- Id
- LocationId
- TaskId
- Type (must or cant)
