package usecase

type TeacherDashboardUseCase struct {
}

func NewTeacherDashboardUseCase() *TeacherDashboardUseCase {
	return &TeacherDashboardUseCase{}
}

func (uc *TeacherDashboardUseCase) GetTeacherDashboardData() string {
	return "OK: Teacher Dashboard Module is ready."
}
