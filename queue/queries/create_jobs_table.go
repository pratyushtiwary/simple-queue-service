package queries

const JOB_CREATE_QUERY = `
	create table if not exists jobs(
		id text not null primary key,
		status text,
		detail_id text,
		createdAt DATE DEFAULT (datetime('now','localtime')),
		constraint fk_detail
		foreign key (detail_id)
		references details(id)
	)
`
