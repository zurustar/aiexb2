-- Create Keycloak database
CREATE DATABASE keycloak;

-- Grant privileges to esms_user
GRANT ALL PRIVILEGES ON DATABASE keycloak TO esms_user;
