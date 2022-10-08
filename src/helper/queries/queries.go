package queries

//Login queries
const LoginQuery = "select id, email, password from usr where  email = $1 and is_active"

//Register queries
const CreateUsrQuery = "insert into usr ( password, email, shortening_url_limit, account_type, zlins_dttm) values ($1, $2, $3, $4, current_timestamp) returning id"

//Url queries
const CreateUrlQuery = "INSERT INTO url (long_version, shortened_version, usr_id, zlins_dttm) VALUES ($1, $2, $3, current_timestamp)"

const ListUrlQuery = "SELECT id,long_version,shortened_version, usr_id FROM url WHERE usr_id = $1"

const DeleteUrlQuery = "DELETE FROM url WHERE usr_id = $1 AND id = $2"

const GetLongUrlFromShortenedQuery = "SELECT long_version FROM url WHERE shortened_version = $1 and usr_id = $2"

const CheckIfUrlExistsQuery = "SELECT EXISTS(SELECT 1 FROM url WHERE id = $1)"
