package tasks

import "gorm.io/gorm"

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

func updateTask(svc TaskApp, updatedTask map[string]interface{}) error {
	var task Task
	if err := svc.Db.Where("id = ?", updatedTask["id"]).First(&task).Error; err != nil {
		return err
	}

	err := svc.Db.Model(&task).Updates(updatedTask).Error
	if err != nil {
		return err
	}
	return nil

}

func getTaskById(svc TaskApp, id int) (Task, error) {
	var task Task
	if err := svc.Db.Where("id = ?", id).First(&task).Error; err != nil {
		return task, err
	}
	if task.ID == 0 {

		return task, &NoTasks{}
	}
	return task, nil

}

func checkUserTask(svc TaskApp, userId, taskId uint) error {
	result := svc.Db.Where("id = ? AND user_id = ?", taskId, userId).Find(&Task{})
	if result.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
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

func allTasks(svc TaskApp, id interface{}) ([]Task, error) {
	var tasks []Task
	if err := svc.Db.Where("user_id = ?", id).Find(&tasks).Error; err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, &NoTasks{}
	}
	return tasks, nil

}

func addFilePath(svc TaskApp, fileName string, taskId int) error {
	err := svc.Db.Model(&Task{}).Where("id = ?", taskId).Update("file", fileName).Error
	if err != nil {
		return err
	}
	return nil
}

func getFilePath(svc TaskApp, taskId int) (string, error) {
	var task Task
	if err := svc.Db.Where("id = ?", taskId).First(&task).Error; err != nil {
		return "", err
	}
	return task.File, nil
}

func deleteFilePath(svc TaskApp, taskId int) error {
	//assign zero value to file path in db
	return addFilePath(svc, "", taskId)
}
