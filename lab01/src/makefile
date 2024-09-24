SQL_DIR := ./sql

all: build insert

build:
	mkdir -p data
	psql -U postgres -d bmstu_db -f $(SQL_DIR)/create.sql
	psql -U postgres -d bmstu_db -f $(SQL_DIR)/constraints.sql
	python3 generator.py

insert:
	psql -U postgres -d bmstu_db -f $(SQL_DIR)/bulk_copy.sql
	
drop:
	psql -U postgres -d bmstu_db -f $(SQL_DIR)/down.sql

clean:
	rm -rf ./data

rebuild: drop clean build