package tasks

type NoTasks struct{}

func (m *NoTasks) Error() string {
	return "no tasks in db"
}
func createTask(svc TaskApp, task Task) error {
	result := svc.Db.Create(&task)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func updateTask() {

}

func getTaskById(svc TaskApp, id int) (Task, error) {
	var task Task
	if err := svc.Db.Where("id = ?", id).First(&task).Error; err != nil {
		return task, err
	}
	return task, nil

}

func deleteTask(svc TaskApp, id int) error {
	var task Task
	if err := svc.Db.Where("id = ?", id).First(&task).Error; err != nil {
		return err
	}

	result := svc.Db.Delete(&task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func allTasks(svc TaskApp) ([]Task, error) {
	var tasks []Task
	if err := svc.Db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, &NoTasks{}
	}

	return tasks, nil
}
