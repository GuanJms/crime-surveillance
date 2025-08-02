-- Create users with each role
INSERT INTO users (username, password_hash, role)
VALUES 
  ('citizen_jane', 'hash_jane', 'CITIZEN'),
  ('citizen_mark', 'hash_mark', 'CITIZEN'),
  ('citizen_sarah', 'hash_sarah', 'CITIZEN'),
  ('citizen_lee', 'hash_lee', 'CITIZEN'),
  ('patrol_mike', 'hash_mike', 'PATROL'),
  ('patrol_julia', 'hash_julia', 'PATROL'),
  ('dispatcher_amy', 'hash_amy', 'DISPATCHER'),
  ('dispatcher_tom', 'hash_tom', 'DISPATCHER'),
  ('admin_lisa', 'hash_lisa', 'ADMIN');

-- Patrol profiles for PATROL users
INSERT INTO patrol_profile (
  user_id, officer_id, officer_name, street, city, state, latitude, longitude
) VALUES 
(
  (SELECT id FROM users WHERE username = 'patrol_mike'),
  'OFF789',
  'Mike Thompson',
  '1st Avenue',
  'New York',
  'NY',
  40.7306,
  -73.9866
),
(
  (SELECT id FROM users WHERE username = 'patrol_julia'),
  'OFF456',
  'Julia Sanchez',
  '7th Avenue',
  'New York',
  'NY',
  40.7440,
  -73.9950
);

-- NEW crimes
INSERT INTO crime (
  reporter_id, description, status, street, city, state, latitude, longitude
) VALUES 
-- Jane
(
  (SELECT id FROM users WHERE username = 'citizen_jane'),
  'Package theft from porch',
  'NEW',
  'Elm Street',
  'New York',
  'NY',
  40.7310,
  -73.9870
),
-- Mark
(
  (SELECT id FROM users WHERE username = 'citizen_mark'),
  'Loitering outside grocery store',
  'NEW',
  'Oak Avenue',
  'New York',
  'NY',
  40.7320,
  -73.9880
),
-- Sarah
(
  (SELECT id FROM users WHERE username = 'citizen_sarah'),
  'Public urination report',
  'NEW',
  'Delancey Street',
  'New York',
  'NY',
  40.7180,
  -73.9900
),
-- Lee
(
  (SELECT id FROM users WHERE username = 'citizen_lee'),
  'Illegal street racing',
  'NEW',
  'Broadway',
  'New York',
  'NY',
  40.7190,
  -74.0020
);

-- ASSIGNED crimes
INSERT INTO crime (
  reporter_id, patrol_id, description, status, street, city, state, latitude, longitude
) VALUES 
-- Jane + Mike
(
  (SELECT id FROM users WHERE username = 'citizen_jane'),
  (SELECT id FROM users WHERE username = 'patrol_mike'),
  'Disturbance at local park',
  'ASSIGNED',
  'Park Lane',
  'New York',
  'NY',
  40.7330,
  -73.9890
),
-- Mark + Mike
(
  (SELECT id FROM users WHERE username = 'citizen_mark'),
  (SELECT id FROM users WHERE username = 'patrol_mike'),
  'Suspicious vehicle near school',
  'ASSIGNED',
  'School Street',
  'New York',
  'NY',
  40.7340,
  -73.9900
),
-- Sarah + Julia
(
  (SELECT id FROM users WHERE username = 'citizen_sarah'),
  (SELECT id FROM users WHERE username = 'patrol_julia'),
  'Fight outside bar',
  'ASSIGNED',
  'Houston Street',
  'New York',
  'NY',
  40.7250,
  -73.9930
),
-- Lee + Julia
(
  (SELECT id FROM users WHERE username = 'citizen_lee'),
  (SELECT id FROM users WHERE username = 'patrol_julia'),
  'Person brandishing weapon',
  'ASSIGNED',
  'Canal Street',
  'New York',
  'NY',
  40.7200,
  -74.0000
);

-- RESOLVED crimes
INSERT INTO crime (
  reporter_id, patrol_id, description, status, street, city, state, latitude, longitude
) VALUES 
-- Jane + Mike
(
  (SELECT id FROM users WHERE username = 'citizen_jane'),
  (SELECT id FROM users WHERE username = 'patrol_mike'),
  'Resolved: Lost child found at mall',
  'RESOLVED',
  'Mall Drive',
  'New York',
  'NY',
  40.7350,
  -73.9910
),
-- Mark + Mike
(
  (SELECT id FROM users WHERE username = 'citizen_mark'),
  (SELECT id FROM users WHERE username = 'patrol_mike'),
  'Resolved: Noise complaint handled',
  'RESOLVED',
  '5th Avenue',
  'New York',
  'NY',
  40.7360,
  -73.9920
),
-- Sarah + Julia
(
  (SELECT id FROM users WHERE username = 'citizen_sarah'),
  (SELECT id FROM users WHERE username = 'patrol_julia'),
  'Resolved: Trespasser removed from lobby',
  'RESOLVED',
  'Lexington Avenue',
  'New York',
  'NY',
  40.7420,
  -73.9820
),
-- Lee + Julia
(
  (SELECT id FROM users WHERE username = 'citizen_lee'),
  (SELECT id FROM users WHERE username = 'patrol_julia'),
  'Resolved: Drug paraphernalia recovered',
  'RESOLVED',
  'Chrystie Street',
  'New York',
  'NY',
  40.7210,
  -73.9915
);

-- Event logs for user actions
INSERT INTO event_log (user_id, service, message)
VALUES 
  ((SELECT id FROM users WHERE username = 'citizen_jane'), 'report_crime', '{"event":"submitted"}'),
  ((SELECT id FROM users WHERE username = 'citizen_sarah'), 'report_crime', '{"event":"submitted"}'),
  ((SELECT id FROM users WHERE username = 'patrol_mike'), 'patrol_status', '{"status":"BUSY"}'),
  ((SELECT id FROM users WHERE username = 'patrol_julia'), 'patrol_status', '{"status":"AVAILABLE"}'),
  ((SELECT id FROM users WHERE username = 'dispatcher_amy'), 'dispatch', '{"action":"assigned unit"}'),
  ((SELECT id FROM users WHERE username = 'dispatcher_tom'), 'dispatch', '{"action":"cleared case"}'),
  ((SELECT id FROM users WHERE username = 'admin_lisa'), 'admin_action', '{"action":"created user"}');
