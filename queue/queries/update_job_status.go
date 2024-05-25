package queries

const UPDATE_JOB_STATUS_QUERY = `
	update jobs set status = ?
	where id = ?
`
