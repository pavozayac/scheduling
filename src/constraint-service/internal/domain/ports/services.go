package ports

import "context"

type ConstraintDTO struct {
	ScheduleId     string
	WorkerId       string
	TaskId         string
	LocationId     string
	StartTime      int
	EndTime        int
	ConstraintType string
}

type ScheduleDTO struct {
	Id          string
	Title       string
	Constraints []ConstraintDTO
}

type LocationDTO struct {
	Id          string
	ScheduleId  string
	Name        string
	Description string
}

type TaskDTO struct {
	Id          string
	ScheduleId  string
	Title       string
	Description string
}

type WorkerDTO struct {
	Id         string
	ScheduleId string
	FirstName  string
	LastName   string
}

type ScheduleService interface {
	GetSchedule(context.Context, string) (ScheduleDTO, error)
	CreateOrModifySchedule(context.Context, ScheduleDTO) (*ScheduleDTO, error)
	RemoveSchedule(context.Context, string) error
}

type LocationService interface {
	GetLocation(context.Context, string) (LocationDTO, error)
	CreateOrModifyLocation(context.Context, LocationDTO) (*LocationDTO, error)
	RemoveLocation(context.Context, string) error
}

type TaskService interface {
	GetTask(context.Context, string) (TaskDTO, error)
	CreateOrModifyTask(context.Context, TaskDTO) (*TaskDTO, error)
	RemoveTask(context.Context, string) error
}

type WorkerService interface {
	GetWorker(context.Context, string) (WorkerDTO, error)
	CreateOrModifyWorker(context.Context, WorkerDTO) (*WorkerDTO, error)
	RemoveWorker(context.Context, string) error
}
