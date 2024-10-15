# Db models outline

## Schedule - Entity

- Id
- Title

## Worker - Entity

- Id
- Title
- ScheduleId

## Task - Entity

- Id
- Title
- Description
- ScheduleId

## Location - Entity

- Id
- Title
- Description
- ScheduleId

## Constraint - Value Object

- WorkerId
- TaskId
- LocationId
- StartTime
- EndTime
- Type
