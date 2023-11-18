# Database Code Generator

Query used to get all the required details: 

```
twitter_clone=# SELECT pg_stat_user_tables.relname                         AS table_name,
(SELECT pg_description.description
FROM pg_description
WHERE pg_stat_user_tables.relid = pg_description.objoid
AND pg_description.objsubid = 0)                 AS table_comment,
information_schema.columns.column_name              AS column_name,
information_schema.columns.column_default           AS column_default,
pg_description.description                          AS column_comment,
information_schema.columns.table_schema             AS table_schema,
information_schema.columns.data_type                AS column_data_type,
information_schema.columns.character_maximum_length AS char_len,
CASE information_schema.columns.is_nullable
WHEN 'NO' THEN FALSE
WHEN 'YES' THEN TRUE
END                                             AS nullable,
information_schema.columns.numeric_precision || ',' ||
information_schema.columns.numeric_scale            AS numeric_length
FROM pg_stat_user_tables
LEFT JOIN information_schema.columns
ON pg_stat_user_tables.relname = information_schema.columns.table_name
LEFT JOIN pg_description
ON pg_description.objoid = pg_stat_user_tables.relid
AND pg_description.objsubid = information_schema.columns.ordinal_position
WHERE information_schema.columns.table_schema NOT IN ('pg_catalog', 'information_schema')
AND table_name != 'schema_migrations'
ORDER BY table_schema, table_name, information_schema.columns.ordinal_position;
```
Result (formatted for Markdown): 

| table_name        | table_comment                                  | column_name           | column_default                           | column_comment                                                                                 | table_schema | column_data_type            | char_len | nullable | numeric_length |
|-------------------|------------------------------------------------|-----------------------|------------------------------------------|------------------------------------------------------------------------------------------------|--------------|-----------------------------|----------|----------|----------------|
| followers         | Table storing follower/following relationships | follower_id           |                                          | Foreign key referencing users_credentials.user_id for the follower                             | app          | uuid                        |          | f        |                |
| followers         | Table storing follower/following relationships | following_id          |                                          | Foreign key referencing users_credentials.user_id for the following                            | app          | uuid                        |          | f        |                |
| likes             | Table storing likes on tweets                  | like_id               | uuid_generate_v4()                       | Unique identifier for likes                                                                    | app          | uuid                        |          | f        |                |
| likes             | Table storing likes on tweets                  | user_id               |                                          | Foreign key referencing users_credentials.user_id for the user who liked                       | app          | uuid                        |          | t        |                |
| likes             | Table storing likes on tweets                  | tweet_id              |                                          | Foreign key referencing tweets.tweet_id for the liked tweet                                    | app          | uuid                        |          | t        |                |
| likes             | Table storing likes on tweets                  | created_at            | CURRENT_TIMESTAMP                        | Timestamp of like creation                                                                     | app          | timestamp without time zone |          | t        |                |
| tweets            | Table storing tweets and comments              | tweet_id              | uuid_generate_v4()                       | Unique identifier for tweets and comments                                                      | app          | uuid                        |          | f        |                |
| tweets            | Table storing tweets and comments              | user_id               |                                          | Foreign key referencing users_credentials.user_id                                              | app          | uuid                        |          | t        |                |
| tweets            | Table storing tweets and comments              | parent_tweet_id       |                                          | Foreign key referencing tweets.tweet_id for comments                                           | app          | uuid                        |          | t        |                |
| tweets            | Table storing tweets and comments              | tweet_text            |                                          | Text content of the tweet or comment                                                           | app          | text                        |          | f        |                |
| tweets            | Table storing tweets and comments              | created_at            | CURRENT_TIMESTAMP                        | Timestamp of tweet creation                                                                    | app          | timestamp without time zone |          | t        |                |
| user_profiles     | Table storing user profile information         | user_id               |                                          | Foreign key referencing users_credentials.user_id                                              | app          | uuid                        |          | f        |                |
| user_profiles     | Table storing user profile information         | description           |                                          | User profile description                                                                       | app          | text                        |          | t        |                |
| user_profiles     | Table storing user profile information         | profile_picture_url   |                                          | URL of the profile picture for the user                                                        | app          | character varying           | 255      | t        |                |
| user_profiles     | Table storing user profile information         | followers_count       | 0                                        | Count of followers for the user                                                                | app          | integer                     |          | t        | 32,0           |
| user_profiles     | Table storing user profile information         | bio                   |                                          | User biography                                                                                 | app          | text                        |          | t        |                |
| user_profiles     | Table storing user profile information         | website_url           |                                          | Website URL of the user                                                                        | app          | character varying           | 255      | t        |                |
| user_profiles     | Table storing user profile information         | created_at            | CURRENT_TIMESTAMP                        | Timestamp of profile creation                                                                  | app          | timestamp without time zone |          | t        |                |
| users_credentials | Table storing user login information           | user_id               | uuid_generate_v4()                       | Unique identifier for users                                                                    | app          | uuid                        |          | f        |                |
| users_credentials | Table storing user login information           | username              |                                          | Username for login                                                                             | app          | character varying           | 50       | f        |                |
| users_credentials | Table storing user login information           | email                 |                                          | Email for login                                                                                | app          | character varying           | 100      | f        |                |
| users_credentials | Table storing user login information           | password_hash         |                                          | Hashed password for login                                                                      | app          | character varying           | 100      | f        |                |
| app_log           | Table to store application logs                | id                    | nextval('logs.app_log_id_seq'::regclass) | Unique log ID, Primary Key, Auto-incrementing number                                           | logs         | bigint                      |          | f        | 64,0           |
| app_log           | Table to store application logs                | log_time              | (now() AT TIME ZONE 'utc'::text)         | UTC time when the log occurred. Defaults to "now()"                                            | logs         | timestamp without time zone |          | f        |                |
| app_log           | Table to store application logs                | log_level             | 'info'::character varying                | Severity level - INFO, WARNING, ERROR, PANIC etc.                                              | logs         | character varying           | 16       | f        |                |
| app_log           | Table to store application logs                | service_name          | 'def_svc'::character varying             | Name of the service that sent this log (defaults to "def_svc")                                 | logs         | character varying           | 64       | f        |                |
| app_log           | Table to store application logs                | service_instance_name | 'def_svc_instance'::character varying    | Name of the service instance (or Pod Name) that sent this log (defaults to "def_svc_instance") | logs         | character varying           | 64       | f        |                |
| app_log           | Table to store application logs                | code                  | '000000'::character varying              | Unique Code of the log (LMID) - to be sent by the caller (defaults to "000000")                | logs         | character varying           | 16       | f        |                |
| app_log           | Table to store application logs                | msg                   | '_no_msg_supplied_'::text                | Actual log message as text                                                                     | logs         | text                        |          | f        |                |
| app_log           | Table to store application logs                | more_data             | '{}'::jsonb                              | Anything else that needs to be saved alongside the log entry                                   | logs         | jsonb                       |          | f        |                |
| key_value         |                                                | key                   |                                          |                                                                                                | public       | character varying           | 512      | f        |                |
| key_value         |                                                | value                 |                                          |                                                                                                | public       | text                        |          | f        |                |
| key_value         |                                                | created_at            | (now() AT TIME ZONE 'utc'::text)         |                                                                                                | public       | timestamp without time zone |          | f        |                |
| key_value         |                                                | updated_at            | (now() AT TIME ZONE 'utc'::text)         |                                                                                                | public       | timestamp without time zone |          | f        |                |
| key_value         |                                                | more_data             | '{}'::jsonb                              |                                                                                                | public       | jsonb                       |          | f        |                |

(35 rows)
