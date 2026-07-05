package jobs

type JobStore interface {
	Create(job Job) error
	Get(id string) (Job, error)
	Update(job Job) error
	List() ([]Job, error)
}
