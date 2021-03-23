package main

import (
	"context"
	"testing"
	"time"
	"wuyuan_exam/internal/biz"
)

var (
	rightDuedate = time.Now().AddDate(1, 1, 1)
	wrongDate    = time.Now().AddDate(-1, 1, 1)
)

func TestAddTask(t *testing.T) {
	taskUsercase, err := initTaskUsercase(1)
	if err != nil {
		panic(err)
	}
	duedate, err := time.Parse("2006-01-02", "2023-01-01")
	if err != nil {
		t.Fatal(err)
	}
	_, err = taskUsercase.Create(context.Background(), &biz.Task{
		Duedate:     duedate,
		Description: "Play with someone",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddTaskBeforeNow(t *testing.T) {
	taskUsercase, err := initTaskUsercase(1)
	if err != nil {
		panic(err)
	}
	_, err = taskUsercase.Create(context.Background(), &biz.Task{
		Duedate:     wrongDate,
		Description: "Play with someone",
	})

	if err != biz.DuedateMustLaterThanNow {
		t.Fatal("DuedateMustLaterThanNow")
	}
}

func TestUpdateTaskDescription(t *testing.T) {
	taskUsercase, err := initTaskUsercase(1)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	id, err := taskUsercase.Create(ctx, &biz.Task{
		Duedate:     rightDuedate,
		Description: "Play with someone",
	})
	if err != nil {
		t.Fatal(err)
	}
	newDescription := "No play"
	err = taskUsercase.UpdateDescription(ctx, id, "No play")
	if err != nil {
		t.Fatal(err)
	}
	tasks, err := taskUsercase.Get(ctx, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if tasks[0].Description != newDescription {
		t.Fatal("Update descrption failed")
	}
}

func TestAddSubTask(t *testing.T) {
	taskUsercase, err := initTaskUsercase(1)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	id, err := taskUsercase.Create(ctx, &biz.Task{
		Duedate:     rightDuedate,
		Description: "A",
	})
	if err != nil {
		t.Fatal(err)
	}

	newID, err := taskUsercase.Create(ctx, &biz.Task{
		ParentId:    id,
		Duedate:     rightDuedate,
		Description: "A",
	})
	if err != nil {
		t.Fatal(err)
	}
	tasks, err := taskUsercase.Get(ctx, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	var subtask biz.Task
	for _, task := range tasks {
		if task.Id == newID {
			subtask = task
			break
		}
	}
	if subtask.ParentId != id {
		t.Fatal("subtask is wrong !")
	}
}

func TestChainComplete(t *testing.T) {
	taskUsercase, err := initTaskUsercase(1)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	id, err := taskUsercase.Create(ctx, &biz.Task{
		Duedate:     rightDuedate,
		Description: "A",
	})

	secondID, err := taskUsercase.Create(ctx, &biz.Task{
		ParentId:    id,
		Duedate:     rightDuedate,
		Description: "A",
	})
	if err != nil {
		t.Fatal(err)
	}

	thirdID, err := taskUsercase.Create(ctx, &biz.Task{
		ParentId:    secondID,
		Duedate:     rightDuedate,
		Description: "A",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = taskUsercase.Complete(ctx, thirdID)
	if err != nil {
		t.Fatal(err)
	}

	tasks, err := taskUsercase.Get(ctx, 0, 10)
	if err != nil {
		t.Fatal(err)
	}

	for _, task := range tasks {
		if task.Status != biz.TaskCompleted {
			t.Fatal("complete faild")
		}
	}
}

func TestReschedule(t *testing.T) {
	taskUsercase, err := initTaskUsercase(1)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	newDuedate := time.Date(2033, 3, 3, 3, 3, 3, 3, time.Local)
	_, err = taskUsercase.Create(ctx, &biz.Task{
		Duedate:     newDuedate,
		Description: "A",
	})
	if err != nil {
		t.Fatal(err)
	}
	tasks, err := taskUsercase.Get(ctx, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if tasks[0].Duedate != newDuedate {
		t.Fatal("Reschedule failed")
	}
}
