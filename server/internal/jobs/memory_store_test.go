package jobs

import "testing"

func TestGetJob(t *testing.T) {
	job, store := MemoryStoreSetup()

	store.Create(job)

	if _, err := store.Get(job.ID); err != nil {
		t.Errorf("store.GET(%q) did not return the correct job", job.ID)
	}
}

func TestGetJobNotFount(t *testing.T) {
	_, store := MemoryStoreSetup()

	if _, err := store.Get("unknown"); err == nil {
		t.Errorf("Store did not throw expected error")
	}
}

func TestCreateJobAlreadyExists(t *testing.T) {
	job, store := MemoryStoreSetup()

	if err := store.Create(job); err != nil {
		t.Errorf("Store failed to create Job")
	}

	if err := store.Create(job); err == nil {
		t.Errorf("Store did not throw error when creating job twice")
	}
}

func TestUpdateJob(t *testing.T) {
	job, store := MemoryStoreSetup()

	job.Status = JobPending

	store.Create(job)
	job.Status = JobSucceeded
	store.Update(job)

	if jobFromStore, _ := store.Get(job.ID); jobFromStore.Status != JobSucceeded {
		t.Errorf("Job has incorrect status %q", job.Status)
	}
}

func MemoryStoreSetup() (Job, JobStore) {
	store := NewInMemoryStore()

	job := Job{
		ID: "1",
	}

	return job, store
}
