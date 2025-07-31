package data

// TODO: Add methods to the repository interface
type Repository interface {
	UpdateUserRole(userId string, role string) error
	InsertPatrol(*Patrol) (string, error)
	GetAllPatrols(*GetPatrolInfoRequest) ([]*Patrol, error)
	PutPatrolInfo(*UpdatePatrolInfoRequest) error
	PatchPatrolInfo(*UpdatePatrolInfoRequest) error
	UpdatePatrolStatus(id string, status string, conditionStatus *string) error
}

type FastRepository interface {
	UpdatePatrolLocation(id string, location *Location) error
	PersistUpdatePatrolLocation(patrolId string, location *Location) error
	SyncPatrolLocation() error
}

type CrimeRepository interface {
	UpdateCrimePatrolIDStatus(crimeId string, patrolId string, status string) error
}
