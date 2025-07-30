package data

// TODO: Add methods to the repository interface
type Repository interface {
	InsertPatrol(*Patrol) (string, error)
	GetAllPatrols(*GetPatrolInfoRequest) ([]*Patrol, error)
	PutPatrolInfo(*UpdatePatrolInfoRequest) error
	PatchPatrolInfo(*UpdatePatrolInfoRequest) error
	UpdatePatrolStatus(id string, status string) error
}

type FastRepository interface {
	UpdatePatrolLocation(id string, location *Location) error
	PersistUpdatePatrolLocation(patrolId string, location *Location) error
	SyncPatrolLocation() error
}
