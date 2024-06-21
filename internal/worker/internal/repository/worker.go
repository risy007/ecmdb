package repository

import (
	"context"
	"github.com/Duke1616/ecmdb/internal/worker/internal/domain"
	"github.com/Duke1616/ecmdb/internal/worker/internal/repository/dao"
)

type WorkerRepository interface {
	CreateWorker(ctx context.Context, req domain.Worker) (int64, error)
	FindByName(ctx context.Context, name string) (domain.Worker, error)
}

func NewWorkerRepository(dao dao.WorkerDAO) WorkerRepository {
	return &workerRepository{
		dao: dao,
	}
}

type workerRepository struct {
	dao dao.WorkerDAO
}

func (repo *workerRepository) CreateWorker(ctx context.Context, req domain.Worker) (int64, error) {
	return repo.dao.CreateWorker(ctx, repo.toEntity(req))
}

func (repo *workerRepository) FindByName(ctx context.Context, name string) (domain.Worker, error) {
	worker, err := repo.dao.FindByName(ctx, name)
	return repo.toDomain(worker), err
}

func (repo *workerRepository) toEntity(req domain.Worker) dao.Worker {
	return dao.Worker{
		Name:  req.Name,
		Topic: req.Topic,
		Desc:  req.Desc,
	}
}

func (repo *workerRepository) toDomain(req dao.Worker) domain.Worker {
	return domain.Worker{
		Name:  req.Name,
		Desc:  req.Desc,
		Topic: req.Topic,
	}
}
