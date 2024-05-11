package queries

const SELECT_JOB_QUERY = `
	select status, detail_id, createdAt from jobs
	where id = ?
`
