-- Remove seeded Opportunities
DELETE FROM opportunities
WHERE name IN ('CRM Deal', 'Follow-up Opportunity');

-- Remove seeded Leads
DELETE FROM leads
WHERE email IN ('lead1@example.com', 'lead2@example.com');

-- Remove seeded Contacts
DELETE FROM contacts
WHERE email IN ('john.doe@example.com', 'jane.smith@example.com');

-- Remove seeded Companies
DELETE FROM companies
WHERE name IN ('TechCorp', 'InnovateLtd');

-- Remove seeded Tasks
DELETE FROM tasks
WHERE title IN ('Send Proposal', 'Prepare Presentation');

-- Remove seeded Activities
DELETE FROM activities
WHERE title IN ('Intro Call', 'Follow-up Meeting');
