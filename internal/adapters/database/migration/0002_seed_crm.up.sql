-- Companies
INSERT INTO companies (name, industry, website, phone, email, country)
VALUES 
('Acme Corp', 'Manufacturing', 'https://acme.com', '123-456-7890', 'info@acme.com', 'USA'),
('Globex Ltd', 'Finance', 'https://globex.com', '987-654-3210', 'contact@globex.com', 'UK');

-- Taxation Details
INSERT INTO taxation_details (country, tax_type, rate, description)
VALUES 
('USA', 'VAT', 7.00, 'Standard VAT'),
('UK', 'GST', 20.00, 'Standard GST');

-- Contacts
INSERT INTO contacts (contact_type, first_name, last_name, email, phone, company_id, country)
VALUES
('individual', 'John', 'Doe', 'john.doe@example.com', '111-222-3333', 1, 'USA'),
('individual', 'Jane', 'Smith', 'jane.smith@example.com', '444-555-6666', 2, 'UK');

-- Leads
INSERT INTO leads (first_name, last_name, email, phone, status, assigned_to, organization_id)
VALUES
('Alice', 'Johnson', 'alice.j@example.com', '555-123-4567', 'New', 101, 201),
('Bob', 'Williams', 'bob.w@example.com', '555-987-6543', 'In Progress', 102, 202);

-- Opportunities
INSERT INTO opportunities (name, description, stage, amount, probability, lead_id, account_id, owner_id)
VALUES
('CRM Upgrade Project', 'Upgrade CRM system for Acme Corp', 'Proposal', 50000.00, 60.0, 1, 1, 1001),
('New Partnership', 'Potential partnership with Globex', 'Negotiation', 75000.00, 40.0, 2, 2, 1002);

-- Activities
INSERT INTO activities (title, description, type, status, due_date, contact_id)
VALUES
('Intro Call', 'Initial introduction call with John Doe', 'Call', 'Completed', GETDATE(), 1),
('Follow-up Meeting', 'Discuss proposal with Jane Smith', 'Meeting', 'Pending', DATEADD(DAY, 3, GETDATE()), 2);

-- Tasks
INSERT INTO tasks (title, description, status, priority, due_date, activity_id)
VALUES
('Send Proposal', 'Email CRM proposal to John Doe', 'Pending', 'High', DATEADD(DAY, 1, GETDATE()), 1),
('Prepare Presentation', 'Prepare deck for meeting with Jane Smith', 'In Progress', 'Medium', DATEADD(DAY, 2, GETDATE()), 2);