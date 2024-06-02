package queries

const DETAIL_CREATE_QUERY = `
	create table if not exists details(
		id text not null primary key,
		data blob
	)
`
