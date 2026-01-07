CREATE VIEW course_grades_view AS
SELECT
  c.id as course_id,
  c.name as course_name,
  c.user_id,
  COUNT(cc.id) FILTER (WHERE cc.is_completed = TRUE) as completed_components,
  COUNT(cc.id) as total_components,
  SUM(cc.achieved_score * cc.weight / NULLIF(cc.max_score, 0)) /
    NULLIF(SUM(cc.weight) FILTER (WHERE cc.is_completed = TRUE), 0) as current_grade,
  SUM(cc.weight) as total_weight,
  SUM(cc.weight) FILTER (WHERE cc.is_completed = TRUE) as completed_weight
FROM courses c
LEFT JOIN course_components cc ON cc.course_id = c.id
GROUP BY c.id, c.name, c.user_id;

CREATE VIEW upcoming_course_deadlines AS
SELECT
  c.id as course_id,
  c.user_id,
  c.name as course_name,
  cc.id as component_id,
  cc.type as component_type,
  cc.name as component_name,
  cc.due_date,
  cc.is_completed,
  cc.weight,
  EXTRACT(DAY FROM (cc.due_date - NOW())) as days_remaining
FROM course_components cc
JOIN courses c ON c.id = cc.course_id
WHERE cc.due_date > NOW()
  AND cc.is_completed = FALSE
ORDER BY cc.due_date ASC;
