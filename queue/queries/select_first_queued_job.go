package queries

const SELECT_FIRST_QUEUED_JOB_QUERY = `
	select id, detail_id, createdAt from jobs
	where status = ? order by createdAt ASC limit 1;
`
