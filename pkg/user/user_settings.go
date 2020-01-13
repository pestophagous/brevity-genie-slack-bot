package user

import ()

type UserSettings interface {
	ChatCapMinutes() int
	SilenceBarrierMinutes() int
	DailyTotalCapMinutes() int

	SetChatCapMinutes(minutes int)
	SetSilenceBarrierMinutes(minutes int)
	SetDailyTotalCapMinutes(minutes int)
}

// This says: a pointer (nil) cast-to-settings is assignable to interface type UserSettings.
// In other words: compile-time enforcement/documentation that settings satisfies UserSettings.
var _ UserSettings = (*settings)(nil)

type settings struct {
	chatCapMinutes        int
	silenceBarrierMinutes int
	dailyTotalCapMinutes  int
}

func SettingsForUser(userId string) UserSettings {
	return &settings{
		chatCapMinutes:        7,
		silenceBarrierMinutes: 20,
		dailyTotalCapMinutes:  30,
	}
}

func (s *settings) ChatCapMinutes() int        { return s.chatCapMinutes }
func (s *settings) SilenceBarrierMinutes() int { return s.silenceBarrierMinutes }
func (s *settings) DailyTotalCapMinutes() int  { return s.dailyTotalCapMinutes }

func (s *settings) SetChatCapMinutes(minutes int)        { panic("not implemented yet") }
func (s *settings) SetSilenceBarrierMinutes(minutes int) { panic("not implemented yet") }
func (s *settings) SetDailyTotalCapMinutes(minutes int)  { panic("not implemented yet") }
