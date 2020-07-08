alter table shop drop column name_norm;
alter table history.shop drop column name_norm;

create table  shop_search(
	"id" int8,
	"name_norm" tsvector,
	"name" text
)
