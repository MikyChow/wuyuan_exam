package main

import (
	"context"
	"github.com/google/wire"
	"time"
	"wuyuan_exam/internal/biz"
)

// ProviderSet provide mock repo
var ProviderSet = wire.NewSet(NewTaskRepo)

type taskRepo struct {
	data map[int64]*biz.Task
}

func (repo *taskRepo) CreateTask(_ context.Context, task *biz.Task) error {
	repo.data[task.Id] = task
	return nil
}

func (repo *taskRepo) UpdateTask(_ context.Context, task *biz.Task) error {
	repo.data[task.Id] = task
	return nil

}

func (repo *taskRepo) QueryTasks(_ context.Context, _ int64, _ int64) ([]biz.Task, error) {
	tasks, err := repo.GetAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (repo *taskRepo) GetSiblingTasks(_ context.Context, _ int64, _ int64) ([]biz.Task, error) {
	tasks, err := repo.GetAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (repo *taskRepo) GetTasks(_ context.Context, _ []int64) ([]biz.Task, error) {
	tasks, err := repo.GetAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (repo *taskRepo) GetTask(c context.Context, _ int64) (biz.Task, error) {
	tasks, err := repo.GetAll()
	if err != nil {
		return biz.Task{}, err
	}
	return tasks[0], nil
}

func (repo *taskRepo) UpdateTaskDescription(_ context.Context, id int64, description string) error {
	repo.data[id].Description = description
	return nil
}

func (repo *taskRepo) CompleteTask(_ context.Context, ids []int64) error {
	for _, id := range ids {
		repo.data[id].Status = biz.TaskCompleted
	}
	return nil
}
func (repo *taskRepo) RescheduleTask(_ context.Context, id int64, newDate time.Time) error {
	repo.data[id].Duedate = newDate
	return nil
}

func (repo *taskRepo) GetAll() ([]biz.Task, error) {
	tasks := make([]biz.Task, 0)
	for _, task := range repo.data {
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

//NewTaskRepo returns mock repo
func NewTaskRepo() biz.TaskRepo {
	d := make(map[int64]*biz.Task)
	return &taskRepo{data: d}
}
