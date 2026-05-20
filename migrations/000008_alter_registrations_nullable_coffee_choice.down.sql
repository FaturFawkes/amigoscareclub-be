UPDATE registrations SET coffee_choice = '' WHERE coffee_choice IS NULL;
ALTER TABLE registrations ALTER COLUMN coffee_choice SET NOT NULL;
