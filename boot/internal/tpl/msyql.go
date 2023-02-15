package tpl

var Mysql = `CREATE TABLE ` + "`${tableName}`" + ` (
${fields}
${keys}
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='${comment}';
`
