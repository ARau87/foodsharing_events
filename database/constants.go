package database

const (

	MYSQL_INSERT_USER = "INSERT INTO users (email, password, firstname, lastname, updated) VALUES (?, ?, ?, ?, UTC_TIMESTAMP())"
	MYSQL_SELECT_USER = "SELECT users.id, users.email, users.firstname, users.lastname FROM users WHERE users.id = ? "
	MYSQL_SELECT_USER_BY_CREDENTIALS = "SELECT users.id, users.email, users.firstname, users.lastname FROM users WHERE users.email = ? AND users.password = ? "

	MYSQL_INSERT_EVENT = "INSERT INTO events (name, description, city, address, gpsLat, gpsLong, creatorId, maxParticipants, updated) VALUES (?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP())"
	MYSQL_SELECT_EVENT_BY_ID = "SELECT events.id, events.name, events.description, events.city, events.address, events.gpsLat, events.gpsLong, events.creatorId, events.maxParticipants FROM events WHERE events.id = ?"
	MYSQL_SELECT_EVENT_PARTICIPANTS_BY_ID = "SELECT users.id, users.email, users.firstname, users.lastname FROM events_users INNER JOIN users ON events_users.userId = users.id WHERE events_users.eventId = ? "

	MYSQL_INSERT_EVENT_USER = "INSERT INTO events_users (userId, eventId) VALUES (?, ?)"
	)
