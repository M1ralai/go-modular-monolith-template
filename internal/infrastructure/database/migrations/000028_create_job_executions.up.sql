-- Job execution tracking with deduplication support
CREATE TABLE job_executions (
    id SERIAL PRIMARY KEY,
    job_name VARCHAR(100) NOT NULL,
    entity_id INTEGER,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'running', 'completed', 'failed')),
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT,
    result JSONB,
    duration_ms INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Unique constraint: one job per entity per day (for deduplication)
CREATE UNIQUE INDEX idx_job_executions_unique 
ON job_executions(job_name, entity_id, DATE(started_at))
WHERE status = 'running' OR (completed_at IS NOT NULL AND DATE(completed_at) = DATE(started_at));

-- Index for active jobs
CREATE INDEX idx_job_executions_active 
ON job_executions(job_name, entity_id) 
WHERE completed_at IS NULL AND status = 'running';

-- Index for cleanup of old completed jobs
CREATE INDEX idx_job_executions_completed 
ON job_executions(completed_at) 
WHERE completed_at IS NOT NULL;

CREATE INDEX idx_job_executions_job_name ON job_executions(job_name);
CREATE INDEX idx_job_executions_status ON job_executions(status);
CREATE INDEX idx_job_executions_started_at ON job_executions(started_at DESC);

-- Function to clean up old job executions (older than 7 days)
CREATE OR REPLACE FUNCTION cleanup_old_job_executions()
RETURNS void AS $$
BEGIN
    DELETE FROM job_executions 
    WHERE completed_at IS NOT NULL 
    AND completed_at < NOW() - INTERVAL '7 days';
END;
$$ LANGUAGE plpgsql;
