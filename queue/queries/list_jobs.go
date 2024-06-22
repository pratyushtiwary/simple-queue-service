package queries

const LIST_JOB_QUERY = `
	select id from jobs
	where status = ?
`
