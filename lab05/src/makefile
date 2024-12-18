DATA_DIR="data/"

json-export:
	psql -U postgres -d bmstu_db -f ./db-json-export.sql

json-import:
	psql -U postgres -d bmstu_db -f ./json-db-import.sql

json-update:
	psql -U postgres -d bmstu_db -f ./json-db-update.sql

main:
	psql -U postgres -d bmstu_db -f ./main.sql

