package http

import (
	"encoding/json"
	"lms_backend/internal/groups/usecase"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GroupHandler struct {
	uc usecase.GroupUseCase
}

func NewGroupHandler(uc usecase.GroupUseCase) *GroupHandler {
	return &GroupHandler{uc: uc}
}

type UpdateGroupRequest struct {
	Name      string  `json:"name"`
	TeacherID *string `json:"teacher_id,omitempty"`
}

type AddStudentRequest struct {
	StudentID string `json:"student_id"`
}

type ChangeGroupRequest struct {
	GroupID string `json:"group_id"`
}

// UpdateGroup godoc
// @Summary Изменить группу (moderator/admin)
// @Tags Groups
// @Param groupId path string true "Group ID"
// @Param body body UpdateGroupRequest true "Group data"
// @Success 200 {object} map[string]string
// @Router /api/groups/{groupId} [patch]
func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupId")

	var req UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.uc.UpdateGroup(r.Context(), groupID, req.Name, req.TeacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Group updated successfully"})
}

// AddStudentToGroup godoc
// @Summary Добавить ученика в группу (moderator/admin)
// @Tags Groups
// @Param groupId path string true "Group ID"
// @Param body body AddStudentRequest true "Student data"
// @Success 200 {object} map[string]string
// @Router /api/groups/{groupId}/students [post]
func (h *GroupHandler) AddStudentToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupId")

	var req AddStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.uc.AddStudentToGroup(r.Context(), groupID, req.StudentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Student added to group successfully"})
}

// RemoveStudentFromGroup godoc
// @Summary Удалить ученика из группы (moderator/admin)
// @Tags Groups
// @Param groupId path string true "Group ID"
// @Param studentId path string true "Student ID"
// @Success 200 {object} map[string]string
// @Router /api/groups/{groupId}/students/{studentId} [delete]
func (h *GroupHandler) RemoveStudentFromGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupId")
	studentID := chi.URLParam(r, "studentId")

	err := h.uc.RemoveStudentFromGroup(r.Context(), groupID, studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Student removed from group successfully"})
}

// ChangeStudentGroup godoc
// @Summary Изменить группу ученика (moderator/admin)
// @Tags Groups
// @Param studentId path string true "Student ID"
// @Param body body ChangeGroupRequest true "New group data"
// @Success 200 {object} map[string]string
// @Router /api/students/{studentId}/group [patch]
func (h *GroupHandler) ChangeStudentGroup(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentId")

	var req ChangeGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.uc.ChangeStudentGroup(r.Context(), studentID, req.GroupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Student group changed successfully"})
}

// ChangeTeacherGroup godoc
// @Summary Изменить группу учителя (moderator/admin)
// @Tags Groups
// @Param teacherId path string true "Teacher ID"
// @Param body body ChangeGroupRequest true "New group data"
// @Success 200 {object} map[string]string
// @Router /api/teachers/{teacherId}/group [patch]
func (h *GroupHandler) ChangeTeacherGroup(w http.ResponseWriter, r *http.Request) {
	teacherID := chi.URLParam(r, "teacherId")

	var req ChangeGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.uc.ChangeTeacherGroup(r.Context(), teacherID, req.GroupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Teacher group changed successfully"})
}
