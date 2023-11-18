package main

import (
	`fmt`

	`github.com/techrail/ground/dbCodegen`
)

func main() {
	cnf := dbCodegen.CodegenConfig{
		TablePackageName: "",
		TablePackagePath: "",
		PgDbUrl:          "postgres://vaibhav:vaibhav@127.0.0.1:5432/twitter_clon?sslmode=disable",
	}
	g, e := dbCodegen.NewCodeGenerator(cnf)
	if e.IsNotBlank() {
		fmt.Printf("I#1NPKZR - Some error when creating new codegenerator: %v\n", e)
	}
	errTy := g.Connect()
	if errTy.IsNotBlank() {
		fmt.Printf("I#1NPLCJ - %v\n", errTy)
	}
}

// SQL script to get the list of tables:
/*
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
*/
