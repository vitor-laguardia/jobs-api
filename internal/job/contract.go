package job

import (
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/shared/assert"
	"github.com/vitor-laguardia/jobs-api/internal/user"
)

type RepositoryContract struct {
	NewRepositories func() (Repository, user.Repository)
}

const fakeID = "7c01f5a4-d927-4f4a-98b2-a4bd9cea868b"

func (rc RepositoryContract) Test(t *testing.T) {
	t.Run("Repository.Create job tests", func(t *testing.T) {
		t.Run("correctly persists a job in repository", func(t *testing.T) {
			jobRepo, userRepo := rc.NewRepositories()

			user := user.NewUser("alef", "alef@gmail.com")
			if _, err := userRepo.Create(user); err != nil {
				t.Fatalf("setup failed: could not create job in repository: %v", err)
			}

			job := NewJob("create game", "roleplay for fun", 2, user.ID)
			rowCounter, err := jobRepo.Create(&job)

			assert.Nil(t, err)
			assert.Equal(t, rowCounter, 1, "did not get correct row counter in repository.Create response")
		})
		t.Run("return error when create a job with wrong userID", func(t *testing.T) {
			jobRepo, userRepo := rc.NewRepositories()

			user := user.NewUser("alef", "alef@gmail.com")
			if _, err := userRepo.Create(user); err != nil {
				t.Fatalf("setup failed: could not create job in repository: %v", err)
			}

			job := NewJob("create game", "", 2, fakeID)
			rowCounter, err := jobRepo.Create(&job)

			assert.ErrorIs(t, err, ErrUserNotFound)
			assert.Equal(t, rowCounter, 0, "did not get correct row counter in repository.Create response")
		})
		t.Run("return error when create a job with existent ID", func(t *testing.T) {
			jobRepo, userRepo := rc.NewRepositories()

			user := user.NewUser("alef", "alef@gmail.com")
			if _, err := userRepo.Create(user); err != nil {
				t.Fatalf("setup failed: could not create job in repository: %v", err)
			}

			job := NewJob("create game", "", 2, user.ID)
			if _, err := jobRepo.Create(&job); err != nil {
				t.Fatalf("setup failed: could not create job in repository: %v", err)
			}
			rowCounter, err := jobRepo.Create(&job)

			assert.ErrorIs(t, err, ErrJobAlreadyExists)
			assert.Equal(t, rowCounter, 0, "did not get correct row counter in repository.Create response")
		})

	})

	t.Run("Repository.GetByID job tests", func(t *testing.T) {
		t.Run("sucessfuly get job", func(t *testing.T) {
			jobRepo, userRepo := rc.NewRepositories()

			user := user.NewUser("alef", "alef@gmail.com")
			if _, err := userRepo.Create(user); err != nil {
				t.Fatalf("setup failed: could not create job in repository: %v", err)
			}

			job := NewJob("create game", "roleplay for fun", 2, user.ID)
			if _, err := jobRepo.Create(&job); err != nil {
				t.Fatalf("setup failed: could not create job in repository: %v", err)
			}

			got, err := jobRepo.GetByID(job.ID)

			assert.Equal(t, *got, job, "did not get correct job in repository response")
			assert.Nil(t, err)
		})

		t.Run("return error when job not found", func(t *testing.T) {
			jobRepo, userRepo := rc.NewRepositories()

			user := user.NewUser("alef", "alef@gmail.com")
			if _, err := userRepo.Create(user); err != nil {
				t.Fatalf("setup failed: could not create job in repository: %v", err)
			}

			job := NewJob("create game", "roleplay for fun", 2, user.ID)
			if _, err := jobRepo.Create(&job); err != nil {
				t.Fatalf("setup failed: could not create job in repository: %v", err)
			}

			got, err := jobRepo.GetByID(fakeID)

			assert.ErrorIs(t, err, ErrJobNotFound)
			assert.Nil(t, got)
		})
	})

	t.Run("repository.Update job tests", func(t *testing.T) {
		t.Run("successfully update job fields", func(t *testing.T) {
			tests := []struct {
				name   string
				update func(j *Job)
			}{
				{
					name: "update title only",
					update: func(j *Job) {
						j.Title = "better game"
					},
				},
				{
					name: "update description only",
					update: func(j *Job) {
						j.Description = "new description for testing"
					},
				},
				{
					name: "update status only",
					update: func(j *Job) {
						j.Status = JobStatusDone
					},
				},
				{
					name: "update priority only",
					update: func(j *Job) {
						j.Priority = JobPriorityHigh
					},
				},
				{
					name: "update all editable fields",
					update: func(j *Job) {
						j.Title = "ultimate game edition"
						j.Description = "completely revamped"
						j.Status = JobStatusDone
						j.Priority = JobPriorityHigh
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					jobRepo, userRepo := rc.NewRepositories()

					user := user.NewUser("alef", "alef@gmail.com")
					if _, err := userRepo.Create(user); err != nil {
						t.Fatalf("setup failed: could not create user in repository: %v", err)
					}

					job := NewJob("create game", "roleplay for fun", 2, user.ID)
					if _, err := jobRepo.Create(&job); err != nil {
						t.Fatalf("setup failed: could not create job in repository: %v", err)
					}

					tt.update(&job)

					updatedJob, err := jobRepo.Update(&job)

					assert.Nil(t, err)
					assert.TimeAfter(t, updatedJob.UpdatedAt, job.UpdatedAt, "did not get correct job.UpdatedAt in repository response")

					expectedJob := job
					expectedJob.UpdatedAt = updatedJob.UpdatedAt

					assert.Equal(t, expectedJob, *updatedJob, "did not get correct job in repository response")
				})
			}
		})

		t.Run("cover error scenarios", func(t *testing.T) {
			t.Run("return error when priority out of range", func(t *testing.T) {
				jobRepo, userRepo := rc.NewRepositories()

				user := user.NewUser("alef", "alef@gmail.com")
				if _, err := userRepo.Create(user); err != nil {
					t.Fatalf("setup failed: could not create user in repository: %v", err)
				}

				job := NewJob("create game", "roleplay for fun", 2, user.ID)
				if _, err := jobRepo.Create(&job); err != nil {
					t.Fatalf("setup failed: could not create user in repository: %v", err)
				}

				job.Priority = JobPriority(5)

				got, err := jobRepo.Update(&job)

				assert.ErrorIs(t, err, ErrPriorityOutOfRange)
				assert.Nil(t, got)
			})

			t.Run("return error when wrong status", func(t *testing.T) {
				jobRepo, userRepo := rc.NewRepositories()

				user := user.NewUser("alef", "alef@gmail.com")
				if _, err := userRepo.Create(user); err != nil {
					t.Fatalf("setup failed: could not create user in repository: %v", err)
				}

				job := NewJob("create game", "roleplay for fun", 2, user.ID)
				if _, err := jobRepo.Create(&job); err != nil {
					t.Fatalf("setup failed: could not create user in repository: %v", err)
				}

				job.Status = JobStatus("new status")

				got, err := jobRepo.Update(&job)

				assert.ErrorIs(t, err, ErrInvalidStatus)
				assert.Nil(t, got)
			})
		})
	})

	t.Run("repository.Update job tests", func(t *testing.T) {
		t.Run("successfully update job fields", func(t *testing.T) {
			tests := []struct {
				name   string
				update func(j *Job)
			}{
				{
					name: "update title only",
					update: func(j *Job) {
						j.Title = "better game"
					},
				},
				{
					name: "update description only",
					update: func(j *Job) {
						j.Description = "new description for testing"
					},
				},
				{
					name: "update status only",
					update: func(j *Job) {
						j.Status = JobStatusDone
					},
				},
				{
					name: "update priority only",
					update: func(j *Job) {
						j.Priority = JobPriorityHigh
					},
				},
				{
					name: "update all editable fields",
					update: func(j *Job) {
						j.Title = "ultimate game edition"
						j.Description = "completely revamped"
						j.Status = JobStatusDone
						j.Priority = JobPriorityHigh
					},
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					jobRepo, userRepo := rc.NewRepositories()

					user := user.NewUser("alef", "alef@gmail.com")
					if _, err := userRepo.Create(user); err != nil {
						t.Fatalf("setup failed: could not create job in repository: %v", err)
					}

					job := NewJob("create game", "roleplay for fun", 2, user.ID)
					if _, err := jobRepo.Create(&job); err != nil {
						t.Fatalf("setup failed: could not create job in repository: %v", err)
					}

					tt.update(&job)

					updatedJob, err := jobRepo.Update(&job)

					assert.Nil(t, err)
					assert.TimeAfter(t, updatedJob.UpdatedAt, job.UpdatedAt, "did not get correct job.UpdatedAt in repository response")

					expectedJob := job
					expectedJob.UpdatedAt = updatedJob.UpdatedAt

					assert.Equal(t, expectedJob, *updatedJob, "did not get correct job in repository response")
				})
			}
		})

		t.Run("cover error scenarios", func(t *testing.T) {
			t.Run("return error when priority out of range", func(t *testing.T) {
				jobRepo, userRepo := rc.NewRepositories()

				user := user.NewUser("alef", "alef@gmail.com")
				if _, err := userRepo.Create(user); err != nil {
					t.Fatalf("setup failed: could not create job in repository: %v", err)
				}

				job := NewJob("create game", "roleplay for fun", 2, user.ID)
				if _, err := jobRepo.Create(&job); err != nil {
					t.Fatalf("setup failed: could not create job in repository: %v", err)
				}

				job.Priority = JobPriority(5)

				got, err := jobRepo.Update(&job)

				assert.ErrorIs(t, err, ErrPriorityOutOfRange)
				assert.Nil(t, got)
			})

			t.Run("return error when wrong status", func(t *testing.T) {
				jobRepo, userRepo := rc.NewRepositories()

				user := user.NewUser("alef", "alef@gmail.com")
				if _, err := userRepo.Create(user); err != nil {
					t.Fatalf("setup failed: could not create job in repository: %v", err)
				}

				job := NewJob("create game", "roleplay for fun", 2, user.ID)
				if _, err := jobRepo.Create(&job); err != nil {
					t.Fatalf("setup failed: could not create job in repository: %v", err)
				}

				job.Status = JobStatus("new status")

				got, err := jobRepo.Update(&job)

				assert.ErrorIs(t, err, ErrInvalidStatus)
				assert.Nil(t, got)
			})
		})
	})

	t.Run("repository.Delete job tests", func(t *testing.T) {
		t.Run("succesfully delete a job", func(t *testing.T) {
			jobRepo, userRepo := rc.NewRepositories()

			user := user.NewUser("alef", "alef@gmail.com")
			if _, err := userRepo.Create(user); err != nil {
				t.Fatalf("setup failed: could not create job in repository: %v", err)
			}

			job := NewJob("create game", "roleplay for fun", 2, user.ID)
			if _, err := jobRepo.Create(&job); err != nil {
				t.Fatalf("setup failed: could not create job in repository: %v", err)
			}

			err := jobRepo.Delete(job.ID)

			assert.Nil(t, err)

			got, err := jobRepo.GetByID(job.ID)

			assert.Nil(t, got)
			assert.ErrorIs(t, err, ErrJobNotFound)
		})

		t.Run("return error when delete non-existent job", func(t *testing.T) {
			jobRepo, _ := rc.NewRepositories()
			err := jobRepo.Delete(fakeID)

			assert.ErrorIs(t, err, ErrJobNotFound)
		})
	})
}
