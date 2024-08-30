package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pingcap/parser/ast"
	formater "github.com/pingcap/parser/format"
	"github.com/stretchr/testify/assert"
)

func Test_ParseNewSupportSQL(t *testing.T) {
	type args struct {
		stmt string
	}
	tests := []struct {
		name       string
		args       args
		wantNewSQL string
		wantErr    bool
	}{
		{
			name: "geometry type check",
			args: args{
				stmt: "CREATE TABLE `gis_table` (  `id` bigint NOT NULL, `gis` geometry NOT NULL COMMENT '空间位置信息',  " +
					"PRIMARY KEY (`id`),  SPATIAL KEY `gis_index` (`gis`))",
			},
			wantNewSQL: "CREATE TABLE `gis_table` (`id` BIGINT NOT NULL,`gis` GEOMETRY NOT NULL COMMENT '空间位置信息',PRIMARY KEY(`id`),SPATIAL KEY `gis_index`(`gis`))",
			wantErr:    false,
		},
		{
			name: "default CURRENT_TIMESTAMP check",
			args: args{
				stmt: "CREATE TABLE `json_table` (  `id` bigint NOT NULL, outsource_details json NOT NULL DEFAULT " +
					"CURRENT_TIMESTAMP() COMMENT 'product_outsource_details', " +
					"PRIMARY KEY (`id`))",
			},
			wantNewSQL: "CREATE TABLE `json_table` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT 'product_outsource_details',PRIMARY KEY(`id`))",
			wantErr:    false,
		},
		{
			name: "default CURRENT_TIMESTAMP in parentheses check",
			args: args{
				stmt: "CREATE TABLE `json_table` (  `id` bigint NOT NULL, outsource_details json NOT NULL DEFAULT " +
					"(CURRENT_TIMESTAMP()) COMMENT 'product_outsource_details', " +
					"PRIMARY KEY (`id`))",
			},
			wantNewSQL: "CREATE TABLE `json_table` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT (CURRENT_TIMESTAMP()) COMMENT 'product_outsource_details',PRIMARY KEY(`id`))",
			wantErr:    false,
		},
		{
			name: "string literal check",
			args: args{
				stmt: "CREATE TABLE `json_table` (  `id` bigint NOT NULL, outsource_details json NOT NULL DEFAULT (" +
					"'haha') COMMENT 'product_outsource_details', " +
					"PRIMARY KEY (`id`))",
			},
			wantNewSQL: "CREATE TABLE `json_table` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT ('haha') COMMENT 'product_outsource_details',PRIMARY KEY(`id`))",
			wantErr:    false,
		},
		{
			name: "JSON_OBJECT check",
			args: args{
				stmt: "CREATE TABLE `json_table` (  `id` bigint NOT NULL, outsource_details json NOT NULL DEFAULT JSON_OBJECT() COMMENT 'product_outsource_details', PRIMARY KEY (`id`))",
			},
			wantNewSQL: "CREATE TABLE `json_table` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT JSON_OBJECT() COMMENT 'product_outsource_details',PRIMARY KEY(`id`))",
			wantErr:    false,
		},
		{
			name: "JSON_OBJECT in parentheses check",
			args: args{
				stmt: "CREATE TABLE `json_table` (  `id` bigint NOT NULL, outsource_details json NOT NULL DEFAULT (JSON_OBJECT()) COMMENT 'product_outsource_details', PRIMARY KEY (`id`))",
			},
			wantNewSQL: "CREATE TABLE `json_table` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT (JSON_OBJECT()) COMMENT 'product_outsource_details',PRIMARY KEY(`id`))",
			wantErr:    false,
		},
		{
			name: "JSON_OBJECT with param check",
			args: args{
				stmt: "CREATE TABLE `json_table` (  `id` bigint NOT NULL, outsource_details json NOT NULL DEFAULT (" +
					"JSON_OBJECT(_utf8mb4'id',87,_utf8mb4'name',_utf8mb4'carrot')) COMMENT 'product_outsource_details', " +
					"PRIMARY KEY (`id`))",
			},
			wantNewSQL: "CREATE TABLE `json_table` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT (JSON_OBJECT(_UTF8MB4'id', 87, _UTF8MB4'name', _UTF8MB4'carrot')) COMMENT 'product_outsource_details',PRIMARY KEY(`id`))",
			wantErr:    false,
		},
		{
			name: "JSON_ARRAY check",
			args: args{
				stmt: "CREATE TABLE `json_table` (  `id` bigint NOT NULL, " +
					"outsource_details json NOT NULL DEFAULT JSON_ARRAY() " +
					"COMMENT 'product_outsource_details', " +
					"PRIMARY KEY (`id`))",
			},
			wantNewSQL: "CREATE TABLE `json_table` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT JSON_ARRAY() COMMENT 'product_outsource_details',PRIMARY KEY(`id`))",
			wantErr:    false,
		},
		{
			name: "JSON_ARRAY in parentheses check",
			args: args{
				stmt: "CREATE TABLE `json_table` (  `id` bigint NOT NULL, " +
					"outsource_details json NOT NULL DEFAULT (JSON_ARRAY()) " +
					"COMMENT 'product_outsource_details', " +
					"PRIMARY KEY (`id`))",
			},
			wantNewSQL: "CREATE TABLE `json_table` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT (JSON_ARRAY()) COMMENT 'product_outsource_details',PRIMARY KEY(`id`))",
			wantErr:    false,
		},
		{
			name: "JSON_ARRAY check with param check",
			args: args{
				stmt: "CREATE TABLE `json_table3` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT(JSON_ARRAY(1, 2, 3)) " +
					"COMMENT 'product_outsource_details',PRIMARY KEY(`id`)) ENGINE = InnoDB DEFAULT CHARACTER SET = UTF8MB4 DEFAULT COLLATE = UTF8MB4_0900_AI_CI",
			},
			wantNewSQL: "CREATE TABLE `json_table3` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT (JSON_ARRAY(1, 2, 3)) COMMENT 'product_outsource_details',PRIMARY KEY(`id`)) ENGINE = InnoDB DEFAULT CHARACTER SET = UTF8MB4 DEFAULT COLLATE = UTF8MB4_0900_AI_CI",
			wantErr:    false,
		},
		{
			name: "JSON_QUOTE check",
			args: args{
				stmt: "CREATE TABLE `json_table` (  `id` bigint NOT NULL, outsource_details json NOT NULL DEFAULT " +
					"JSON_QUOTE('haha') COMMENT 'product_outsource_details', " +
					"PRIMARY KEY (`id`))",
			},
			wantNewSQL: "CREATE TABLE `json_table` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT JSON_QUOTE(_UTF8MB4'haha') COMMENT 'product_outsource_details',PRIMARY KEY(`id`))",
			wantErr:    false,
		},
		{
			name: "JSON_QUOTE in parentheses check",
			args: args{
				stmt: "CREATE TABLE `json_table` (  `id` bigint NOT NULL, outsource_details json NOT NULL DEFAULT (" +
					"JSON_QUOTE('haha')) COMMENT 'product_outsource_details', " +
					"PRIMARY KEY (`id`))",
			},
			wantNewSQL: "CREATE TABLE `json_table` (`id` BIGINT NOT NULL,`outsource_details` JSON NOT NULL DEFAULT (JSON_QUOTE(_UTF8MB4'haha')) COMMENT 'product_outsource_details',PRIMARY KEY(`id`))",
			wantErr:    false,
		},
	}

	sqlParser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := sqlParser.ParseOneStmt(tt.args.stmt, "", "")
			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("ParseOneStmt() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			gotNewSQL := RestoreSQLNode(stmt)
			assert.ErrorIs(t, err, nil)
			assert.Equal(t, tt.wantNewSQL, gotNewSQL)
		})
	}
}

func RestoreSQLNode(node ast.Node) string {
	writer := &strings.Builder{}
	restoreCtx := formater.NewRestoreCtx(
		formater.RestoreStringSingleQuotes|formater.RestoreKeyWordUppercase|formater.RestoreNameBackQuotes|formater.RestoreBracketAroundBinaryOperation,
		writer)
	err := node.Restore(restoreCtx)
	if err != nil {
		panic(fmt.Sprintf("restore SQL node %s failed <-- %s", node.Text(), err))
	}
	return writer.String()
}
