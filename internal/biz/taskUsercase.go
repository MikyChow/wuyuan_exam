package biz

import (
	"context"
	"errors"
	"time"
)

type TaskStatus int32

const (
	TaskPending   TaskStatus = 1
	TaskOverdue   TaskStatus = 2
	TaskCompleted TaskStatus = 3
)

var (
	CreateTaskIdMustBeZero  = errors.New("taskId is must be 0")
	DuedateMustLaterThanNow = errors.New("duedate must later than now")
)

type Task struct {
	Id          int64
	ParentId    int64 `json:"parentId"`
	IdPath      string
	Status      TaskStatus `json:"status"`
	Duedate     time.Time  `json:"duedate"`
	Description string     `json:"description"`
}

type TaskRepo interface {
	CreateTask(c context.Context, task *Task) error
	UpdateTask(c context.Context, task *Task) error
	QueryTasks(c context.Context, greateThanId int64, count int64) ([]Task, error)
	GetSiblingTasks(c context.Context, parentId int64, taskId int64) ([]Task, error)
	GetTasks(c context.Context, ids []int64) ([]Task, error)
	GetTask(c context.Context, id int64) (Task, error)
	UpdateTaskDescription(c context.Context, id int64, description string) error
	RescheduleTask(c context.Context, id int64, newDate time.Time) error
	CompleteTask(c context.Context, ids []int64) error
}

type TaskUsercase struct {
	repo        TaskRepo
	idGenerator *IdGenerator
}

func NewTaskUsercase(repo TaskRepo, generator *IdGenerator) *TaskUsercase {
	return &TaskUsercase{repo: repo, idGenerator: generator}
}

func (usercase *TaskUsercase) Create(c context.Context, task *Task) (int64, error) {
	if task.Id != 0 {
		return 0, CreateTaskIdMustBeZero
	}

	if !time.Now().Before(task.Duedate) {
		return 0, DuedateMustLaterThanNow
	}

	newId, err := usercase.idGenerator.Generate()
	if err != nil {
		return 0, err
	}

	task.Id = newId
	task.Status = TaskPending
	return newId, usercase.repo.CreateTask(c, task)
}

func (usercase *TaskUsercase) UpdateDescription(c context.Context, id int64, description string) error {
	return usercase.repo.UpdateTaskDescription(c, id, description)
}

func (usercase *TaskUsercase) Reschedule(c context.Context, id int64, newDate time.Time) error {
	return usercase.repo.RescheduleTask(c, id, newDate)
}

func (usercase *TaskUsercase) Get(c context.Context, greateThanId int64, count int64) ([]Task, error) {
	return usercase.repo.QueryTasks(c, greateThanId, count)
}
func (usercase *TaskUsercase) Complete(c context.Context, id int64) error {
	task, err := usercase.repo.GetTask(c, id)
	if err != nil {
		return err
	}
	completeIds := make([]int64, 0)
	completeIds = append(completeIds, task.Id)

	err = usercase.completeParent(c, task, &completeIds)
	if err != nil {
		return err
	}

	err = usercase.repo.CompleteTask(c, completeIds)
	if err != nil {
		return err
	}

	return nil
}

func (usercase *TaskUsercase) completeParent(c context.Context, task Task, ids *[]int64) error {
	if task.ParentId == 0 {
		return nil
	}

	siblingTasks, err := usercase.repo.GetSiblingTasks(c, task.ParentId, task.Id)
	if err != nil {
		return err
	}
	isAllCompleted := true
	for _, siblingTask := range siblingTasks {
		if siblingTask.Status != TaskCompleted {
			isAllCompleted = false
		}
	}
	if !isAllCompleted {
		return nil
	}
	*ids = append(*ids, task.ParentId)
	parentTask, err := usercase.repo.GetTask(c, task.ParentId)
	if err != nil {
		return err
	}
	if parentTask.ParentId == 0 {
		return nil
	}
	err = usercase.completeParent(c, parentTask, ids)
	if err != nil {
		return err
	}

	return nil
}
