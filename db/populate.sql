CREATE TABLE `first_table` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `first_col` VARCHAR(64) NULL,
    `second_col` DATE NULL
);

INSERT INTO first_table VALUES (0, 'hello', '2011-03-12');
INSERT INTO first_table VALUES (1, 'foo', '2016-08-30');
INSERT INTO first_table VALUES (2, 'bar', '2013-08-23');
