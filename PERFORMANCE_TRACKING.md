# Response Time Tracking & Admin Reporting

This feature adds response time tracking middleware to monitor Telegram bot handler performance and sends automated daily and weekly reports to admin users.

## 🏗️ Architecture

```
Telegram Message → Bot Router → Metrics Middleware → Handler → UseCase
                        ↓                 ↓
                   Track Time         Store Metrics
                        ↓                 ↓
                   Calculate         PostgreSQL
                        ↓
                   Performance Data → Daily/Weekly Reports → Admin Users
```

## 📊 Components

### 1. **Middleware Layer**
- `internal/middleware/metrics.go` - Records handler response times

### 2. **Performance Domain**
- `internal/domain/performance/metrics.go` - Performance metrics and reports

### 3. **Persistence Layer**
- `internal/persistence/performance/repository.go` - Metrics storage
- `internal/persistence/performance/admin_repository.go` - Admin user management
- `internal/persistence/performance/daemon.go` - Background metrics processing

### 4. **Reports**
- `internal/usecase/reports/daily/` - Daily report generation
- `internal/usecase/reports/weekly/` - Weekly report generation

### 5. **Scheduler**
- `internal/scheduler/` - Automated daily/weekly report scheduling

### 6. **Admin Handler**
- `internal/handler/admin/admin.go` - Admin commands for manual reports

## 🔧 Setup Instructions

### 1. Run Migration
```bash
# Performance metrics and admin tables
make shell
migrate -path migrations -database postgres://youruser:yourpassword@db:5432/task_tracker?sslmode=disable up
```

### 2. Register Admin Users
```sql
-- Add admin users by Telegram user ID
INSERT INTO admin_users (telegram_user_id) VALUES (123456789);
```

### 3. Configure Scheduling
Edit `cmd/service/main.go` to configure report timing:
```go
schedulerConfig := scheduler.SchedulerConfig{
    DailyReportHour:   9,  // 9 AM
    WeeklyReportDay:   time.Sunday,
    WeeklyReportHour:  10, // 10 AM
}
```

## 📈 Metrics Collected

### Per Handler
- Total requests count
- Average response time
- Min/Max response time
- Percentiles (P50, P75, P90, P95, P99)
- Week-over-week trends

### System Wide
- Total system requests
- Overall average response time
- Slowest requests list
- Performance anomalies

## 📅 Reports

### Daily Reports
- Sent at configured hour (default: 9 AM)
- Previous day's performance data
- Handler statistics
- Top 10 slowest requests
- Performance trends vs previous day

### Weekly Reports
- Sent every Sunday at configured hour (default: 10 AM)
- Complete week's performance data
- Week-over-week comparison
- Performance anomalies detection
- Top 20 slowest requests

## 🎛️ Admin Commands

- `/metrics` or `/stats` - Trigger daily report
- `/slowest` - Information about slowest requests
- `/testdaily` - Send test daily report
- `/testweekly` - Send test weekly report
- `/help` - Show admin commands help

## 🔄 Integration Example

```go
// Initialize performance tracking
metricsRepo := performance.NewMetricsRepository(db)
adminRepo := performance.NewAdminRepository(db)
metricsDaemon := performance.NewMetricsDaemon(metricsRepo, 1000, 100, sugar, 5*time.Second)
metricsDaemon.StartWorker(ctx)

// Create middleware
metricsMiddleware := middleware.NewMetricsMiddleware(metricsDaemon, sugar, analyticsDaemon)

// Initialize reporting
dailyUsecase := daily.NewUsecase(metricsRepo, adminRepo, sugar)
weeklyUsecase := weekly.NewUsecase(metricsRepo, adminRepo, sugar)
reportScheduler := scheduler.NewScheduler(dailyUsecase, weeklyUsecase, botApi, schedulerConfig, sugar)

// Add middleware to bot
bot.AddMiddleware(metricsMiddleware)

// Start scheduler
reportScheduler.Start(ctx)
```

## 🗄️ Database Schema

```sql
-- Handler performance metrics
CREATE TABLE handler_performance_metrics (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    handler_name TEXT NOT NULL,
    command TEXT NOT NULL,
    duration_ms BIGINT NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

-- Admin users
CREATE TABLE admin_users (
    id BIGSERIAL PRIMARY KEY,
    telegram_user_id BIGINT UNIQUE NOT NULL,
    registered_at TIMESTAMP DEFAULT NOW()
);
```

## 📝 Notes

- Metrics are collected asynchronously to prevent blocking handlers
- Batch processing (100 metrics or 5-second intervals) for efficiency
- Old metrics are automatically cleaned up (90-day retention default)
- Supports multiple admin users
- Configurable report scheduling
- Detailed performance analytics with percentiles

## 🚀 Monitoring

- Slow requests (>1s) are logged with warnings
- Metrics collection failures are logged
- Report delivery failures are logged
- Scheduler health is monitored

## 🛠️ Troubleshooting

### No Reports Received
1. Check admin users are registered in database
2. Verify scheduler is running correctly
3. Check bot has permissions to send messages to admins
4. Review error logs for connection issues

### High Response Times
1. Check slowest requests in daily reports
2. Review database query performance
3. Monitor system resources
4. Check for external API delays

### Database Performance
1. Ensure indexes are created on timestamp columns
2. Regular cleanup of old metrics
3. Monitor table size and query performance
