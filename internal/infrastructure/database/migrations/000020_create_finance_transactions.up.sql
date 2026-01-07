CREATE TABLE finance_transactions (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  type VARCHAR(20) NOT NULL,
  amount DECIMAL(10,2) NOT NULL,
  currency VARCHAR(3) DEFAULT 'TRY',
  category VARCHAR(100),
  description TEXT,
  transaction_date DATE DEFAULT CURRENT_DATE,
  related_goal_id INTEGER REFERENCES goals(id) ON DELETE SET NULL,
  is_recurring BOOLEAN DEFAULT FALSE,
  recurrence_config JSONB,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_finance_transactions_user ON finance_transactions(user_id);
CREATE INDEX idx_finance_transactions_date ON finance_transactions(transaction_date);
CREATE INDEX idx_finance_transactions_type ON finance_transactions(type);
CREATE INDEX idx_finance_transactions_category ON finance_transactions(category);
CREATE INDEX idx_finance_transactions_goal ON finance_transactions(related_goal_id);
