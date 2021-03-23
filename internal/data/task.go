package data

import (
	"context"
	"errors"
	"github.com/lib/pq"
	"time"
	"wuyuan_exam/internal/biz"
)

type taskRepo struct {
	data *Data
}

func NewTaskRepo(data *Data) biz.TaskRepo {
	return &taskRepo{data}
}

func (repo *taskRepo) CreateTask(c context.Context, task *biz.Task) error {
	insertSql :=
		"INSERT INTO tasks(id,parent_id,id_path,status,duedate,description) VALUES($1,$2,$3,$4,$5,$6)"

	_, err := repo.data.ExecContext(
		c,
		insertSql,
		task.Id,
		task.ParentId,
		task.IdPath,
		task.Status,
		task.Duedate,
		task.Description)
	return err
}

func (repo *taskRepo) UpdateTask(c context.Context, task *biz.Task) error {
	updateSql := "UPDATE tasks SET status=$1,duedate=$2,description=$3 WHERE Id=$4"
	_, err := repo.data.ExecContext(c, updateSql, task.Status, task.Duedate, task.Description, task.Id)
	return err
}
func (repo *taskRepo) UpdateTaskDescription(c context.Context, id int64, description string) error {
	updateSql := "UPDATE tasks SET description=$1 WHERE Id=$2"
	_, err := repo.data.ExecContext(c, updateSql, description, id)
	return err
}

func (repo *taskRepo) CompleteTask(c context.Context, ids []int64) error {
	updateSql := "UPDATE tasks SET status=$1 WHERE Id= ANY($2)"
	_, err := repo.data.ExecContext(c, updateSql, biz.TaskCompleted, pq.Int64Array(ids))
	return err
}

func (repo *taskRepo) QueryTasks(c context.Context, greateThanId int64, count int64) ([]biz.Task, error) {
	querySql :=
		"SELECT id,parent_id,id_path,status,duedate,description FROM tasks WHERE id>$1 limit $2"

	rows, err := repo.data.QueryContext(c, querySql, greateThanId, count)
	defer rows.Close()
	tasks := make([]biz.Task, 0)
	for rows.Next() {
		var task biz.Task
		if err := rows.Scan(&task.Id, &task.ParentId, &task.IdPath, &task.Status, &task.Duedate, &task.Description); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, err
}

func (repo *taskRepo) GetSiblingTasks(c context.Context, parentId int64, taskId int64) ([]biz.Task, error) {
	querySql :=
		"SELECT id,parent_id,id_path,status,duedate,description FROM tasks WHERE parent_id=$1 and id!=$2"

	rows, err := repo.data.QueryContext(c, querySql, parentId, taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := make([]biz.Task, 0)
	for rows.Next() {
		var task biz.Task
		if err := rows.Scan(&task.Id, &task.ParentId, &task.IdPath, &task.Status, &task.Duedate, &task.Description); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, err
}

func (repo *taskRepo) GetTasks(c context.Context, ids []int64) ([]biz.Task, error) {
	querySql :=
		"SELECT id,parent_id,id_path,status,duedate,description FROM tasks WHERE id=ANY($1)"
	rows, err := repo.data.QueryContext(c, querySql, pq.Int64Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	tasks := make([]biz.Task, 0)
	for rows.Next() {
		var task biz.Task
		if err := rows.Scan(&task.Id, &task.ParentId, &task.IdPath, &task.Status, &task.Duedate, &task.Description); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, err
}

func (repo *taskRepo) RescheduleTask(_ context.Context, id int64, newDate time.Time) error {
	updateSql := "UPDATE tasks SET duedate=$1 WHERE id=$2"
	_, err := repo.data.Exec(updateSql, newDate, id)
	return err
}

func (repo *taskRepo) GetTask(c context.Context, id int64) (biz.Task, error) {
	tasks, err := repo.GetTasks(c, []int64{id})
	if err != nil {
		return biz.Task{}, err
	}
	if cap(tasks) == 0 {
		return biz.Task{}, errors.New("not found")
	}

	return tasks[0], err
}
