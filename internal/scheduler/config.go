package scheduler

import "time"

type SchedulerConfig struct {
	DailyReportHour  int          // Hour (0-23) to send daily reports
	WeeklyReportHour int          // Hour (0-23) to send weekly reports
	WeeklyReportDay  time.Weekday // Day of week for weekly reports (Sunday = 0)
}

func DefaultConfig() SchedulerConfig {
	return SchedulerConfig{
		DailyReportHour:  9,  // 9 AM
		WeeklyReportDay:  0,  // Sunday
		WeeklyReportHour: 10, // 10 AM
	}
}
