package dbcodegen

const tableInfoQuery = `
SELECT pg_stat_user_tables.relname                         AS table_name,
       (SELECT pg_description.description
        FROM pg_description
        WHERE pg_stat_user_tables.relid = pg_description.objoid
          AND pg_description.objsubid = 0)                 AS table_comment,
       information_schema.columns.column_name              AS column_name,
       information_schema.columns.column_default           AS column_default,
       pg_description.description                          AS column_comment,
       information_schema.columns.table_schema             AS table_schema,
       information_schema.columns.data_type                AS column_data_type,
	   information_schema.columns.is_generated             AS is_generated,
	   information_schema.columns.generation_expression    AS generation_expression,
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
                       AND pg_stat_user_tables.schemaname = information_schema.columns.table_schema
         LEFT JOIN pg_description
                   ON pg_description.objoid = pg_stat_user_tables.relid
                       AND pg_description.objsubid = information_schema.columns.ordinal_position
WHERE information_schema.columns.table_schema NOT IN ('pg_catalog', 'information_schema') 
--   AND table_name != 'schema_migrations'
ORDER BY table_schema, table_name, information_schema.columns.ordinal_position;
`

const primaryKeyInfoQuery = `
SELECT a.attname
FROM   pg_index i
JOIN   pg_attribute a ON a.attrelid = i.indrelid
					 AND a.attnum = ANY(i.indkey)
WHERE  i.indrelid = '%v'::regclass
AND    i.indisprimary;
`

const tableIndexQuery = `
SELECT n.nspname                                  AS schema_name,
	   t.relname                                  AS table_name,
	   ix.relname                                 AS index_name,
	   i.indisunique                              AS is_unique,
	   i.indisprimary                             AS pkey,
	   ARRAY_TO_STRING(ARRAY_AGG(a.attname), ',') AS column_names
FROM pg_class t,
	 pg_class ix,
	 pg_index i,
	 pg_attribute a,
	 pg_namespace n
WHERE t.relnamespace = n.oid
  AND t.oid = i.indrelid
  AND ix.oid = i.indexrelid
  AND a.attrelid = t.oid
  AND a.attnum = ANY (i.indkey)
  AND t.relkind = 'r'
  AND n.nspname != 'pg_catalog'
  AND n.nspname = '%v'
  AND t.relname = '%v'
GROUP BY n.nspname,
		 t.relname,
		 ix.relname,
		 i.indisunique,
		 i.indisprimary
ORDER BY n.nspname,
		 t.relname,
		 ix.relname;	
`

const tableForeignKeyQuery = `
SELECT kcu.table_name       AS from_table,
	   kcu.table_schema     AS from_schema,
	   kcu.column_name      AS from_column,
	   rel_kcu.table_name   AS to_table,
	   rel_kcu.table_schema AS to_schema,
	   rel_kcu.column_name  AS to_column,
	   kcu.ordinal_position AS ordinal_position,
	   kcu.constraint_name
FROM information_schema.table_constraints tco
		 JOIN information_schema.key_column_usage kcu
			  ON tco.constraint_schema = kcu.constraint_schema
				  AND tco.constraint_name = kcu.constraint_name
		 JOIN information_schema.referential_constraints rco
			  ON tco.constraint_schema = rco.constraint_schema
				  AND tco.constraint_name = rco.constraint_name
		 JOIN information_schema.key_column_usage rel_kcu
			  ON rco.unique_constraint_schema = rel_kcu.constraint_schema
				  AND rco.unique_constraint_name = rel_kcu.constraint_name
				  AND kcu.ordinal_position = rel_kcu.ordinal_position
WHERE tco.constraint_type = 'FOREIGN KEY'
ORDER BY kcu.table_schema,
		 kcu.table_name,
		 kcu.constraint_name,
		 kcu.ordinal_position;
`
