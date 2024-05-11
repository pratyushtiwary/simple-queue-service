package queries

const INSERT_DETAIL_QUERY = `
	insert into details(id, data) values (
		?,
		?
	)
`
