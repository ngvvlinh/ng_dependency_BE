create table ticket (
	id int8 primary key,
	code text,
	assigned_user_ids int8[],
	account_id int8,
	label_ids int8[],
	title text,
	description text,
	note text,
	ref_id int8,
	ref_type int2,
	"source" int2,
	state int2,
	status int2,
	created_by int8,
	updated_by int8,
	confirmed_by int8,
	closed_by int8,
	created_at timestamptz,
	updated_at timestamptz,
	confirmed_at timestamptz,
	closed_at timestamptz
);

create table ticket_search (
	id int8 primary key,
	title_norm tsvector
);

create table ticket_comment(
	id int8 primary key,
	ticket_id int8,
	account_id int8,
	created_by int8,
	message text,
	image_url text,
	parent_id int8,
	description text,
	deleted_at timestamptz,
	deleted_by int8,
	created_at timestamptz,
	updated_at timestamptz
);

create table ticket_label(
	id int8,
	"name" text,
	parent_id int8
);


insert into ticket_label(id, "name", parent_id)
values
(1114377203503360427, 'Topship', 0),
(1143772035031952261, 'Đơn giao hàng', 	1114377203503360427),
(1143772035038117662, 'Giục lấy', 		1143772035031952261),
(1143772035039707243, 'Giục giao', 		1143772035031952261),
(1143772035041002422, 'Giục giao lại', 	1143772035031952261),
(1143772035039984275, 'Đổi COD',	 	1143772035031952261),
(1143772113550994045, 'Đổi SDT người nhận', 1143772035031952261),
(1143772113554218820, 'Đổi Tên người nhận', 1143772035031952261),
(1143772113555730296, 'Yêu cầu khác', 		1143772035031952261),
(1143772113558311281, 'Đối soát', 			1114377203503360427),
(1143772113563232772, 'Số dư', 				1114377203503360427),
(1143772113562859008, 'Hướng dẫn sử dụng', 0),
(1143772113554509729, 'Khác', 0),
(1143772113551196577, 'eTop', 0),
(1143772113562488928, 'Subscription', 		1143772113551196577),
(1143772113548134095, 'Tài khoản', 			1143772113551196577),
(1143772113555377865, 'eComify', 0),
(1143772113551161883, 'Subscription', 		1143772113555377865),
(1143772113558398845, 'eOrder', 0),
(1143772113562946217, 'eB2B (trading)', 0),
(1143772113561348063, 'Đơn hàng', 			1143772113562946217),
(1143772113560893102, 'Internal', 0);
