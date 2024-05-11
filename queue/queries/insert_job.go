package queries

const INSERT_JOB_QUERY = `
	insert into jobs(id, status, detail_id, createdAt) values (
		?,
		"pending",
		?,
		datetime('now','localtime')
	)
`
