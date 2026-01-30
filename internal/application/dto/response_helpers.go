package dto

import "fmt"

// ErrorUserResponse creates an error response for User operations
func ErrorUserResponse(err error) *UserResponse {
	return &UserResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessUserResponse creates a success response for User operations
func SuccessUserResponse(user *UserDTO) *UserResponse {
	return &UserResponse{
		Success: true,
		User:    user,
	}
}

// ErrorUsersResponse creates an error response for Users list operations
func ErrorUsersResponse(err error) *UsersResponse {
	return &UsersResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessUsersResponse creates a success response for Users list operations
func SuccessUsersResponse(users []*UserDTO) *UsersResponse {
	return &UsersResponse{
		Success: true,
		Users:   users,
	}
}

// ErrorSessionResponse creates an error response for Session operations
func ErrorSessionResponse(err error) *SessionResponse {
	return &SessionResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessSessionResponse creates a success response for Session operations
func SuccessSessionResponse(session *SessionDTO) *SessionResponse {
	return &SessionResponse{
		Success: true,
		Session: session,
	}
}

// ErrorSessionsResponse creates an error response for Sessions list operations
func ErrorSessionsResponse(err error) *SessionsResponse {
	return &SessionsResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessSessionsResponse creates a success response for Sessions list operations
func SuccessSessionsResponse(sessions []*SessionDTO) *SessionsResponse {
	return &SessionsResponse{
		Success:  true,
		Sessions: sessions,
	}
}

// ErrorMessageResponse creates an error response for Message operations
func ErrorMessageResponse(err error) *MessagesResponse {
	return &MessagesResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessMessagesResponse creates a success response for Message list operations
func SuccessMessagesResponse(messages []*MessageDTO) *MessagesResponse {
	return &MessagesResponse{
		Success:  true,
		Messages: messages,
	}
}

// ErrorTaskResponse creates an error response for Task operations
func ErrorTaskResponse(err error) *TasksResponse {
	return &TasksResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessTasksResponse creates a success response for Task list operations
func SuccessTasksResponse(tasks []*TaskDTO) *TasksResponse {
	return &TasksResponse{
		Success: true,
		Tasks:   tasks,
	}
}

// ErrorSkillResponse creates an error response for Skill operations
func ErrorSkillResponse(err error) *SkillResponse {
	return &SkillResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessSkillResponse creates a success response for Skill operations
func SuccessSkillResponse(skill *SkillDTO) *SkillResponse {
	return &SkillResponse{
		Success: true,
		Skill:   skill,
	}
}

// ErrorSkillsResponse creates an error response for Skills list operations
func ErrorSkillsResponse(err error) *SkillsResponse {
	return &SkillsResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessSkillsResponse creates a success response for Skills list operations
func SuccessSkillsResponse(skills []*SkillDTO) *SkillsResponse {
	return &SkillsResponse{
		Success: true,
		Skills:  skills,
	}
}

// ErrorScheduleResponse creates an error response for Schedule operations
func ErrorScheduleResponse(err error) *ScheduleResponse {
	return &ScheduleResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessScheduleResponse creates a success response for Schedule operations
func SuccessScheduleResponse(schedule *ScheduleDTO) *ScheduleResponse {
	return &ScheduleResponse{
		Success:  true,
		Schedule: schedule,
	}
}

// ErrorSchedulesResponse creates an error response for Schedules list operations
func ErrorSchedulesResponse(err error) *SchedulesResponse {
	return &SchedulesResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessSchedulesResponse creates a success response for Schedules list operations
func SuccessSchedulesResponse(schedules []*ScheduleDTO) *SchedulesResponse {
	return &SchedulesResponse{
		Success:   true,
		Schedules: schedules,
	}
}

// ErrorSkillExecutionResponse creates an error response for SkillExecution operations
func ErrorSkillExecutionResponse(err error) *SkillExecutionResponse {
	return &SkillExecutionResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessSkillExecutionResponse creates a success response for SkillExecution operations
func SuccessSkillExecutionResponse(output string) *SkillExecutionResponse {
	return &SkillExecutionResponse{
		Success: true,
		Output:  output,
	}
}

// ErrorSendMessageResponse creates an error response for SendMessage operations
func ErrorSendMessageResponse(err error) *SendMessageResponse {
	return &SendMessageResponse{
		Success: false,
		Error:   fmt.Sprintf("operation failed: %v", err),
	}
}

// SuccessSendMessageResponse creates a success response for SendMessage operations
func SuccessSendMessageResponse(message *MessageDTO, messages []*MessageDTO) *SendMessageResponse {
	return &SendMessageResponse{
		Success:  true,
		Message:  message,
		Messages: messages,
	}
}
